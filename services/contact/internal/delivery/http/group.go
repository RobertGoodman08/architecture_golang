package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"architecture_go/pkg/tools/converter"
	"architecture_go/pkg/type/context"
	"architecture_go/pkg/type/pagination"
	"architecture_go/pkg/type/query"
	"architecture_go/pkg/type/queryParameter"
	jsonGroup "architecture_go/services/contact/internal/delivery/http/group"
	domainGroup "architecture_go/services/contact/internal/domain/group"
	"architecture_go/services/contact/internal/domain/group/description"
	"architecture_go/services/contact/internal/domain/group/name"
	"architecture_go/services/contact/internal/useCase"
)

var mappingSortsGroup = query.SortsOptions{
	"id":           {},
	"name":         {},
	"description":  {},
	"contactCount": {},
}

// CreateGroup
// @Summary Метод позволяет создать группу контактов.
// @Description Метод позволяет создать группу контактов.
// @Tags 	groups
// @Accept  json
// @Produce json
// @Param   group 		body 		jsonGroup.ShortGroup 	true	"Данные по группе"
// @Success 200			{object}  	jsonGroup.GroupResponse	true
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse					"404 Not Found"
// @Router /groups/ [post]
func (d *Delivery) CreateGroup(c *gin.Context) {

	var ctx = context.New(c)

	var group = &jsonGroup.ShortGroup{}

	if err := c.ShouldBindJSON(&group); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	groupName, err := name.New(group.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}
	groupDescription, err := description.New(group.Description)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}
	newGroup, err := d.ucGroup.Create(ctx, domainGroup.New(
		groupName,
		groupDescription,
	))
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonGroup.GroupResponse{
		ID:         newGroup.ID().String(),
		CreatedAt:  newGroup.CreatedAt(),
		ModifiedAt: newGroup.ModifiedAt(),
		Group: jsonGroup.Group{
			ShortGroup: jsonGroup.ShortGroup{
				Name:        newGroup.Name().Value(),
				Description: newGroup.Description().Value(),
			},
			ContactsAmount: newGroup.ContactCount(),
		},
	})
}

// UpdateGroup
// @Summary Метод позволяет обновить данные группы.
// @Description Метод позволяет обновить данные группы.
// @Tags 	groups
// @Accept  json
// @Produce json
// @Param   id 			path 		string 					true	"Идентификатор группы"
// @Param   group 		body 		jsonGroup.ShortGroup 	true	"Данные по группе"
// @Success 200			{object}  	jsonGroup.GroupResponse
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse					"404 Not Found"
// @Router /groups/{id} [put]
func (d *Delivery) UpdateGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	group := jsonGroup.ShortGroup{}
	if err := c.ShouldBindJSON(&group); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	groupName, err := name.New(group.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}
	groupDescription, err := description.New(group.Description)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucGroup.Update(ctx, domainGroup.NewWithID(
		converter.StringToUUID(id.Value),
		time.Now().UTC(),
		time.Now().UTC(),
		groupName,
		groupDescription,
		0,
	))
	if err != nil {
		if errors.Is(err, useCase.ErrGroupNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, jsonGroup.ProtoToGroupResponse(response))
}

// DeleteGroup
// @Summary Метод позволяет удалить группу.
// @Description Метод позволяет удалить группу.
// @Tags 	groups
// @Accept  json
// @Produce json
// @Param   id 			path 		string 			true 	"Идентификатор группы"
// @Success 200			{object}  	string
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Router /groups/{id} [delete]
func (d *Delivery) DeleteGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	if err := d.ucGroup.Delete(ctx, converter.StringToUUID(id.Value)); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// ListGroup
// @Summary Метод позволяет получить список групп.
// @Description Метод позволяет получить список групп.
// @Tags 	groups
// @Accept  json
// @Produce json
// @Param 	limit 		query 		int 					false "Количество записей" default(10) mininum(0) maxinum(100)
// @Param 	offset 		query 		int 					false "Смещение при получении записей" default(0) mininum(0)
// @Param 	sort 		query 		string 					false "Сортировка по полю" default(name)
// @Success 200			{object}  	jsonGroup.GroupList
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Router /groups/ [get]
func (d *Delivery) ListGroup(c *gin.Context) {

	var ctx = context.New(c)

	params, err := query.ParseQuery(c, query.Options{
		Sorts: mappingSortsGroup,
	})

	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	groups, err := d.ucGroup.List(ctx, queryParameter.QueryParameter{
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

	var list = make([]*jsonGroup.GroupResponse, len(groups))

	for i, elem := range groups {
		list[i] = jsonGroup.ProtoToGroupResponse(elem)
	}

	c.JSON(http.StatusOK, jsonGroup.GroupList{
		Total:  count,
		Limit:  params.Limit,
		Offset: params.Offset,
		List:   list,
	})
}

// ReadGroupByID
// @Summary Метод позволяет получить данные по группе.
// @Description Метод позволяет получить данные по группе.
// @Tags 	groups
// @Accept  json
// @Produce json
// @Param   id 			path 		string 					true 	"Идентификатор группы контактов"
// @Success 200			{object}  	jsonGroup.GroupResponse
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse					"404 Not Found"
// @Router /groups/{id} [get]
func (d *Delivery) ReadGroupByID(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucGroup.ReadByID(ctx, converter.StringToUUID(id.Value))
	if err != nil {
		if errors.Is(err, useCase.ErrGroupNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonGroup.ProtoToGroupResponse(response))
}
