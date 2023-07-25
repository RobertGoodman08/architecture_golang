package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"

	"architecture_go/pkg/tools/transaction"
	"architecture_go/pkg/type/columnCode"
	"architecture_go/pkg/type/context"
	log "architecture_go/pkg/type/logger"
	"architecture_go/pkg/type/queryParameter"
	"architecture_go/services/contact/internal/domain/contact"
	"architecture_go/services/contact/internal/repository/storage/postgres/dao"
	"architecture_go/services/contact/internal/useCase"
)

var mappingSortContact = map[columnCode.ColumnCode]string{
	"id":          "id",
	"fullName":    "full_name",
	"phoneNumber": "phone_number",
	"name":        "name",
	"surname":     "surname",
	"patronymic":  "patronymic",
	"email":       "email",
	"gender":      "gender",
	"age":         "age",
}

func (r *Repository) CreateContact(c context.Context, contacts ...*contact.Contact) ([]*contact.Contact, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	response, err := r.createContactTx(ctx, tx, contacts...)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Repository) createContactTx(ctx context.Context, tx pgx.Tx, contacts ...*contact.Contact) ([]*contact.Contact, error) {
	if len(contacts) == 0 {
		return []*contact.Contact{}, nil
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"slurm", "contact"},
		dao.CreateColumnContact,
		r.toCopyFromSource(contacts...))
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	return contacts, nil
}

func (r *Repository) UpdateContact(c context.Context, ID uuid.UUID, updateFn func(c *contact.Contact) (*contact.Contact, error)) (*contact.Contact, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	upContact, err := r.oneContactTx(ctx, tx, ID)
	if err != nil {
		return nil, err
	}
	in, err := updateFn(upContact)
	if err != nil {
		return nil, err
	}

	return r.updateContactTx(ctx, tx, in)
}

func (r *Repository) updateContactTx(ctx context.Context, tx pgx.Tx, in *contact.Contact) (*contact.Contact, error) {

	builder := r.genSQL.Update("slurm.contact").
		Set("email", in.Email().String()).
		Set("phone_number", in.PhoneNumber().String()).
		Set("age", in.Age()).
		Set("gender", in.Gender()).
		Set("modified_at", in.ModifiedAt()).
		Set("name", in.Name().String()).
		Set("surname", in.Surname().String()).
		Set("patronymic", in.Patronymic().String()).
		Where(squirrel.And{
			squirrel.Eq{
				"id":          in.ID(),
				"is_archived": false,
			},
		}).
		Suffix(`RETURNING
			id,
			created_at,
			modified_at,
			phone_number,
			email,
			name,
			surname,
			patronymic,
			age,
			gender`,
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	var daoContacts []*dao.Contact
	if err = pgxscan.ScanAll(&daoContacts, rows); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	return r.toDomainContact(daoContacts[0])
}

func (r *Repository) DeleteContact(c context.Context, ID uuid.UUID) error {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	if err = r.deleteContactTx(ctx, tx, ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) deleteContactTx(ctx context.Context, tx pgx.Tx, ID uuid.UUID) error {
	builder := r.genSQL.Update("slurm.contact").
		Set("is_archived", true).
		Set("modified_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_archived": false, "id": ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	var daoContacts []*dao.Contact
	if err = pgxscan.ScanAll(&daoContacts, rows); err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	if err = r.updateGroupsContactCountByFilters(ctx, tx, ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListContact(c context.Context, parameter queryParameter.QueryParameter) ([]*contact.Contact, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	span, tmp := opentracing.StartSpanFromContext(c, "ListContact")
	defer span.Finish()
	ctx = context.New(tmp)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	if parameter.Pagination.Limit == 0 {
		parameter.Pagination.Limit = r.options.DefaultLimit
	}

	contacts, err := r.listContactTx(ctx, tx, parameter)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *Repository) listContactTx(ctx context.Context, tx pgx.Tx, parameter queryParameter.QueryParameter) ([]*contact.Contact, error) {
	var builder = r.genSQL.Select(
		"id",
		"created_at",
		"modified_at",
		"phone_number",
		"email",
		"name",
		"surname",
		"patronymic",
		"age",
		"gender",
	).From("slurm.contact")

	builder = builder.Where(squirrel.Eq{"is_archived": false})

	if len(parameter.Sorts) > 0 {
		builder = builder.OrderBy(parameter.Sorts.Parsing(mappingSortContact)...)
	} else {
		builder = builder.OrderBy("created_at DESC")
	}

	if parameter.Pagination.Limit > 0 {
		builder = builder.Limit(parameter.Pagination.Limit)
	}
	if parameter.Pagination.Offset > 0 {
		builder = builder.Offset(parameter.Pagination.Offset)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	var daoContacts []*dao.Contact
	if err = pgxscan.ScanAll(&daoContacts, rows); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	return r.toDomainContacts(daoContacts)
}

func (r *Repository) ReadContactByID(c context.Context, ID uuid.UUID) (response *contact.Contact, err error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	response, err = r.oneContactTx(ctx, tx, ID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Repository) oneContactTx(ctx context.Context, tx pgx.Tx, ID uuid.UUID) (*contact.Contact, error) {
	var builder = r.genSQL.Select(
		"id",
		"created_at",
		"modified_at",
		"phone_number",
		"email",
		"name",
		"surname",
		"patronymic",
		"age",
		"gender",
	).From("slurm.contact")

	builder = builder.Where(squirrel.Eq{"is_archived": false, "id": ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	var daoContact []*dao.Contact
	if err = pgxscan.ScanAll(&daoContact, rows); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	if len(daoContact) == 0 {
		return nil, useCase.ErrContactNotFound
	}

	return r.toDomainContact(daoContact[0])
}

func (r *Repository) CountContact(c context.Context) (uint64, error) {

	span, tmp := opentracing.StartSpanFromContext(c, "CountContact")
	defer span.Finish()
	ctx := context.New(tmp)

	var builder = r.genSQL.Select(
		"COUNT(id)",
	).From("slurm.contact")

	builder = builder.Where(squirrel.Eq{"is_archived": false})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, log.ErrorWithContext(ctx, err)
	}

	var row = r.db.QueryRow(ctx, query, args...)
	var total uint64

	if err = row.Scan(&total); err != nil {
		return 0, log.ErrorWithContext(ctx, err)
	}

	return total, nil
}
