package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"architecture_go/pkg/tools/transaction"
	"architecture_go/services/contact/internal/domain/contact"
	"architecture_go/services/contact/internal/repository/storage/postgres/dao"
)

func (r *Repository) CreateContactIntoGroup(groupID uuid.UUID, contacts ...*contact.Contact) ([]*contact.Contact, error) {
	var ctx = context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	response, err := r.repoContact.CreateContactTx(ctx, tx, contacts...)
	if err != nil {
		return nil, err
	}
	var contactIDs = make([]uuid.UUID, len(response))
	for i, c := range response {
		contactIDs[i] = c.ID()
	}

	if err = r.fillGroupTx(ctx, tx, groupID, contactIDs...); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Repository) DeleteContactFromGroup(groupID, contactID uuid.UUID) error {
	var ctx = context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	query, args, err := r.genSQL.
		Delete("slurm.contact_in_group").
		Where(squirrel.Eq{"contact_id": contactID, "group_id": groupID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if err = r.updateGroupContactCount(ctx, tx, groupID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddContactsToGroup(groupID uuid.UUID, contactIDs ...uuid.UUID) error {
	var ctx = context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(ctx context.Context, t pgx.Tx) {
		err = transaction.Finish(ctx, t, err)
	}(ctx, tx)

	if err = r.fillGroupTx(ctx, tx, groupID, contactIDs...); err != nil {
		return err
	}
	return nil
}

func (r *Repository) fillGroupTx(ctx context.Context, tx pgx.Tx, groupID uuid.UUID, contactIDs ...uuid.UUID) error {
	_, mapExist, err := r.checkExistContactInGroup(ctx, tx, groupID, contactIDs...)
	if err != nil {
		return err
	}

	for i := 0; i < len(contactIDs); {
		var contactID = contactIDs[i]
		if exist := mapExist[contactID]; exist {
			contactIDs[i] = contactIDs[len(contactIDs)-1]
			contactIDs = contactIDs[:len(contactIDs)-1]
			continue
		}
		i++
	}

	if len(contactIDs) == 0 {
		return nil
	}

	var rows [][]interface{}
	var timeNow = time.Now().UTC()
	for _, contactID := range contactIDs {
		rows = append(rows, []interface{}{
			timeNow,
			timeNow,
			groupID,
			contactID,
		})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"slurm", "contact_in_group"},
		dao.CreateColumnContactInGroup,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	if err = r.updateGroupContactCount(ctx, tx, groupID); err != nil {
		return err
	}

	return nil
}

// checkExistContactInGroup
// return listExist -- list existing contactID
// return mapExist -- mapping contact ID how exist or not exist
func (r *Repository) checkExistContactInGroup(ctx context.Context, tx pgx.Tx, groupID uuid.UUID, contactIDs ...uuid.UUID) (listExist []uuid.UUID, mapExist map[uuid.UUID]bool, err error) {
	listExist = make([]uuid.UUID, 0)
	mapExist = make(map[uuid.UUID]bool)

	if len(contactIDs) == 0 {
		return listExist, mapExist, nil
	}

	query, args, err := r.genSQL.
		Select("contact_id").
		From("slurm.contact_in_group").
		Where(squirrel.Eq{"contact_id": contactIDs, "group_id": groupID}).ToSql()

	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var contactID = uuid.UUID{}

		if err = rows.Scan(&contactID); err != nil {
			return nil, nil, err
		}

		listExist = append(listExist, contactID)
		mapExist[contactID] = true
	}

	for _, contactID := range contactIDs {
		if _, ok := mapExist[contactID]; !ok {
			mapExist[contactID] = false
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return listExist, mapExist, nil
}
