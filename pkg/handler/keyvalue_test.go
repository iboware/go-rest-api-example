package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/iboware/go-rest-api-example/pkg/model"
)

func TestKeyValueHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		store          *KVStore
		request        *model.Tuple
		wantStatusCode int
		wantBody       *model.Tuple
	}{
		{
			name: "success",
			store: &KVStore{
				m:       make(map[string]string),
				RWMutex: &sync.RWMutex{},
			},
			request: &model.Tuple{
				Key:   "foo",
				Value: "bar",
			},
			wantStatusCode: http.StatusOK,
			wantBody: &model.Tuple{
				Key:   "foo",
				Value: "bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			b, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatal(err)
				return
			}
			req := httptest.NewRequest("POST", "/in-memory", bytes.NewReader(b))

			h := &KeyValueHandler{
				store: tt.store,
			}
			h.Create(rr, req)

			if rr.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, tt.wantStatusCode)
			}

			var respModel *model.Tuple
			if err := json.NewDecoder(rr.Result().Body).Decode(&respModel); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(respModel, tt.wantBody) {
				t.Errorf("Response body, %+v, did not match expected code %+v", respModel, tt.wantBody)
			}
		})
	}
}

func TestKeyValueHandler_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		store          *KVStore
		key            string
		wantStatusCode int
		wantBody       *model.Tuple
	}{
		{
			name: "success",
			store: &KVStore{
				m: map[string]string{
					"foo": "bar",
				},
				RWMutex: &sync.RWMutex{},
			},
			key:            "foo",
			wantStatusCode: http.StatusOK,
			wantBody: &model.Tuple{
				Key:   "foo",
				Value: "bar",
			},
		},
		{
			name: "miss",
			store: &KVStore{
				m: map[string]string{
					"foo": "bar",
				},
				RWMutex: &sync.RWMutex{},
			},
			key:            "bar",
			wantStatusCode: http.StatusNotFound,
			wantBody:       nil,
		},
		{
			name: "key query parameter not found",
			store: &KVStore{
				m: map[string]string{
					"foo": "bar",
				},
				RWMutex: &sync.RWMutex{},
			},
			key:            "",
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/in-memory", nil)

			q := req.URL.Query()
			if tt.key != "" {
				q.Add("key", tt.key)
			}
			req.URL.RawQuery = q.Encode()

			h := &KeyValueHandler{
				store: tt.store,
			}
			h.Fetch(rr, req)

			if rr.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, tt.wantStatusCode)
			}

			if rr.Result().StatusCode != http.StatusOK {
				return
			}

			var respModel *model.Tuple
			if err := json.NewDecoder(rr.Result().Body).Decode(&respModel); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(respModel, tt.wantBody) {
				t.Errorf("Response body, %+v, did not match expected code %+v", respModel, tt.wantBody)
			}
		})
	}
}
