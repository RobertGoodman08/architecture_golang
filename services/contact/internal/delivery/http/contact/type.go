package contact

import (
	"time"

	"architecture_go/pkg/type/email"
	"architecture_go/pkg/type/gender"
)

type ID struct {
	Value string `json:"id" uri:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type ContactResponse struct {
	// Идетификатор записи
	ID string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	// Дата создания контакта
	CreatedAt time.Time `json:"createdAt"  binding:"required"`
	// Дата последнего изменения контакта
	ModifiedAt time.Time `json:"modifiedAt"  binding:"required"`
	ShortContact
}

type ShortContact struct {
	// Мобильный телефон
	PhoneNumber string `json:"phoneNumber" binding:"required,max=50" maxLength:"50" example:"78002002020"`
	// Электронная почта
	Email email.Email `json:"email" binding:"omitempty,max=250,email" maxLength:"250" example:"example@gmail.com" format:"email" swaggertype:"string"`
	// Пол
	Gender gender.Gender `json:"gender" example:"1" enums:"1,2" swaggertype:"integer"`
	// Возраст
	Age uint8 `json:"age" binding:"min=0,max=200" minimum:"0" maximum:"200" default:"0" example:"42"`
	// Имя клиента
	Name string `json:"name" binding:"max=50" maxLength:"50" example:"Иван"`
	// Фамилия клиента
	Surname string `json:"surname" binding:"max=100" maxLength:"100" example:"Иванов"`
	// Отчество клиента
	Patronymic string `json:"patronymic" binding:"max=100" maxLength:"100" example:"Иванович"`
}

type ListContact struct {
	// Всего
	Total uint64 `json:"total" example:"10" default:"0" binding:"min=0" minimum:"0"`
	// Количество записей
	Limit uint64 `json:"limit"  example:"10" default:"10" binding:"min=0" minimum:"0"`
	// Смещение при получении записей
	Offset uint64 `json:"offset" example:"20" default:"0" binding:"min=0" minimum:"0"`

	List []*ContactResponse `json:"list"`
}
