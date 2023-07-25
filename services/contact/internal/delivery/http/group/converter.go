package group

import (
	"architecture_go/services/contact/internal/domain/group"
)

func ProtoToGroupResponse(response *group.Group) *GroupResponse {
	return &GroupResponse{
		ID:         response.ID().String(),
		CreatedAt:  response.CreatedAt(),
		ModifiedAt: response.ModifiedAt(),
		Group: Group{
			ShortGroup: ShortGroup{
				Name:        response.Name().Value(),
				Description: response.Description().Value(),
			},
			ContactsAmount: response.ContactCount(),
		},
	}
}
