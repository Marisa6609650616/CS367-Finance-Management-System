package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"CS367-Finance-Management-System/internal/handlers"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(db *sql.DB) *gin.Engine {
	handlers.DB = db
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	r.PUT("/api/transactions/:id", handlers.UpdateTransaction)
	r.DELETE("/api/transactions/:id", handlers.DeleteTransaction)
	return r
}

func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	return db, mock
}

// ============================================================
// PUT /api/transactions/:id
// ============================================================

func TestUpdateTransaction_Success(t *testing.T) {
	db, mock := newMock(t)
	cols := []string{"id", "user_id", "type", "amount", "category", "description", "date", "created_at"}

	mock.ExpectQuery("SELECT user_id").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
	mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows(cols).
			AddRow(1, 1, "expense", 150.0, "อาหาร", "", "2025-05-01", "2025-05-01"))

	body, _ := json.Marshal(map[string]interface{}{"amount": 150.0})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateTransaction_NoFields(t *testing.T) {
	db, _ := newMock(t)
	body, _ := json.Marshal(map[string]interface{}{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTransaction_NegativeAmount(t *testing.T) {
	db, _ := newMock(t)
	body, _ := json.Marshal(map[string]interface{}{"amount": -50.0})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTransaction_InvalidType(t *testing.T) {
	db, _ := newMock(t)
	body, _ := json.Marshal(map[string]interface{}{"type": "invalid"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTransaction_NotFound(t *testing.T) {
	db, mock := newMock(t)
	mock.ExpectQuery("SELECT user_id").WillReturnError(sql.ErrNoRows)

	body, _ := json.Marshal(map[string]interface{}{"amount": 150.0})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTransaction_Forbidden(t *testing.T) {
	db, mock := newMock(t)
	mock.ExpectQuery("SELECT user_id").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(2))

	body, _ := json.Marshal(map[string]interface{}{"amount": 150.0})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestUpdateTransaction_InvalidID(t *testing.T) {
	db, _ := newMock(t)
	body, _ := json.Marshal(map[string]interface{}{"amount": 150.0})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/transactions/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ============================================================
// DELETE /api/transactions/:id
// ============================================================

func TestDeleteTransaction_Success(t *testing.T) {
	db, mock := newMock(t)
	mock.ExpectQuery("SELECT user_id").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
	mock.ExpectExec("DELETE").
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/transactions/1", nil)
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	db, mock := newMock(t)
	mock.ExpectQuery("SELECT user_id").WillReturnError(sql.ErrNoRows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/transactions/999", nil)
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTransaction_AlreadyDeleted(t *testing.T) {
	db, mock := newMock(t)
	mock.ExpectQuery("SELECT user_id").WillReturnError(sql.ErrNoRows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/transactions/1", nil)
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeleteTransaction_Forbidden(t *testing.T) {
	db, mock := newMock(t)
	mock.ExpectQuery("SELECT user_id").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(2))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/transactions/1", nil)
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDeleteTransaction_InvalidID(t *testing.T) {
	db, _ := newMock(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/transactions/abc", nil)
	setupRouter(db).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
