package handlers

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"tasks-api/internal/models"
	"tasks-api/internal/storage/memory"
	"testing"
)

func TestAPI(t *testing.T) {
	h := New(memory.New())
	cases := []struct {
		method, path, body string
		status             int
		contains           string
	}{
		{"GET", "/tasks", "", 200, "[]"},
		{"POST", "/tasks", `{"title":"Learn Go"}`, 201, `"id":1`},
		{"GET", "/tasks/1", "", 200, `"done":false`},
		{"PUT", "/tasks/1", `{"title":"Done","done":true}`, 200, `"done":true`},
		{"GET", "/tasks", "", 200, `"title":"Done"`},
		{"PUT", "/tasks/1", `{"title":"Missing done"}`, 400, "error"},
		{"GET", "/tasks/1", "", 200, `"title":"Done"`},
		{"POST", "/tasks", `{"title":" "}`, 400, "error"},
		{"POST", "/tasks", `{"title":3}`, 400, "error"},
		{"POST", "/tasks", `{"title":"x","done":"yes"}`, 400, "error"},
		{"POST", "/tasks", `{"title":"x","id":9}`, 400, "error"},
		{"POST", "/tasks", `{"title":"x"} {}`, 400, "error"},
		{"POST", "/tasks", `null`, 400, "error"},
		{"POST", "/tasks", `{"title":"x","done":null}`, 400, "error"},
		{"POST", "/tasks", `[]`, 400, "error"},
		{"POST", "/tasks", ``, 400, "error"},
		{"POST", "/tasks", `{"title":"` + strings.Repeat("a", 1<<20) + `"}`, 400, "error"},
		{"GET", "/tasks/no", "", 400, "error"},
		{"GET", "/tasks/0", "", 400, "error"},
		{"GET", "/tasks/999", "", 404, "error"},
		{"PUT", "/tasks/999", `{"title":"x","done":false}`, 404, "error"},
		{"PATCH", "/tasks/1", "", 405, "error"},
		{"DELETE", "/tasks", "", 405, "error"},
		{"GET", "/missing", "", 404, "error"},
		{"GET", "/tasks/", "", 404, "error"},
		{"GET", "/tasks/1/extra", "", 404, "error"},
		{"GET", "/health", "", 200, `"status":"ok"`},
		{"POST", "/health", "", 405, "error"},
		{"DELETE", "/tasks/1", "", 204, ""},
		{"DELETE", "/tasks/1", "", 404, "error"},
		{"GET", "/tasks", "", 200, "[]"},
	}
	for _, tc := range cases {
		t.Run(tc.method+tc.path+tc.contains, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status %d, want %d: %s", w.Code, tc.status, w.Body)
			}
			if w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("missing JSON content type")
			}
			if tc.status == 204 {
				if w.Body.Len() != 0 {
					t.Fatal("204 must have no body")
				}
			} else if !json.Valid(w.Body.Bytes()) || !strings.Contains(w.Body.String(), tc.contains) {
				t.Fatalf("unexpected body: %s", w.Body)
			}
			if tc.status == 405 && w.Header().Get("Allow") == "" {
				t.Fatal("missing Allow")
			}
			if tc.status == 201 && w.Header().Get("Location") != "/tasks/1" {
				t.Fatal("missing Location")
			}
		})
	}
}

func TestContentType(t *testing.T) {
	w := httptest.NewRecorder()
	New(memory.New()).ServeHTTP(w, httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"x"}`)))
	if w.Code != 400 {
		t.Fatalf("got %d", w.Code)
	}
}

type brokenStore struct{ *memory.Store }

func (s brokenStore) Create(models.Task) (models.Task, error) {
	return models.Task{}, errors.New("private details")
}
func (s brokenStore) Update(int, models.Task) (models.Task, error) {
	return models.Task{}, errors.New("private details")
}
func (s brokenStore) Delete(int) error { return errors.New("private details") }

func TestStorageFailures(t *testing.T) {
	for _, method := range []string{"POST", "PUT", "DELETE"} {
		path := "/tasks/1"
		if method == "POST" {
			path = "/tasks"
		}
		r := httptest.NewRequest(method, path, strings.NewReader(`{"title":"x","done":false}`))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		New(brokenStore{memory.New()}).ServeHTTP(w, r)
		if w.Code != 500 || w.Body.String() != "{\"error\":\"internal server error\"}\n" {
			t.Fatalf("%s: %d %s", method, w.Code, w.Body)
		}
	}
}
