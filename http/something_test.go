package http

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"app-go/common"
	"app-go/http/api"
	"app-go/service/mock"
)

func TestRead(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		title         string
		somethingFunc func(context.Context, string) (string, error)
		common.HttpTest[api.Response]
	}{
		{
			title: "success",
			somethingFunc: func(ctx context.Context, str string) (string, error) {
				assert.Equal(t, "request", str)
				return str, nil
			},
			HttpTest: common.HttpTest[api.Response]{
				RequestQuery: "?query=request",
				ResponseCode: 200,
				Response: &api.Response{
					Data: "request",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			srv := &mock.AppMock{}
			srv.ReadFunc = tc.somethingFunc

			s := NewServer(8080, srv)
			tc.HttpTest.Run(t, s.read)
		})
	}

}

func TestSomething(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		title         string
		somethingFunc func(context.Context, string) (string, error)
		common.HttpTest[api.Response]
	}{
		{
			title: "success",
			somethingFunc: func(ctx context.Context, str string) (string, error) {
				assert.Equal(t, "request", str)
				return str, nil
			},
			HttpTest: common.HttpTest[api.Response]{
				RequestBody: `
				{
					"field": "request"
				}
				`,
				ResponseCode: 200,
				Response: &api.Response{
					Data: "request",
				},
			},
		},
		{
			title: "bad request",
			HttpTest: common.HttpTest[api.Response]{
				RequestBody: `
				{
					"field-unknown": "something"
				}
				`,
				ResponseCode: 400,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			srv := &mock.AppMock{}
			srv.SomethingFunc = tc.somethingFunc

			s := NewServer(8080, srv)
			tc.HttpTest.Run(t, s.something)
		})
	}

}

func TestSomethingId(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		title         string
		somethingFunc func(context.Context, string, string) (string, error)
		common.HttpTest[api.Response]
	}{
		{
			title: "success",
			somethingFunc: func(ctx context.Context, str string, id string) (string, error) {
				assert.Equal(t, "12345", id)
				assert.Equal(t, "request", str)
				return str, nil
			},
			HttpTest: common.HttpTest[api.Response]{
				RequestParams: map[string]string{
					"id": "12345",
				},
				RequestBody: `
				{
					"field": "request"
				}
				`,
				ResponseCode: 200,
				Response: &api.Response{
					Data: "request",
				},
			},
		},
		{
			title: "bad request",
			HttpTest: common.HttpTest[api.Response]{
				RequestBody: `
				{
					"field-unknown": "something"
				}
				`,
				ResponseCode: 400,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			srv := &mock.AppMock{}
			srv.SomethingIdFunc = tc.somethingFunc

			s := NewServer(8080, srv)
			tc.HttpTest.Run(t, s.somethingId)
		})
	}

}
