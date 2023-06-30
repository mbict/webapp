package cqs

import (
	"context"
	"github.com/mbict/go-commandbus/v2"
	"github.com/mbict/go-querybus"
	"github.com/mbict/webapp"
	"net/http"
	"reflect"
)

type CommandInitializer interface {
	Init()
}

type QueryHandler = querybus.QueryHandler

type CommandHandler = commandbus.CommandHandler

type chc struct {
	constructor func() any
	handler     any
}

type commandHandler[C any] func(ctx context.Context, cmd C) error

func WrapCommandHandler[C any](handler commandHandler[C]) *chc {

	if reflect.TypeOf(handler).In(1).Kind() != reflect.Struct {
		panic("invalid command type in handler, please provide a (non pointer) structure type for the 2nd argument of the handler")
	}

	return &chc{
		constructor: func() any {
			return *new(C)
		},
		handler: handler,
	}
}

type qhc struct {
	constructor     func() any
	originalHandler any
	handler         QueryHandler
}

type queryHandler[Q any, R any] func(ctx context.Context, query Q) (R, error)

func WrapQueryHandler[Q any, R any](handler queryHandler[Q, R]) *qhc {

	if reflect.TypeOf(handler).In(1).Kind() != reflect.Struct {
		panic("invalid command type in handler, please provide a (non pointer) structure type for the 2nd argument of the handler")
	}

	return &qhc{
		constructor: func() any {
			return *new(Q)
		},
		originalHandler: handler,
		handler: querybus.QueryHandlerFunc(func(ctx context.Context, q any) (any, error) {
			return handler(ctx, q.(Q))
		}),
	}
}

type SuccessHandler func(context webapp.Context, cmd any) error

// CQSApi is a Command Query Segregation helper that ties handlers to an api endpoint
type CQSApi interface {
	HandleQuery(method, path string, queryHandler any, middleware ...webapp.MiddlewareFunc) error
	HandleCommand(method, path string, commandHandler any, middleware ...webapp.MiddlewareFunc) error
	HandleCommandWithSuccessResponseHandler(method, path string, commandHandler any, successHandler SuccessHandler, middleware ...webapp.MiddlewareFunc) error
}

func New(cb commandbus.CommandBus, qb querybus.QueryBus, app webapp.WebApp) CQSApi {
	return &cqsApi{
		app:        app,
		commandBus: cb,
		queryBus:   qb,
	}
}

type cqsApi struct {
	app        webapp.WebApp
	commandBus commandbus.CommandBus
	queryBus   querybus.QueryBus
}

func (cq *cqsApi) HandleQuery(method, path string, queryHandler any, middleware ...webapp.MiddlewareFunc) error {
	var newQuery func() any
	var handler QueryHandler

	qc, isQueryHandlerConstructor := queryHandler.(*qhc)

	//when the provided query handler is a qhc (Query handler constructor), we handle the creation of the command via native new func
	if isQueryHandlerConstructor {
		newQuery = qc.constructor
		queryHandler = qc.originalHandler
		handler = qc.handler
	}

	// extract the type of the query
	t := reflect.TypeOf(queryHandler)
	tQry := t.In(1)
	isPtr := tQry.Kind() == reflect.Ptr

	//when the provided query handler is a QueryHandler type, we handle the creation of the command via reflection
	if !isQueryHandlerConstructor {
		if isPtr == true {
			tQry = tQry.Elem()
		}
		newQuery = func() any {
			return reflect.New(tQry).Elem().Interface()
		}

		handler = querybus.WrapHandler(queryHandler)

	}

	//register handler on the querybus for this query with a dummy
	dQry := reflect.New(tQry)
	if isPtr == true {
		dQry.Elem()
	}
	if err := cq.queryBus.Register(dQry.Interface(), handler); err != nil {
		return err
	}

	httpHandler := func(c webapp.Context) error {
		query := newQuery()

		// when command implements the initializer interface, we call it here
		if qi, ok := query.(CommandInitializer); ok {
			qi.Init()
		}

		//bind the query to get it filled in
		if err := c.Bind(&query); err != nil {
			return err
		}

		//validate the request
		if err := c.Validate(query); err != nil {
			return err
		}

		//if the query bus expects a ptr query we get th addr here
		if isPtr {
			query = &query
		}

		//dispatch over the query bus and wait for a result
		res, err := cq.queryBus.Handle(c.Request().Context(), query)

		if err != nil {
			return err
		}

		//return the result
		return c.JSON(http.StatusOK, res)
	}

	cq.app.Add(method, path, httpHandler, middleware...)
	return nil
}

func (cq *cqsApi) HandleCommand(method, path string, commandHandler any, middleware ...webapp.MiddlewareFunc) error {
	return cq.HandleCommandWithSuccessResponseHandler(method, path, commandHandler, nil, middleware...)
}

func (cq *cqsApi) HandleCommandWithSuccessResponseHandler(method, path string, commandHandler any, successHandler SuccessHandler, middleware ...webapp.MiddlewareFunc) error {
	var newCommand func() any

	// extract the type of the query
	t := reflect.TypeOf(commandHandler)
	tCmd := t.In(1)
	isPtr := tCmd.Kind() == reflect.Ptr

	//when the provided query handler is a CommandHandler type, we handle the creation of the command via reflection
	handler, isCommandHandler := commandHandler.(CommandHandler)
	if isCommandHandler {
		if isPtr == true {
			tCmd = tCmd.Elem()
		}

		newCommand = func() any {
			return reflect.New(tCmd).Elem().Interface()
		}
	}

	//when the provided query handler is a qhc (Query handler constructor), we handle the creation of the command via native new func
	cc, isCommandHandlerConstructor := commandHandler.(chc)
	if isCommandHandlerConstructor {
		newCommand = cc.constructor
		//handler = cc.handler
	}

	//register handler on the commandbus for this command with a dummy
	dCmd := reflect.New(tCmd)
	if isPtr == true {
		dCmd.Elem()
	}
	if err := cq.commandBus.Register(dCmd.Interface(), handler); err != nil {
		return err
	}

	httpHandler := func(c webapp.Context) error {
		cmd := newCommand()

		// when command implements the initializer interface, we call it here
		if ci, ok := cmd.(CommandInitializer); ok {
			ci.Init()
		}

		//bind the command to get it filled in
		if err := c.Bind(c); err != nil {
			return err
		}

		//validate the request
		if err := c.Validate(c); err != nil {
			return err
		}

		//dispatch over the commandbus and wait for a result
		if err := cq.commandBus.Handle(c.Request().Context(), cmd); err != nil {
			return err
		}

		if successHandler != nil {
			return successHandler(c, cmd)
		}

		//return the result
		return c.NoContent()
	}

	cq.app.Add(method, path, httpHandler, middleware...)
	return nil
}
