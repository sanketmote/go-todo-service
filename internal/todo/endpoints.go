package todo

import (
	"context"

	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
	"github.com/sanketmote/go-todo-service/internal/service"
)

// MakeEndpoints creates Go-kit endpoints from the service.
func MakeEndpoints(svc service.TodoService) spectodo.Endpoints {
	return spectodo.Endpoints{
		Create:  makeCreateEndpoint(svc),
		List:    makeListEndpoint(svc),
		GetByID: makeGetByIDEndpoint(svc),
		Update:  makeUpdateEndpoint(svc),
		Delete:  makeDeleteEndpoint(svc),
	}
}

func makeCreateEndpoint(svc service.TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(spectodo.CreateTodoRequest)
		return svc.Create(ctx, req.Title, req.Description, req.DueDate, req.RepeatType)
	}
}

func makeListEndpoint(svc service.TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		return svc.List(ctx, req.IncludeCompleted, req.Sort, req.Limit, req.Offset)
	}
}

func makeGetByIDEndpoint(svc service.TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getByIDRequest)
		return svc.GetByID(ctx, req.ID)
	}
}

func makeUpdateEndpoint(svc service.TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRequest)
		return svc.Update(ctx, req.ID, req.Title, req.Description, req.Completed, req.DueDate, req.RepeatType)
	}
}

func makeDeleteEndpoint(svc service.TodoService) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteRequest)
		err := svc.Delete(ctx, req.ID)
		return nil, err
	}
}
