package cqs

import (
	"context"
	"github.com/mbict/go-commandbus/v2"
	"github.com/mbict/go-querybus"
	"github.com/mbict/webapp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type testQuery struct {
	Name string `query:"name"`
}

func testQueryHandler(ctx context.Context, query testQuery) (string, error) {
	return "hi " + query.Name, nil
}

func TestCqsApi_HandleQueryWithInterfaceHandler(t *testing.T) {
	qb := querybus.New()
	cb := commandbus.New()
	app := webapp.New()

	cqs := New(cb, qb, app)

	err := cqs.HandleQuery(http.MethodGet, "/query", testQueryHandler)

	assert.NoError(t, err)

	a, err := qb.Handle(context.Background(), testQuery{Name: "world"})

	assert.NoError(t, err)
	assert.Equal(t, "hi world", a)
}

func TestCqsApi_HandleQueryWithGenericWrap(t *testing.T) {
	qb := querybus.New()
	cb := commandbus.New()
	app := webapp.New()

	cqs := New(cb, qb, app)

	err := cqs.HandleQuery(http.MethodGet, "/query", WrapQueryHandler(testQueryHandler))

	assert.NoError(t, err)

	a, err := qb.Handle(context.Background(), testQuery{Name: "world"})

	assert.NoError(t, err)
	assert.Equal(t, "hi world", a)
}
