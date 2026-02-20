package todo

import (
	"context"

	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
)

// MakeEndpoints creates Go-kit endpoints from the service.
func MakeEndpoints(svc TodoService) spectodo.Endpoints {
	return spectodo.Endpoints{
		Create:  makeCreateEndpoint(svc),
		List:    makeListEndpoint(svc),
		GetByID: makeGetByIDEndpoint(svc),
		Update:  makeUpdateEndpoint(svc),
		Delete:  makeDeleteEndpoint(svc),
	}
}

func makeCreateEndpoint(svc TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(spectodo.CreateTodoRequest)
		return svc.Create(ctx, req.Title, req.Description, req.DueDate, req.RepeatType)
	}
}

func makeListEndpoint(svc TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		return svc.List(ctx, req.IncludeCompleted, req.Sort, req.Limit, req.Offset)
	}
}

func makeGetByIDEndpoint(svc TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getByIDRequest)
		return svc.GetByID(ctx, req.ID)
	}
}

func makeUpdateEndpoint(svc TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRequest)
		return svc.Update(ctx, req.ID, req.Title, req.Description, req.Completed, req.DueDate, req.RepeatType)
	}
}

func makeDeleteEndpoint(svc TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteRequest)
		err := svc.Delete(ctx, req.ID)
		return nil, err
	}
}
