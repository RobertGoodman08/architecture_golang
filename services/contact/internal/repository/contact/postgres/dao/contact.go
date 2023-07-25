package dao

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"modified_at"`

	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`

	Name       string `db:"name"`
	Surname    string `db:"surname"`
	Patronymic string `db:"patronymic"`

	Age    uint64 `db:"age"`
	Gender uint8  `db:"gender"`
}

var CreateColumnContact = []string{
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
}

var CreateColumnContactInGroup = []string{
	"created_at",
	"modified_at",
	"group_id",
	"contact_id",
}
