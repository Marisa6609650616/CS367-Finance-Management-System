package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"CS367-Finance-Management-System/middleware"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	// 1. Mock Database
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to open mock: %s", err)
	}
	defer db.Close()

	// 2. Mock SQL
	mock.ExpectQuery("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'income'").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"SUM(amount)"}).AddRow(1000.0))

	mock.ExpectQuery("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = ? AND type = 'expense'").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"SUM(amount)"}).AddRow(400.0))

	// 3. สร้าง Request และ Response Recorder
	req, _ := http.NewRequest("GET", "/api/summary/balance", nil)

	// จำลองค่า UserID เข้าไปใน Context (เหมือนที่ Middleware ของเพื่อนทำ)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := GetBalance(db)

	// 4. รัน Handler
	handler.ServeHTTP(rr, req)

	// 5. ตรวจสอบผล
	assert.Equal(t, http.StatusOK, rr.Code)

	var response SummaryResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, 600.0, response.Balance)
	assert.Equal(t, "success", response.Status)
}
