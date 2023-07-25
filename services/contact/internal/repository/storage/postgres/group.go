package postgres

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"architecture_go/pkg/tools/converter"
	"architecture_go/pkg/tools/transaction"
	"architecture_go/pkg/type/columnCode"
	"architecture_go/pkg/type/context"
	log "architecture_go/pkg/type/logger"
	"architecture_go/pkg/type/queryParameter"
	"architecture_go/services/contact/internal/domain/group"
	"architecture_go/services/contact/internal/repository/storage/postgres/dao"
	"architecture_go/services/contact/internal/useCase"
)

var mappingSortGroup = map[columnCode.ColumnCode]string{
	"id":           "id",
	"name":         "name",
	"description":  "description",
	"contactCount": "contact_count",
}

func (r *Repository) CreateGroup(c context.Context, group *group.Group) (*group.Group, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	query, args, err := r.genSQL.Insert("slurm.group").
		Columns(
			"id",
			"name",
			"description",
			"created_at",
			"modified_at",
		).
		Values(
			group.ID(),
			group.Name().Value(),
			group.Description().Value(),
			group.CreatedAt(),
			group.ModifiedAt()).
		ToSql()
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}
	return group, nil
}

func (r *Repository) UpdateGroup(c context.Context, ID uuid.UUID, updateFn func(group *group.Group) (*group.Group, error)) (*group.Group, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	upGroup, err := r.oneGroupTx(ctx, tx, ID)
	if err != nil {
		return nil, err
	}
	groupForUpdate, err := updateFn(upGroup)
	if err != nil {
		return nil, err
	}

	query, args, err := r.genSQL.Update("slurm.group").
		Set("name", groupForUpdate.Name().Value()).
		Set("description", groupForUpdate.Description().Value()).
		Set("modified_at", groupForUpdate.ModifiedAt()).
		Where(squirrel.And{
			squirrel.Eq{
				"id":          ID,
				"is_archived": false,
			},
		}).
		Suffix(`RETURNING
			id,
			name,
			description,
			created_at,
			modified_at`,
		).
		ToSql()
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	var daoGroup []*dao.Group
	if err = pgxscan.ScanAll(&daoGroup, rows); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	return groupForUpdate, nil
}

func (r *Repository) DeleteGroup(c context.Context, ID uuid.UUID) error {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	if err = r.deleteGroupTx(ctx, tx, ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) deleteGroupTx(ctx context.Context, tx pgx.Tx, ID uuid.UUID) error {
	query, args, err := r.genSQL.Update("slurm.group").
		Set("is_archived", true).
		Set("modified_at", time.Now().UTC()).
		Where(squirrel.Eq{
			"id":          ID,
			"is_archived": false,
		}).ToSql()

	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	if err = r.clearGroupTx(ctx, tx, ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) clearGroupTx(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) error {
	query, args, err := r.genSQL.
		Delete("slurm.contact_in_group").
		Where(squirrel.Eq{"group_id": groupID}).
		ToSql()
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	if err = r.updateGroupContactCount(ctx, tx, groupID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListGroup(c context.Context, parameter queryParameter.QueryParameter) ([]*group.Group, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	response, err := r.listGroupTx(ctx, tx, parameter)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Repository) listGroupTx(ctx context.Context, tx pgx.Tx, parameter queryParameter.QueryParameter) ([]*group.Group, error) {
	var result []*group.Group

	var builder = r.genSQL.Select(
		"id",
		"name",
		"description",
		"created_at",
		"modified_at",
		"contact_count",
		"is_archived",
	).
		From("slurm.group")

	builder = builder.Where(squirrel.Eq{"is_archived": false})

	if len(parameter.Sorts) > 0 {
		builder = builder.OrderBy(parameter.Sorts.Parsing(mappingSortGroup)...)
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
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	var groups []*dao.Group
	if err = pgxscan.ScanAll(&groups, rows); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	for _, g := range groups {
		domainGroup, err := g.ToDomainGroup()
		if err != nil {
			return nil, log.ErrorWithContext(ctx, err)
		}
		result = append(result, domainGroup)
	}
	return result, nil
}

func (r *Repository) ReadGroupByID(c context.Context, ID uuid.UUID) (*group.Group, error) {

	ctx := c.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	response, err := r.oneGroupTx(ctx, tx, ID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Repository) oneGroupTx(ctx context.Context, tx pgx.Tx, ID uuid.UUID) (response *group.Group, err error) {

	var builder = r.genSQL.Select(
		"id",
		"name",
		"description",
		"created_at",
		"modified_at",
		"contact_count",
		"is_archived",
	).
		From("slurm.group")

	builder = builder.Where(squirrel.Eq{"is_archived": false, "id": ID})

	query, args, err := builder.ToSql()
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	var daoGroup []*dao.Group
	if err = pgxscan.ScanAll(&daoGroup, rows); err != nil {
		return nil, log.ErrorWithContext(ctx, err)
	}

	if len(daoGroup) == 0 {
		return nil, useCase.ErrGroupNotFound
	}

	return daoGroup[0].ToDomainGroup()

}

func (r *Repository) CountGroup(ctx context.Context) (uint64, error) {
	var builder = r.genSQL.Select(
		"COUNT(id)",
	).From("slurm.group")

	builder = builder.Where(squirrel.Eq{"is_archived": false})

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, log.ErrorWithContext(ctx, err)
	}

	row := r.db.QueryRow(ctx, query, args...)
	var total uint64

	if err = row.Scan(&total); err != nil {
		return 0, log.ErrorWithContext(ctx, err)
	}

	return total, nil
}

func (r *Repository) updateGroupsContactCountByFilters(ctx context.Context, tx pgx.Tx, ID uuid.UUID) error {

	builder := r.genSQL.Select("contact_in_group.group_id").
		From("slurm.contact_in_group").
		InnerJoin("slurm.contact ON contact_in_group.contact_id = contact.id").
		GroupBy("contact_in_group.group_id")

	builder = builder.Where(squirrel.Eq{"contact_in_group.contact_id": ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}
	var groupIDs []uuid.UUID
	for rows.Next() {
		var groupID sql.NullString
		if err = rows.Scan(&groupID); err != nil {
			return log.ErrorWithContext(ctx, err)
		}
		groupIDs = append(groupIDs, converter.StringToUUID(groupID.String))
	}

	for _, groupID := range groupIDs {
		if err = r.updateGroupContactCount(ctx, tx, groupID); err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	return nil
}

func (r *Repository) updateGroupContactCount(ctx context.Context, tx pgx.Tx, groupID uuid.UUID) error {
	subSelect := r.genSQL.Select("count(contact_in_group.id)").
		From("slurm.contact_in_group").
		InnerJoin("slurm.contact ON contact_in_group.contact_id = contact.id").
		Where(squirrel.Eq{"group_id": groupID, "is_archived": false})

	query, _, err := r.genSQL.
		Update("slurm.group").
		Set("contact_count", subSelect).
		Where(squirrel.Eq{"id": groupID}).
		ToSql()
	if err != nil {
		return log.ErrorWithContext(ctx, err)
	}

	var args = []interface{}{groupID, false}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return log.ErrorWithContext(ctx, err)
	}
	return nil
}
