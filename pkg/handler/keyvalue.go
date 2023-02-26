package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/iboware/go-rest-api-example/pkg/model"
)

type KVStore struct {
	m map[string]string
	*sync.RWMutex
}

type KeyValueHandler struct {
	store *KVStore
}

func NewKeyValueHandler() *KeyValueHandler {
	return &KeyValueHandler{
		store: &KVStore{
			m:       make(map[string]string, 0),
			RWMutex: &sync.RWMutex{},
		},
	}
}

func (h *KeyValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
}

// Create ... creates a new cache entry.
//
//	@Summary		creates a cache entry
//	@Description	creates a cache entry in-memory
//	@Tags			In-Memory
//	@Accept			json
//	@Param			tuple	body		model.Tuple	true	"Key Value data"
//	@Success		200		{object}	model.Tuple
//	@Failure		500		{object}	object
//	@Router			/in-memory [post]
func (h *KeyValueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.Tuple
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.response(w, http.StatusInternalServerError, "Could not parse request body", nil)
		return
	}

	h.store.Lock()
	h.store.m[t.Key] = t.Value
	h.store.Unlock()

	h.response(w, http.StatusOK, "Success", &t)
}

// Fetch ... fetches the value of a key from in memory cache.
//
//	@Summary		fetches a cache entry
//	@Description	fetches cache entries in-memory
//	@Tags			In-Memory
//	@Param			key	query		string	true	"Key"
//	@Success		200	{object}	model.Tuple
//	@Failure		404	{string}	string
//	@Failure		500	{string}	string
//	@Router			/in-memory [get]
func (h *KeyValueHandler) Fetch(w http.ResponseWriter, r *http.Request) {

	if !r.URL.Query().Has("key") {
		h.response(w, http.StatusBadRequest, "Request does not contain required query parameter: key", nil)
		return
	}

	key := r.URL.Query().Get("key")

	h.store.RLock()
	val, ok := h.store.m[key]
	h.store.RUnlock()

	if !ok {
		h.response(w, http.StatusNotFound, "key not found", nil)
		return
	}
	h.response(w, http.StatusOK, "", &model.Tuple{
		Key:   key,
		Value: val,
	})

}

// response formats body and header and writes it as json.
func (h *KeyValueHandler) response(w http.ResponseWriter, status int, message string, tuple *model.Tuple) {
	w.WriteHeader(status)
	if status != http.StatusOK {
		w.Write([]byte(message))
		return
	}

	if err := json.NewEncoder(w).Encode(tuple); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not serialize response"))
	}
}
