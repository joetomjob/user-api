package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttp(t *testing.T) {
	t.Run("Handler Create", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		body := bytes.NewBufferString(`{"name":"joe", "email":"joetomjob@gmail.com", "age":37}`)
		req := httptest.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h.Create(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Exptected %d, got: %d", http.StatusCreated, rr.Code)
		}

		var got User
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("decode %v", got)
		}

		if got.Id != 1 {
			t.Fatalf("Expected %d, got %d", 1, got.Id)
		}
	})

	t.Run("Handler create - error tests", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		body := bytes.NewBufferString(`{"name":"joe"`)
		req := httptest.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h.Create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

		body = bytes.NewBufferString(`{"name":"joe", "email":"joetomjob@gmail.com", "age":-1}`)
		req = httptest.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		h.Create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

		body = bytes.NewBufferString(`{"name":"", "email":"joetomjob@gmail.com", "age":36}`)
		req = httptest.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		h.Create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

		body = bytes.NewBufferString(`{"name":"joe", "email":"", "age":36}`)
		req = httptest.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		h.Create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

	})

	t.Run("Handler Get", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", "1")

		rr := httptest.NewRecorder()
		h.GetById(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Exptected %d, got: %d", http.StatusOK, rr.Code)
		}

		var got User
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("decode %v", got)
		}

		if got.Id != 1 {
			t.Fatalf("Expected %d, got %d", 1, got.Id)
		}
	})

	t.Run("Handler Get - error test", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		req := httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
		req.SetPathValue("id", "100")

		rr := httptest.NewRecorder()
		h.GetById(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("Exptected %d, got: %d", http.StatusNotFound, rr.Code)
		}

		req = httptest.NewRequest(http.MethodGet, "/users", http.NoBody)
		req.SetPathValue("id", "-1")
		rr = httptest.NewRecorder()

		h.GetById(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("Exptected %d, got: %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("Handler Update", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		body := bytes.NewBufferString(`{"name":"joe", "email":"joetomjob@gmail.com", "age":37}`)
		req := httptest.NewRequest(http.MethodPut, "/users", body)
		req.SetPathValue("id", "1")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h.Update(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Exptected %d, got: %d", http.StatusOK, rr.Code)
		}

		var got User
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("decode %v", got)
		}

		if got.Id != 1 {
			t.Fatalf("Expected %d, got %d", 1, got.Id)
		}
	})

	t.Run("Handler Update - error tests", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		body := bytes.NewBufferString(`{"name":"joe"`)
		req := httptest.NewRequest(http.MethodPut, "/users", body)
		req.SetPathValue("id", "1")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		h.Update(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

		body = bytes.NewBufferString(`{"name":"joe", "email":"joetomjob@gmail.com", "age":-1}`)
		req = httptest.NewRequest(http.MethodPut, "/users", body)
		req.SetPathValue("id", "1")
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		h.Update(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

		body = bytes.NewBufferString(`{"name":"", "email":"joetomjob@gmail.com", "age":36}`)
		req = httptest.NewRequest(http.MethodPut, "/users", body)
		req.SetPathValue("id", "1")
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		h.Update(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

		body = bytes.NewBufferString(`{"name":"joe", "email":"", "age":36}`)
		req = httptest.NewRequest(http.MethodPut, "/users", body)
		req.SetPathValue("id", "1")
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		h.Update(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}

	})

	t.Run("Handler Delete", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		req := httptest.NewRequest(http.MethodDelete, "/users", http.NoBody)
		req.SetPathValue("id", "1")

		rr := httptest.NewRecorder()
		h.Delete(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Exptected %d, got: %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("Handler Delete - error tests", func(t *testing.T) {
		h := NewHandler(NewService(&FakeRepo{}))

		req := httptest.NewRequest(http.MethodDelete, "/users", http.NoBody)
		req.SetPathValue("id", "-1")

		rr := httptest.NewRecorder()
		h.Delete(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Exptected %d, got: %d", http.StatusBadRequest, rr.Code)
		}
	})
}
