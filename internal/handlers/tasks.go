package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		status, body = http.StatusInternalServerError, []byte(`{"error":"internal server error"}`)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}

func fail(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	fail(w, http.StatusMethodNotAllowed, "method not allowed")
}

// ServeHTTP keeps unknown routes and unsupported methods in the JSON contract.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/tasks":
		h.TasksCollection(w, r)
	case strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.TaskItem(w, r)
	case r.URL.Path == "/health":
		if r.Method != http.MethodGet {
			methodNotAllowed(w, "GET")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	default:
		fail(w, http.StatusNotFound, "resource not found")
	}
}

func decodeTask(w http.ResponseWriter, r *http.Request) (models.Task, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		fail(w, http.StatusBadRequest, "Content-Type must be application/json")
		return models.Task{}, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	// RawMessage distinguishes an omitted done from an explicit null.
	var input struct {
		Title *string         `json:"title"`
		Done  json.RawMessage `json:"done"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON body")
		return models.Task{}, false
	}
	if err := decoder.Decode(new(interface{})); err != io.EOF {
		fail(w, http.StatusBadRequest, "body must contain a single JSON object")
		return models.Task{}, false
	}
	if input.Title == nil || strings.TrimSpace(*input.Title) == "" {
		fail(w, http.StatusBadRequest, "title is required and must not be blank")
		return models.Task{}, false
	}
	if r.Method == http.MethodPut && input.Done == nil {
		fail(w, http.StatusBadRequest, "done is required for PUT")
		return models.Task{}, false
	}
	task := models.Task{Title: strings.TrimSpace(*input.Title)}
	if input.Done != nil {
		if string(input.Done) == "null" || json.Unmarshal(input.Done, &task.Done) != nil {
			fail(w, http.StatusBadRequest, "done must be a boolean")
			return models.Task{}, false
		}
	}
	return task, true
}

func storageError(w http.ResponseWriter, err error) {
	if errors.Is(err, storage.ErrNotFound) {
		fail(w, http.StatusNotFound, "task not found")
		return
	}
	log.Printf("storage error: %v", err)
	fail(w, http.StatusInternalServerError, "internal server error")
}

func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := h.Store.List()
		if tasks == nil {
			tasks = []models.Task{}
		}
		writeJSON(w, http.StatusOK, tasks)
	case http.MethodPost:
		task, ok := decodeTask(w, r)
		if !ok {
			return
		}
		created, err := h.Store.Create(task)
		if err != nil {
			storageError(w, err)
			return
		}
		w.Header().Set("Location", "/tasks/"+strconv.Itoa(created.ID))
		writeJSON(w, http.StatusCreated, created)
	default:
		methodNotAllowed(w, "GET, POST")
	}
}

func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	rawID := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if rawID == "" || strings.Contains(rawID, "/") {
		fail(w, http.StatusNotFound, "resource not found")
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 || strings.Trim(rawID, "0123456789") != "" {
		fail(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}
	switch r.Method {
	case http.MethodGet:
		task, ok := h.Store.Get(id)
		if !ok {
			fail(w, http.StatusNotFound, "task not found")
			return
		}
		writeJSON(w, http.StatusOK, task)
	case http.MethodPut:
		task, ok := decodeTask(w, r)
		if !ok {
			return
		}
		updated, err := h.Store.Update(id, task)
		if err != nil {
			storageError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if err := h.Store.Delete(id); err != nil {
			storageError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(w, "GET, PUT, DELETE")
	}
}
