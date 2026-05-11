// handlers/transactions_test.go
package handlers

import (
	"CS367-Finance-Management-System/middleware"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// สร้าง in-memory DB สำหรับ test
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			user_id INTEGER
		);
		CREATE TABLE transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			amount REAL NOT NULL,
			note TEXT,
			date TEXT NOT NULL
		);
		INSERT INTO categories (name, type) VALUES ('อาหาร', 'expense');
	`)
	if err != nil {
		t.Fatalf("failed to setup schema: %v", err)
	}
	return db
}

// helper: สร้าง request พร้อม user_id ใน context
func newRequestWithUser(body interface{}, userID int) *http.Request {
	b, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(b))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	return r.WithContext(ctx)
}

// ✅ เคส 1: เพิ่มรายการสำเร็จ
func TestCreateTransaction_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewTransactionHandler(db)

	body := map[string]interface{}{
		"category_id": 1,
		"type":        "expense",
		"amount":      120.0,
		"note":        "ข้าวกลางวัน",
		"date":        "2026-05-01",
	}

	r := newRequestWithUser(body, 1)
	w := httptest.NewRecorder()
	h.CreateTransaction(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

// ❌ เคส 2: ไม่มี user_id ใน context (Unauthorized)
func TestCreateTransaction_Unauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewTransactionHandler(db)

	body := map[string]interface{}{
		"category_id": 1,
		"type":        "expense",
		"amount":      120.0,
		"date":        "2026-05-01",
	}

	b, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(b))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTransaction(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// ❌ เคส 3: ข้อมูลไม่ครบ (ไม่มี amount)
func TestCreateTransaction_MissingFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewTransactionHandler(db)

	body := map[string]interface{}{
		"category_id": 1,
		"type":        "expense",
		"date":        "2026-05-01",
		// ไม่มี amount
	}

	r := newRequestWithUser(body, 1)
	w := httptest.NewRecorder()
	h.CreateTransaction(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ❌ เคส 4: type ไม่ถูกต้อง
func TestCreateTransaction_InvalidType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewTransactionHandler(db)

	body := map[string]interface{}{
		"category_id": 1,
		"type":        "invalid",
		"amount":      100.0,
		"date":        "2026-05-01",
	}

	r := newRequestWithUser(body, 1)
	w := httptest.NewRecorder()
	h.CreateTransaction(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ❌ เคส 5: category ไม่มีอยู่ใน DB
func TestCreateTransaction_CategoryNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewTransactionHandler(db)

	body := map[string]interface{}{
		"category_id": 999,
		"type":        "expense",
		"amount":      100.0,
		"date":        "2026-05-01",
	}

	r := newRequestWithUser(body, 1)
	w := httptest.NewRecorder()
	h.CreateTransaction(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// ❌ เคส 6: amount เป็น 0 หรือติดลบ
func TestCreateTransaction_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := NewTransactionHandler(db)

	body := map[string]interface{}{
		"category_id": 1,
		"type":        "expense",
		"amount":      -50.0,
		"date":        "2026-05-01",
	}

	r := newRequestWithUser(body, 1)
	w := httptest.NewRecorder()
	h.CreateTransaction(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
