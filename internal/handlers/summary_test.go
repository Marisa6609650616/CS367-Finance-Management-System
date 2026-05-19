package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"CS367-Finance-Management-System/internal/handlers"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// makeBalanceContext สร้าง gin.Context สำหรับทดสอบ GetBalance
func makeBalanceContext(userID int) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/summary/balance", nil)
	c.Set("user_id", userID)
	return c, rr
}

// ─────────────────────────────────────────
// GET /api/summary/balance Tests
// ─────────────────────────────────────────

func TestGetBalance_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(`.*income.*`).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(1000.0))
	mock.ExpectQuery(`.*expense.*`).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(400.0))

	h := handlers.NewBalanceHandler(db)
	c, rr := makeBalanceContext(1)
	h.GetBalance(c)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["balance"] != 600.0 {
		t.Errorf("expected balance 600, got %v", resp["balance"])
	}
	if resp["status"] != "success" {
		t.Errorf("expected status=success, got %v", resp["status"])
	}
}

func TestGetBalance_AllExpense(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`.*income.*`).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(0.0))
	mock.ExpectQuery(`.*expense.*`).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(500.0))

	h := handlers.NewBalanceHandler(db)
	c, rr := makeBalanceContext(1)
	h.GetBalance(c)

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["balance"] != -500.0 {
		t.Errorf("expected balance -500, got %v", resp["balance"])
	}
}

func TestGetBalance_NoTransactions(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`.*income.*`).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(0.0))
	mock.ExpectQuery(`.*expense.*`).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(0.0))

	h := handlers.NewBalanceHandler(db)
	c, rr := makeBalanceContext(1)
	h.GetBalance(c)

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["balance"] != 0.0 {
		t.Errorf("expected balance 0, got %v", resp["balance"])
	}
}
