package contact

import (
	domainContact "architecture_go/services/contact/internal/domain/contact"
)

func ToContactResponse(response *domainContact.Contact) *ContactResponse {
	return &ContactResponse{
		ID:         response.ID().String(),
		CreatedAt:  response.CreatedAt(),
		ModifiedAt: response.ModifiedAt(),
		ShortContact: ShortContact{
			PhoneNumber: response.PhoneNumber().String(),
			Email:       response.Email(),
			Gender:      response.Gender(),
			Age:         uint8(response.Age()),
			Name:        response.Name().String(),
			Surname:     response.Surname().String(),
			Patronymic:  response.Patronymic().String(),
		},
	}
}
