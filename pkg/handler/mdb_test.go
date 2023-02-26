//go:build unit

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iboware/go-rest-api-example/pkg/model"
	"github.com/iboware/go-rest-api-example/pkg/store/mocks"
)

func TestMDBHandler_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		request        *model.MDBRequest
		mockStore      func(f *mocks.MockStore)
		wantStatusCode int
		wantBody       *model.MDBResponse
	}{
		{
			name: "invalid start date",
			request: &model.MDBRequest{
				StartDate: "201901-01",
				EndDate:   "2019-01-01",
				MinCount:  90,
				MaxCount:  100,
			},
			mockStore: func(f *mocks.MockStore) {
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: &model.MDBResponse{
				Code:    -1,
				Msg:     "invalid format for startdate",
				Records: []model.Record{},
			},
		},
		{
			name: "invalid end date",
			request: &model.MDBRequest{
				StartDate: "2019-01-01",
				EndDate:   "2019-0101",
				MinCount:  90,
				MaxCount:  100,
			},
			mockStore: func(f *mocks.MockStore) {
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: &model.MDBResponse{
				Code:    -1,
				Msg:     "invalid format for enddate",
				Records: []model.Record{},
			},
		},
		{
			name: "error from database",
			request: &model.MDBRequest{
				StartDate: "2019-01-01",
				EndDate:   "2019-01-01",
				MinCount:  90,
				MaxCount:  100,
			},
			mockStore: func(f *mocks.MockStore) {
				f.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: &model.MDBResponse{
				Code:    -1,
				Msg:     "could not retrieve the records from the database",
				Records: []model.Record{},
			},
		},
		{
			name: "success",
			request: &model.MDBRequest{
				StartDate: "2019-01-01",
				EndDate:   "2019-01-01",
				MinCount:  90,
				MaxCount:  100,
			},
			mockStore: func(f *mocks.MockStore) {
				f.EXPECT().Find(gomock.Any(), gomock.Any()).Return([]model.Record{
					{
						Key:        "foo1",
						CreatedAt:  "bar1",
						TotalCount: 100,
					},
					{
						Key:        "foo2",
						CreatedAt:  "bar2",
						TotalCount: 200,
					},
				}, nil)
			},
			wantStatusCode: http.StatusOK,
			wantBody: &model.MDBResponse{
				Code: 0,
				Msg:  "success",
				Records: []model.Record{
					{
						Key:        "foo1",
						CreatedAt:  "bar1",
						TotalCount: 100,
					},
					{
						Key:        "foo2",
						CreatedAt:  "bar2",
						TotalCount: 200,
					},
				},
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
			req := httptest.NewRequest("POST", "/mdb", bytes.NewReader(b))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ms := mocks.NewMockStore(ctrl)

			if tt.mockStore != nil {
				tt.mockStore(ms)
			}

			h := &MDBHandler{
				ctx:   context.Background(),
				store: ms,
			}
			h.Fetch(rr, req)

			if rr.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, tt.wantStatusCode)
			}

			var respModel *model.MDBResponse
			if err := json.NewDecoder(rr.Result().Body).Decode(&respModel); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(respModel, tt.wantBody) {
				t.Errorf("Response body, %+v, did not match expected code %+v", respModel, tt.wantBody)
			}

		})
	}
}
