package group

import "time"

type GroupResponse struct {
	// Идентификатор группы
	ID string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	// Дата создания группы
	CreatedAt time.Time `json:"createdAt"  binding:"required"`
	// Дата последнего изменения группы
	ModifiedAt time.Time `json:"modifiedAt"  binding:"required"`
	Group
}

// Group
// Описывает объект, который содержит информацию о группе.
type Group struct {
	ShortGroup
	// Кол-во контактов в группе
	ContactsAmount uint64 `json:"contactsAmount" default:"10" binding:"min=0" minimum:"0"`
}

type ShortGroup struct {
	// Название группы
	Name string `json:"name" binding:"required,max=100" example:"Название группы" maxLength:"100"`
	// Описание
	Description string `json:"description" example:"Описание группы" binding:"max=1000" maxLength:"1000"`
}

// GroupList
// Описывает объект, который содержит информацию о группе.
type GroupList struct {
	// Всего
	Total uint64 `json:"total" example:"10" default:"0" binding:"min=0" minimum:"0"`
	// Количество записей
	Limit uint64 `json:"limit"  example:"10" default:"10" binding:"min=0" minimum:"0"`
	// Смещение при получении записей
	Offset uint64 `json:"offset" example:"20" default:"0" binding:"min=0" minimum:"0"`
	// Список групп
	List []*GroupResponse `json:"list" binding:"min=0" minimum:"0"`
}

type ID struct {
	// Идентификатор группы
	Value string `json:"id" uri:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type ContactID struct {
	// Идентификатор контакта
	Value string `json:"id" uri:"contactId" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}
