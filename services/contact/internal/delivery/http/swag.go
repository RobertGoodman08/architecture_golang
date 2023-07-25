//go:build swag
// +build swag

package http

//go:generate swag init --parseDependency  --generalInfo delivery.go --output swagger/docs/
