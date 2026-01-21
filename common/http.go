package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpTest[T any] struct {
	RequestParams map[string]string
	RequestQuery  string
	RequestBody   string
	ResponseCode  int
	Response      *T
}

func (p HttpTest[T]) Run(t *testing.T, handler func(http.ResponseWriter, *http.Request)) {
	var body io.Reader
	if p.RequestBody != "" {
		body = bytes.NewBufferString(p.RequestBody)
	}

	req, err := http.NewRequest("", p.RequestQuery, body)
	assert.NoError(t, err)

	if p.RequestParams != nil {
		for k, v := range p.RequestParams {
			req.SetPathValue(k, v)
		}
	}

	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, p.ResponseCode, rec.Code)

	if p.Response != nil {
		var got *T
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, p.Response, got)
	}
}
