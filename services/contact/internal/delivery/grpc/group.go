package grpc

import (
	"context"

	contact "architecture_go/services/contact/internal/delivery/grpc/interface"
)

func (d *Delivery) CreateGroup(ctx context.Context, request *contact.CreateGroupRequest) (*contact.CreateGroupResponse, error) {
	panic("implement me")
}

func (d *Delivery) UpdateGroup(ctx context.Context, request *contact.UpdateGroupRequest) (*contact.UpdateGroupResponse, error) {
	panic("implement me")
}

func (d *Delivery) DeleteGroup(ctx context.Context, request *contact.DeleteGroupRequest) (*contact.DeleteGroupResponse, error) {
	panic("implement me")
}
