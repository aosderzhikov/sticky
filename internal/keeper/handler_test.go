package keeper

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetHandle(t *testing.T) {
	cases := []struct {
		name        string
		serviceFunc func(t *testing.T) Service
		reqFunc     func(t *testing.T) *http.Request
		wantFunc    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "get empty key query",
			reqFunc: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://test", http.NoBody)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				return NewMockService(ctrl)
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
				require.Contains(t, rec.Body.String(), "key")
			},
		},
		{
			name: "get success",
			reqFunc: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://test?key=key1", http.NoBody)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				service := NewMockService(ctrl)
				service.EXPECT().Get("key1").Return([]byte("data"))
				return service
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Result().StatusCode)
				require.Equal(t, rec.Body.String(), "data")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewHandler(c.serviceFunc(t))
			rec := httptest.NewRecorder()
			h.GetHandle(rec, c.reqFunc(t))
			c.wantFunc(t, rec)
		})
	}
}

func TestSetHandle(t *testing.T) {
	cases := []struct {
		name        string
		serviceFunc func(t *testing.T) Service
		reqFunc     func(t *testing.T) *http.Request
		wantFunc    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "set empty key",
			reqFunc: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://test", http.NoBody)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				return NewMockService(ctrl)
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
				require.Contains(t, rec.Body.String(), "key")
			},
		},
		{
			name: "set success",
			reqFunc: func(t *testing.T) *http.Request {
				body := bytes.NewReader([]byte("data"))
				req, err := http.NewRequest(http.MethodGet, "http://test?key=key1&ttl=5s", body)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				service := NewMockService(ctrl)
				service.EXPECT().Set("key1", []byte("data"), 5*time.Second)
				return service
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Result().StatusCode)
			},
		},
		{
			name: "invalid ttl query param",
			reqFunc: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://test?key=a&ttl=a", http.NoBody)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				return NewMockService(ctrl)
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
				require.Contains(t, rec.Body.String(), "ttl")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewHandler(c.serviceFunc(t))
			rec := httptest.NewRecorder()
			h.SetHandle(rec, c.reqFunc(t))
			c.wantFunc(t, rec)
		})
	}
}

func TestDeleteHandle(t *testing.T) {
	cases := []struct {
		name        string
		serviceFunc func(t *testing.T) Service
		reqFunc     func(t *testing.T) *http.Request
		wantFunc    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "get empty key query",
			reqFunc: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://test", http.NoBody)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				return NewMockService(ctrl)
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
				require.Contains(t, rec.Body.String(), "key")
			},
		},
		{
			name: "delete success",
			reqFunc: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "http://test?key=key1", http.NoBody)
				require.NoError(t, err)
				return req
			},
			serviceFunc: func(t *testing.T) Service {
				ctrl := gomock.NewController(t)
				service := NewMockService(ctrl)
				service.EXPECT().Delete("key1").Return()
				return service
			},
			wantFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Result().StatusCode)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewHandler(c.serviceFunc(t))
			rec := httptest.NewRecorder()
			h.DeleteHandle(rec, c.reqFunc(t))
			c.wantFunc(t, rec)
		})
	}
}
