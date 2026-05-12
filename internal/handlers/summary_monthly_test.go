package handlers_test

import (
	"CS367-Finance-Management-System/internal/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"

	"github.com/gin-gonic/gin"
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

// makeGinContext builds a gin.Context with userID set and optional ?month= query param
func makeGinContext(userID int, month string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	url := "/api/summary/monthly"
	if month != "" {
		url += "?month=" + month
	}
	c.Request = httptest.NewRequest(http.MethodGet, url, nil)

	// inject userID the same way AuthMiddleware does: c.Set("userID", userID)
	c.Set("userID", userID)

	return c, rr
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
	c, rr := makeGinContext(1, "2025-04")
	h.SummaryMonthly(c)

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
		{"user_id": 2, "category_id": 1, "type": "income", "amount": 99999.0, "note": "", "date": "2025-04-01"}, // other user
	})

	h := handlers.NewSumaryMonthyHandler(db)
	c, rr := makeGinContext(1, "2025-04")
	h.SummaryMonthly(c)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

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

	h := handlers.NewSumaryMonthyHandler(db)
	c, rr := makeGinContext(1, "2025-04")
	h.SummaryMonthly(c)

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
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 5000.0, "note": "", "date": "2025-03-31"},  // prev month
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 20000.0, "note": "", "date": "2025-04-01"}, // this month
		{"user_id": 1, "category_id": 1, "type": "income", "amount": 8000.0, "note": "", "date": "2025-05-01"},  // next month
	})

	h := handlers.NewSumaryMonthyHandler(db)
	c, rr := makeGinContext(1, "2025-04")
	h.SummaryMonthly(c)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

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

	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/summary/monthly?month=2025-04", nil)
	// ไม่ Set("userID") — จำลอง request ที่ไม่ผ่าน auth

	h := handlers.NewSumaryMonthyHandler(db)
	h.SummaryMonthly(c)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
