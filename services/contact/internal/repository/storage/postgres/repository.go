package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/pressly/goose"

	log "architecture_go/pkg/type/logger"
)

func init() {
	viper.SetDefault("MIGRATIONS_DIR", "./services/contact/internal/repository/storage/postgres/migrations")
}

type Repository struct {
	db      *pgxpool.Pool
	genSQL  squirrel.StatementBuilderType
	options Options
}

type Options struct {
	Timeout       time.Duration
	DefaultLimit  uint64
	DefaultOffset uint64
}

func New(db *pgxpool.Pool, o Options) (*Repository, error) {
	if err := migrations(db); err != nil {
		return nil, err
	}

	var r = &Repository{
		genSQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		db:     db,
	}

	r.SetOptions(o)
	return r, nil
}

func (r *Repository) SetOptions(options Options) {
	if options.DefaultLimit == 0 {
		options.DefaultLimit = 10
		log.Debug("set default options.DefaultLimit", zap.Any("defaultLimit", options.DefaultLimit))
	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 30
		log.Debug("set default options.Timeout", zap.Any("timeout", options.Timeout))
	}

	if r.options != options {
		r.options = options
		log.Info("set new options", zap.Any("options", r.options))
	}
}

func migrations(pool *pgxpool.Pool) (err error) {
	db, err := goose.OpenDBWithDriver("postgres", pool.Config().ConnConfig.ConnString())
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			log.Error(errClose)
			err = errClose
			return
		}
	}()

	dir := viper.GetString("MIGRATIONS_DIR")
	goose.SetTableName("contact_version")
	if err = goose.Run("up", db, dir); err != nil {
		log.Error(err, zap.String("command", "up"))
		return err
	}
	return
}
