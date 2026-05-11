package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)) // บังคับเช็ค SQL เป๊ะๆ
	if err != nil {
		t.Fatalf("failed to open mock: %s", err)
	}
	defer db.Close()

	// Mock SQL ให้ตรงกับใน summary.go
	mock.ExpectQuery("SELECT SUM(amount) FROM transactions WHERE user_id = ? AND type = 'income'").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"SUM(amount)"}).AddRow(1000.0))

	mock.ExpectQuery("SELECT SUM(amount) FROM transactions WHERE user_id = ? AND type = 'expense'").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"SUM(amount)"}).AddRow(400.0))

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	ctx.Set("user_id", 1)

	r.GET("/api/summary/balance", GetBalance(db))
	req, _ := http.NewRequest("GET", "/api/summary/balance", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response SummaryResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 600.0, response.Balance)
	assert.Equal(t, "success", response.Status)
}
