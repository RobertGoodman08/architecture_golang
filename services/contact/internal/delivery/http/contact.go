package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"architecture_go/pkg/tools/converter"
	"architecture_go/pkg/type/context"
	"architecture_go/pkg/type/logger"
	"architecture_go/pkg/type/pagination"
	"architecture_go/pkg/type/phoneNumber"
	"architecture_go/pkg/type/query"
	"architecture_go/pkg/type/queryParameter"
	jsonContact "architecture_go/services/contact/internal/delivery/http/contact"
	domainContact "architecture_go/services/contact/internal/domain/contact"
	"architecture_go/services/contact/internal/domain/contact/age"
	"architecture_go/services/contact/internal/domain/contact/name"
	"architecture_go/services/contact/internal/domain/contact/patronymic"
	"architecture_go/services/contact/internal/domain/contact/surname"
	"architecture_go/services/contact/internal/useCase"
)

var mappingSortsContact = query.SortsOptions{
	"name":        {},
	"surname":     {},
	"patronymic":  {},
	"phoneNumber": {},
	"email":       {},
	"gender":      {},
	"age":         {},
}

// CreateContact
// @Summary Метод позволяет создать контакт.
// @Description Метод позволяет создать контакт.
// @Tags contacts
// @Accept  json
// @Produce json
// @Param   contact 	body 		jsonContact.ShortContact 		    true  "Данные по контакту"
// @Success 201			{object}  	jsonContact.ContactResponse 		true  "Структура контакта"
// @Success 200
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Router /contacts/ [post]
func (d *Delivery) CreateContact(c *gin.Context) {

	var ctx = context.New(c)

	contact := jsonContact.ShortContact{}
	if err := c.ShouldBindJSON(&contact); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactAge, err := age.New(contact.Age)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactName, err := name.New(contact.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactSurname, err := surname.New(contact.Surname)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactPatronymic, err := patronymic.New(contact.Patronymic)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	dContact, err := domainContact.New(
		*phoneNumber.New(contact.PhoneNumber),
		contact.Email,
		*contactName,
		*contactSurname,
		*contactPatronymic,
		*contactAge,
		contact.Gender,
	)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucContact.Create(ctx, dContact)
	if err != nil {

		SetError(c, http.StatusInternalServerError, err)
		return
	}
	logger.InfoWithContext(ctx, "test log", zap.Any("Test", "Test"))
	if len(response) > 0 {
		c.JSON(http.StatusCreated, jsonContact.ToContactResponse(response[0]))
	} else {
		c.Status(http.StatusOK)
	}
}

// UpdateContact
// @Summary Метод позволяет обновить данные контакта.
// @Description Метод позволяет обновить данные контакта.
// @Tags contacts
// @Accept  json
// @Produce json
// @Param   id 			path 		string 						true  "Идентификатор контакта"
// @Param   contact 	body 		jsonContact.ShortContact	true  "Данные по контакту"
// @Success 200			{object}  	jsonContact.ContactResponse true  "Структура контакта"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			  		  "404 Not Found"
// @Router /contacts/{id} [put]
func (d *Delivery) UpdateContact(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonContact.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contact := jsonContact.ShortContact{}
	if err := c.ShouldBindJSON(&contact); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	contactAge, err := age.New(contact.Age)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactName, err := name.New(contact.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactSurname, err := surname.New(contact.Surname)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactPatronymic, err := patronymic.New(contact.Patronymic)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	var dContact, _ = domainContact.NewWithID(
		converter.StringToUUID(id.Value),
		time.Now().UTC(),
		time.Now().UTC(),
		*phoneNumber.New(contact.PhoneNumber),
		contact.Email,
		*contactName,
		*contactSurname,
		*contactPatronymic,
		*contactAge,
		contact.Gender,
	)

	response, err := d.ucContact.Update(ctx, *dContact)
	if err != nil {
		if errors.Is(err, useCase.ErrContactNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonContact.ToContactResponse(response))

}

// DeleteContact
// @Summary Метод позволяет удалить контакт.
// @Description Метод позволяет удалить контакт.
// @Tags contacts
// @Accept  json
// @Produce json
// @Param   id 			path 		string 			true 	"Идентификатор контакта"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Router /contacts/{id} [delete]
func (d *Delivery) DeleteContact(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonContact.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	if err := d.ucContact.Delete(ctx, converter.StringToUUID(id.Value)); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// ListContact
// @Summary Получить список контактов.
// @Description Метод позволяет получить список контактов.
// @Tags contacts
// @Accept  json
// @Produce json
// @Param 	limit 		query 		int 					false "Количество записей" default(10) mininum(0) maxinum(100)
// @Param 	offset 		query 		int 					false "Смещение при получении записей" default(0) mininum(0)
// @Param 	sort 		query 		string 					false "Сортировка по полю" default(name)
// @Success 200			{object}  	jsonContact.ListContact true  "Список контактов"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Router /contacts/ [get]
func (d *Delivery) ListContact(c *gin.Context) {

	var ctx = context.New(c)
	params, err := query.ParseQuery(c, query.Options{
		Sorts: mappingSortsContact,
	})

	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contacts, err := d.ucContact.List(ctx, queryParameter.QueryParameter{
		Sorts: params.Sorts,
		Pagination: pagination.Pagination{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	})
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	count, err := d.ucContact.Count(ctx)
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	var result = jsonContact.ListContact{
		Total:  count,
		Limit:  params.Limit,
		Offset: params.Offset,
		List:   []*jsonContact.ContactResponse{},
	}
	for _, value := range contacts {
		result.List = append(result.List, jsonContact.ToContactResponse(value))
	}

	c.JSON(http.StatusOK, result)
}

// ReadContactByID
// @Summary Получить контакт.
// @Description Метод позволяет получить контакт по мдентификатору контакта.
// @Tags contacts
// @Accept  json
// @Produce json
// @Param   id 			path 		string 						true "Идентификатор контакта"
// @Success 200			{object}  	jsonContact.ContactResponse true "Структура контакта"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse					  "404 Not Found"
// @Router /contacts/{id} [get]
func (d *Delivery) ReadContactByID(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonContact.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucContact.ReadByID(ctx, converter.StringToUUID(id.Value))
	if err != nil {
		if errors.Is(err, useCase.ErrContactNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonContact.ToContactResponse(response))

}
