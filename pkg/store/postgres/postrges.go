package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	_ "github.com/spf13/viper/remote"
)

func init() {
	if err := initDefaultEnv(); err != nil {
		panic(err)
	}
}

//  Env to pg params
// 	"PGHOST":               "host",
// 	"PGPORT":               "port",
// 	"PGDATABASE":           "database",
// 	"PGUSER":               "user",
// 	"PGPASSWORD":           "password",
// 	"PGPASSFILE":           "passfile",
// 	"PGAPPNAME":            "application_name",
// 	"PGCONNECT_TIMEOUT":    "connect_timeout",
// 	"PGSSLMODE":            "sslmode",
// 	"PGSSLKEY":             "sslkey",
// 	"PGSSLCERT":            "sslcert",
// 	"PGSSLROOTCERT":        "sslrootcert",
// 	"PGTARGETSESSIONATTRS": "target_session_attrs",
// 	"PGSERVICE":            "service",
// 	"PGSERVICEFILE":        "servicefile",
func initDefaultEnv() error {
	if len(os.Getenv("PGHOST")) == 0 {
		if err := os.Setenv("PGHOST", "postgres"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGPORT")) == 0 {
		if err := os.Setenv("PGPORT", "5838"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGDATABASE")) == 0 {
		if err := os.Setenv("PGDATABASE", "postgres"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGUSER")) == 0 {
		if err := os.Setenv("PGUSER", "postgres"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGPASSWORD")) == 0 {
		if err := os.Setenv("PGPASSWORD", "password"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGSSLMODE")) == 0 {
		if err := os.Setenv("PGSSLMODE", "disable"); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

type Store struct {
	Pool *pgxpool.Pool
}

type Settings struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	SSLMode  string
}

func (s Settings) toDSN() string {
	var args []string

	if len(s.Host) > 0 {
		args = append(args, fmt.Sprintf("host=%s", s.Host))
	}

	if s.Port > 0 {
		args = append(args, fmt.Sprintf("port=%d", s.Port))
	}

	if len(s.Database) > 0 {
		args = append(args, fmt.Sprintf("dbname=%s", s.Database))
	}

	if len(s.User) > 0 {
		args = append(args, fmt.Sprintf("user=%s", s.User))
	}

	if len(s.Password) > 0 {
		args = append(args, fmt.Sprintf("password=%s", s.Password))
	}

	if len(s.SSLMode) > 0 {
		args = append(args, fmt.Sprintf("sslmode=%s", s.SSLMode))
	}

	return strings.Join(args, " ")
}

func New(settings Settings) (*Store, error) {

	config, err := pgxpool.ParseConfig(settings.toDSN())
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{Pool: conn}, nil
}
