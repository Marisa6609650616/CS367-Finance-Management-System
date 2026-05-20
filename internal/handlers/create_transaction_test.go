package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"CS367-Finance-Management-System/internal/handlers"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

// setupCreateTxDB สร้าง in-memory DB พร้อม transactions + categories table
func setupCreateTxDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test DB: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			type       TEXT NOT NULL,
			user_id    INTEGER
		);
		CREATE TABLE IF NOT EXISTS transactions (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id     INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			type        TEXT    NOT NULL,
			amount      REAL    NOT NULL,
			note        TEXT,
			date        TEXT    NOT NULL,
			created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
		);
		INSERT INTO categories (id, name, type) VALUES (1, 'อาหาร', 'expense');
		INSERT INTO categories (id, name, type) VALUES (2, 'เงินเดือน', 'income');
	`)
	if err != nil {
		t.Fatalf("setup DB: %v", err)
	}
	return db
}

func newCreateTxRouter(db *sql.DB) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})
	h := handlers.NewTransactionHandler(db)
	r.POST("/api/transactions", h.CreateTransaction)
	return r
}

func txRequest(r *gin.Engine, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// ─────────────────────────────────────────
// POST /api/transactions
// ─────────────────────────────────────────

func TestCreateTransaction_Success(t *testing.T) {
	db := setupCreateTxDB(t)
	r := newCreateTxRouter(db)
	w := txRequest(r, `{"category_id":1,"type":"expense","amount":120,"note":"ข้าว","date":"2026-05-19"}`)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d | %s", w.Code, w.Body.String())
	}
	var body map[string]interface{}
	json.NewDecoder(w.Body).Decode(&body)
	if body["message"] != "เพิ่มรายการสำเร็จ" {
		t.Errorf("unexpected message: %v", body["message"])
	}
}

func TestCreateTransaction_InvalidType(t *testing.T) {
	db := setupCreateTxDB(t)
	r := newCreateTxRouter(db)
	w := txRequest(r, `{"category_id":1,"type":"invalid","amount":120,"date":"2026-05-19"}`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid type, got %d", w.Code)
	}
}

func TestCreateTransaction_MissingFields(t *testing.T) {
	db := setupCreateTxDB(t)
	r := newCreateTxRouter(db)
	// ไม่มี date
	w := txRequest(r, `{"category_id":1,"type":"expense","amount":120}`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing date, got %d", w.Code)
	}
}

func TestCreateTransaction_CategoryNotFound(t *testing.T) {
	db := setupCreateTxDB(t)
	r := newCreateTxRouter(db)
	w := txRequest(r, `{"category_id":999,"type":"expense","amount":120,"date":"2026-05-19"}`)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for invalid category, got %d", w.Code)
	}
}

func TestCreateTransaction_InvalidJSON(t *testing.T) {
	db := setupCreateTxDB(t)
	r := newCreateTxRouter(db)
	// ส่ง JSON ที่ malformed
	w := txRequest(r, `{invalid-json}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestCreateTransaction_ZeroAmount(t *testing.T) {
	db := setupCreateTxDB(t)
	r := newCreateTxRouter(db)
	w := txRequest(r, `{"category_id":1,"type":"expense","amount":0,"date":"2026-05-19"}`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero amount, got %d", w.Code)
	}
}
