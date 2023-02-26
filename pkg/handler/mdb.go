package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/iboware/go-rest-api-example/pkg/model"
	"github.com/iboware/go-rest-api-example/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
)

type MDBHandler struct {
	ctx   context.Context
	store store.Store
}

var dateRegexp = `^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`

func NewMDBHandler(ctx context.Context, store store.Store) *MDBHandler {
	return &MDBHandler{
		ctx:   ctx,
		store: store,
	}
}

func (h *MDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
}

// Fetch ... fetches the data in the provided MongoDB collection and returns the results
//
//	@Summary		Fetch records
//	@Description	fetch records by filtering
//	@Param			request	body	model.MDBRequest	true	"query params"
//	@Tags			MongoDB
//	@Success		200	{object}	model.MDBResponse
//	@Failure		500	{object}	model.MDBResponse
//	@Router			/mdb [get]
func (h *MDBHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	var reqModel model.MDBRequest
	if err := json.NewDecoder(r.Body).Decode(&reqModel); err != nil {
		h.response(w, http.StatusBadRequest, "could not parse request body", []model.Record{})
		return
	}

	// validate the format of the date fields
	match, _ := regexp.MatchString(dateRegexp, reqModel.EndDate)
	if !match {
		h.response(w, http.StatusBadRequest, "invalid format for enddate", []model.Record{})
		return
	}

	match, _ = regexp.MatchString(dateRegexp, reqModel.StartDate)
	if !match {
		h.response(w, http.StatusBadRequest, "invalid format for startdate", []model.Record{})
		return
	}

	// set the filters to query mongodb collection.
	filter := bson.D{
		{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: reqModel.StartDate}}},
		{Key: "createdAt", Value: bson.D{{Key: "$lte", Value: reqModel.EndDate}}},
		{Key: "totalCount", Value: bson.D{{Key: "$gte", Value: reqModel.MinCount}}},
		{Key: "totalCount", Value: bson.D{{Key: "$lte", Value: reqModel.MaxCount}}},
	}

	records, err := h.store.Find(h.ctx, filter)
	if err != nil {
		h.response(w, http.StatusInternalServerError, "could not retrieve the records from the database", []model.Record{})
		return
	}

	h.response(w, http.StatusOK, "success", records)
}

func (h *MDBHandler) response(w http.ResponseWriter, status int, message string, records []model.Record) {
	code := 0

	if status != http.StatusOK {
		code = -1
	}

	response := model.MDBResponse{
		Records: records,
		Code:    code,
		Msg:     message,
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not serialize response"))
	}
}
