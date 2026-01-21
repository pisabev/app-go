package service

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoSomething(t *testing.T) {
	log.SetOutput(io.Discard)
	
	tests := []struct {
		name   string
		input  string
		result string
		error  string
	}{
		{name: "test1", input: "hello world", result: "hello-world"},
		{name: "test2", input: "go is awesome", result: "go-is-awesome"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewApp()
			res, err := srv.Something(context.Background(), tt.input)
			if tt.error != "" {
				assert.ErrorContains(t, err, tt.error)
				return
			}
			assert.Equal(t, tt.result, res)
		})
	}
}
