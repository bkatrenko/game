package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func Test_httpServer_handleJoin(t *testing.T) {
	type (
		fields struct {
			proc   func() Processor
			server *http.Server
		}

		args struct {
			w *httptest.ResponseRecorder
			r *http.Request
		}
	)

	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
	}{
		{
			name: "bad JSON input",
			fields: fields{
				proc: func() Processor {
					return NewProcessorMock(t)
				},
				server: &http.Server{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: mustMakeRequest(http.NewRequest(http.MethodPost, gameJoinRoute, bytes.NewBuffer([]byte("test")))),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "join error",
			fields: fields{
				proc: func() Processor {
					mockedProcessor := NewProcessorMock(t)
					mockedProcessor.On("Join", mock.Anything).Return(State{}, errors.New("Houston we have a problem"))

					return mockedProcessor
				},
				server: &http.Server{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: mustMakeRequest(http.NewRequest(http.MethodPost, gameJoinRoute, bytes.NewBuffer([]byte("{}")))),
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &httpServer{
				proc:   tt.fields.proc(),
				server: tt.fields.server,
			}

			s.handleJoin(tt.args.w, tt.args.r)
			assert.Equal(t, tt.expectedCode, tt.args.w.Code)
		})
	}
}

func mustMakeRequest(request *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}

	return request
}
