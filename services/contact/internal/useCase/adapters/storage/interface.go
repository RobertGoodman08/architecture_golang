package storage

import (
	"github.com/google/uuid"

	"architecture_go/pkg/type/context"
	"architecture_go/pkg/type/queryParameter"
	"architecture_go/services/contact/internal/domain/contact"
	"architecture_go/services/contact/internal/domain/group"
)

type Storage interface {
	Contact
	Group
}

type Contact interface {
	CreateContact(ctx context.Context, contacts ...*contact.Contact) ([]*contact.Contact, error)
	UpdateContact(ctx context.Context, ID uuid.UUID, updateFn func(c *contact.Contact) (*contact.Contact, error)) (*contact.Contact, error)
	DeleteContact(ctx context.Context, ID uuid.UUID) error

	ContactReader
}

type ContactReader interface {
	ListContact(ctx context.Context, parameter queryParameter.QueryParameter) ([]*contact.Contact, error)
	ReadContactByID(ctx context.Context, ID uuid.UUID) (response *contact.Contact, err error)
	CountContact(ctx context.Context /*Тут можно передавать фильтр*/) (uint64, error)
}

type Group interface {
	CreateGroup(ctx context.Context, group *group.Group) (*group.Group, error)
	UpdateGroup(ctx context.Context, ID uuid.UUID, updateFn func(group *group.Group) (*group.Group, error)) (*group.Group, error)
	DeleteGroup(ctx context.Context, ID uuid.UUID /*Тут можно передавать фильтр*/) error

	GroupReader
	ContactInGroup
}

type GroupReader interface {
	ListGroup(ctx context.Context, parameter queryParameter.QueryParameter) ([]*group.Group, error)
	ReadGroupByID(ctx context.Context, ID uuid.UUID) (*group.Group, error)
	CountGroup(ctx context.Context /*Тут можно передавать фильтр*/) (uint64, error)
}

type ContactInGroup interface {
	CreateContactIntoGroup(ctx context.Context, groupID uuid.UUID, contacts ...*contact.Contact) ([]*contact.Contact, error)
	DeleteContactFromGroup(ctx context.Context, groupID, contactID uuid.UUID) error
	AddContactsToGroup(ctx context.Context, groupID uuid.UUID, contactIDs ...uuid.UUID) error
}
