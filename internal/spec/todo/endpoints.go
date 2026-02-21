package todo

import "github.com/go-kit/kit/endpoint"

// Endpoints holds one Go-kit Endpoint per API.
type Endpoints struct {
	Create   endpoint.Endpoint
	List     endpoint.Endpoint
	GetByID  endpoint.Endpoint
	Update   endpoint.Endpoint
	Delete   endpoint.Endpoint
}
