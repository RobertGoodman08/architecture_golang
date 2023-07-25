package useCase

import (
	"github.com/google/uuid"

	"architecture_go/pkg/type/context"
	"architecture_go/pkg/type/queryParameter"
	"architecture_go/services/contact/internal/domain/contact"
	"architecture_go/services/contact/internal/domain/group"
)

type Contact interface {
	Create(c context.Context, contacts ...*contact.Contact) ([]*contact.Contact, error)
	Update(c context.Context, contactUpdate contact.Contact) (*contact.Contact, error)
	Delete(c context.Context, ID uuid.UUID /*Тут можно передавать фильтр*/) error

	ContactReader
}

type ContactReader interface {
	List(c context.Context, parameter queryParameter.QueryParameter) ([]*contact.Contact, error)
	ReadByID(c context.Context, ID uuid.UUID) (response *contact.Contact, err error)
	Count(c context.Context /*Тут можно передавать фильтр*/) (uint64, error)
}

type Group interface {
	Create(c context.Context, groupCreate *group.Group) (*group.Group, error)
	Update(c context.Context, groupUpdate *group.Group) (*group.Group, error)
	Delete(c context.Context, ID uuid.UUID /*Тут можно передавать фильтр*/) error

	GroupReader
	ContactInGroup
}

type GroupReader interface {
	List(c context.Context, parameter queryParameter.QueryParameter) ([]*group.Group, error)
	ReadByID(c context.Context, ID uuid.UUID) (*group.Group, error)
	Count(c context.Context /*Тут можно передавать фильтр*/) (uint64, error)
}

type ContactInGroup interface {
	CreateContactIntoGroup(c context.Context, groupID uuid.UUID, contacts ...*contact.Contact) ([]*contact.Contact, error)
	AddContactToGroup(c context.Context, groupID, contactID uuid.UUID) error
	DeleteContactFromGroup(c context.Context, groupID, contactID uuid.UUID) error
}
