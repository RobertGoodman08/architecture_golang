package group

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"robovoice/micro-services/_helpers/context"
)

type Group interface {
	UpdateGroupsContactCountByFilters(ctx context.Context, tx pgx.Tx, ID uuid.UUID) error
}
