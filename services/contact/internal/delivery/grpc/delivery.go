package grpc

import (
	contact "architecture_go/services/contact/internal/delivery/grpc/interface"
	"architecture_go/services/contact/internal/useCase"
)

type Delivery struct {
	contact.UnimplementedContactServiceServer
	ucContact useCase.Contact
	ucGroup   useCase.Group

	options Options
}

type Options struct{}

func New(ucContact useCase.Contact, ucGroup useCase.Group, o Options) *Delivery {
	var d = &Delivery{
		ucContact: ucContact,
		ucGroup:   ucGroup,
	}

	d.SetOptions(o)
	return d
}

func (d *Delivery) SetOptions(options Options) {
	if d.options != options {
		d.options = options
	}
}
