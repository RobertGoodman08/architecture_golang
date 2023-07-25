package postgres

import (
	"github.com/jackc/pgx/v4"

	"architecture_go/pkg/type/email"
	"architecture_go/pkg/type/gender"
	"architecture_go/pkg/type/phoneNumber"
	"architecture_go/services/contact/internal/domain/contact"
	"architecture_go/services/contact/internal/domain/contact/age"
	"architecture_go/services/contact/internal/domain/contact/name"
	"architecture_go/services/contact/internal/domain/contact/patronymic"
	"architecture_go/services/contact/internal/domain/contact/surname"
	"architecture_go/services/contact/internal/repository/storage/postgres/dao"
)

func (r Repository) toCopyFromSource(contacts ...*contact.Contact) pgx.CopyFromSource {
	rows := make([][]interface{}, len(contacts))

	for i, val := range contacts {
		rows[i] = []interface{}{
			val.ID(),
			val.CreatedAt(),
			val.ModifiedAt(),
			val.PhoneNumber().String(),
			val.Email().String(),
			val.Name().String(),
			val.Surname().String(),
			val.Patronymic().String(),
			val.Age(),
			val.Gender(),
		}
	}
	return pgx.CopyFromRows(rows)
}

func (r Repository) toDomainContact(dao *dao.Contact) (*contact.Contact, error) {

	nameObject, err := name.New(dao.Name)
	if err != nil {
		return nil, err
	}

	surnameObject, err := surname.New(dao.Surname)
	if err != nil {
		return nil, err
	}

	patronymicObject, err := patronymic.New(dao.Patronymic)
	if err != nil {
		return nil, err
	}

	ageObject, err := age.New(dao.Age)
	if err != nil {
		return nil, err
	}

	localEmail, err := email.New(dao.Email)
	if err != nil {
		return nil, err
	}

	result, err := contact.NewWithID(
		dao.ID,
		dao.CreatedAt,
		dao.ModifiedAt,
		*phoneNumber.New(dao.PhoneNumber),
		localEmail,
		*nameObject,
		*surnameObject,
		*patronymicObject,
		*ageObject,
		gender.Gender(dao.Gender),
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r Repository) toDomainContacts(dao []*dao.Contact) ([]*contact.Contact, error) {
	var result = make([]*contact.Contact, len(dao))
	for i, v := range dao {
		c, err := r.toDomainContact(v)
		if err != nil {
			return nil, err
		}
		result[i] = c
	}
	return result, nil
}
