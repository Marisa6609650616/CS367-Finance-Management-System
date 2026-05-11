package handlers_test

import (
	"CS367-Finance-Management-System/handlers"
	"CS367-Finance-Management-System/middleware"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"

	_ "modernc.org/sqlite"
)

// setupTestDB creates an in-memory SQLite DB with schema + seed data
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE transactions (
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
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

// seedTransactions inserts test rows
func seedTransactions(t *testing.T, db *sql.DB, rows []map[string]interface{}) {
	for _, row := range rows {
		_, err := db.Exec(
			"INSERT INTO transactions (user_id, category_id, type, amount, note, date) VALUES (?, ?, ?, ?, ?, ?)",
			row["user_id"], row["category_id"], row["type"], row["amount"], row["note"], row["date"],
		)
		if err != nil {
			t.Fatalf("failed to seed transaction: %v", err)
		}
	}
}

// makeRequest builds a request with userID injected into context (mimics AuthMiddleware)
func makeRequest(userID int, month string) *http.Request {
	url := "/api/summary/monthly"
	if month != "" {
		url += "?month=" + month
	}
	req := httptest.NewRequest(http.MethodGet, url, nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

// --- Tests ---

func TestSummaryMonthly_BasicIncomeAndExpense(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	seedTransactions(t, db, []map[string]interface{}{
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 30000.0, "note": "", "date": "2025-04-01"},
		{"user_id": 1, "category_id": 2, "type": "expense", "amount": 5000.0, "note": "", "date": "2025-04-15"},
		{"user_id": 1, "category_id": 2, "type": "expense", "amount": 2000.0, "note": "", "date": "2025-04-20"},
	})

	h := handlers.NewSumaryMonthyHandler(db)
	rr := httptest.NewRecorder()
	h.SummaryMonthly(rr, makeRequest(1, "2025-04"))

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	summary := resp["summary"].(map[string]interface{})

	if summary["income"] != 30000.0 {
		t.Errorf("expected income 30000, got %v", summary["income"])
	}
	if summary["expense"] != 7000.0 {
		t.Errorf("expected expense 7000, got %v", summary["expense"])
	}
	if summary["balance"] != 23000.0 {
		t.Errorf("expected balance 23000, got %v", summary["balance"])
	}
}

func TestSummaryMonthly_OnlySeesOwnTransactions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	seedTransactions(t, db, []map[string]interface{}{
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 10000.0, "note": "", "date": "2025-04-01"},
		{"user_id": 2, "category_id": 1, "type": "income", "amount": 99999.0, "note": "", "date": "2025-04-01"}, // คนอื่น
	})

	h := handlers.NewSumaryMonthyHandler(db)
	rr := httptest.NewRecorder()
	h.SummaryMonthly(rr, makeRequest(1, "2025-04"))

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	summary := resp["summary"].(map[string]interface{})

	if summary["income"] != 10000.0 {
		t.Errorf("expected income 10000 (own only), got %v", summary["income"])
	}
}

func TestSummaryMonthly_NoTransactionsReturnsZero(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// ไม่มีข้อมูลเลย
	h := handlers.NewSumaryMonthyHandler(db)
	rr := httptest.NewRecorder()
	h.SummaryMonthly(rr, makeRequest(1, "2025-04"))

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	summary := resp["summary"].(map[string]interface{})

	if summary["income"] != 0.0 {
		t.Errorf("expected income 0, got %v", summary["income"])
	}
	if summary["expense"] != 0.0 {
		t.Errorf("expected expense 0, got %v", summary["expense"])
	}
	if summary["balance"] != 0.0 {
		t.Errorf("expected balance 0, got %v", summary["balance"])
	}
}

func TestSummaryMonthly_FiltersByMonthCorrectly(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	seedTransactions(t, db, []map[string]interface{}{
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 5000.0, "note": "", "date": "2025-03-31"},  // เดือนก่อน
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 20000.0, "note": "", "date": "2025-04-01"}, // เดือนนี้
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 8000.0, "note": "", "date": "2025-05-01"},  // เดือนหน้า
	})

	h := handlers.NewSumaryMonthyHandler(db)
	rr := httptest.NewRecorder()
	h.SummaryMonthly(rr, makeRequest(1, "2025-04"))

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	summary := resp["summary"].(map[string]interface{})

	if summary["income"] != 20000.0 {
		t.Errorf("expected income 20000 (April only), got %v", summary["income"])
	}
}

func TestSummaryMonthly_Unauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	h := handlers.NewSumaryMonthyHandler(db)
	rr := httptest.NewRecorder()

	// ไม่ inject userID เข้า context
	req := httptest.NewRequest(http.MethodGet, "/api/summary/monthly?month=2025-04", nil)
	h.SummaryMonthly(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
