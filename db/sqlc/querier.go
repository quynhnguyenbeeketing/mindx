package db

import "context"

type Querier interface {
	GetUsers(ctx context.Context, req *Users) ([]*Users, error)
	CreateUsers(ctx context.Context, req *Users) (*Users, error)
	CreateListUsers(ctx context.Context, req []*Users) error
	UpdateUsers(ctx context.Context, req *Users) error

	CreateLocations(ctx context.Context, req *Locations) (*Locations, error)
	GetLocations(ctx context.Context, req *Locations) ([]*Locations, error)

	CreateUsersLocations(ctx context.Context, req *UsersLocations) error
	GetUsersLocations(ctx context.Context, req *UsersLocations) ([]*UsersLocations, error)
}

var _ Querier = (*Queries)(nil)
