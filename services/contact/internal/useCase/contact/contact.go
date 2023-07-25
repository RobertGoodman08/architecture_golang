package contact

import (
	"time"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"

	"architecture_go/pkg/type/context"
	"architecture_go/pkg/type/queryParameter"
	"architecture_go/services/contact/internal/domain/contact"
)

func (uc *UseCase) Create(ctx context.Context, contacts ...*contact.Contact) ([]*contact.Contact, error) {
	return uc.adapterStorage.CreateContact(ctx, contacts...)
}

func (uc *UseCase) Update(ctx context.Context, contactUpdate contact.Contact) (*contact.Contact, error) {
	return uc.adapterStorage.UpdateContact(ctx, contactUpdate.ID(), func(oldContact *contact.Contact) (*contact.Contact, error) {
		return contact.NewWithID(
			oldContact.ID(),
			oldContact.CreatedAt(),
			time.Now().UTC(),
			contactUpdate.PhoneNumber(),
			contactUpdate.Email(),
			contactUpdate.Name(),
			contactUpdate.Surname(),
			contactUpdate.Patronymic(),
			contactUpdate.Age(),
			contactUpdate.Gender(),
		)
	})
}

func (uc *UseCase) Delete(ctx context.Context, ID uuid.UUID) error {
	return uc.adapterStorage.DeleteContact(ctx, ID)
}

func (uc *UseCase) List(c context.Context, parameter queryParameter.QueryParameter) ([]*contact.Contact, error) {

	span, ctx := opentracing.StartSpanFromContext(c, "List")
	defer span.Finish()

	return uc.adapterStorage.ListContact(context.New(ctx), parameter)
}

func (uc *UseCase) ReadByID(ctx context.Context, ID uuid.UUID) (response *contact.Contact, err error) {

	return uc.adapterStorage.ReadContactByID(ctx, ID)
}

func (uc *UseCase) Count(c context.Context) (uint64, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "Count")
	defer span.Finish()

	return uc.adapterStorage.CountContact(context.New(ctx))
}
