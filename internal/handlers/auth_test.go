package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"CS367-Finance-Management-System/internal/handlers"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

// setupAuthDB สร้าง in-memory DB พร้อมตาราง users สำหรับทดสอบ auth
func setupAuthDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test DB: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			email      TEXT    NOT NULL UNIQUE,
			password   TEXT    NOT NULL,
			name       TEXT    NOT NULL,
			role       TEXT    NOT NULL DEFAULT 'user',
			created_at TEXT    NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT    NOT NULL DEFAULT (datetime('now'))
		);
	`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}
	return db
}

func newAuthTestRouter(db *sql.DB) *gin.Engine {
	r := gin.New()
	h := &handlers.AuthHandler{DB: db}
	r.POST("/api/auth/register", h.Register)
	r.POST("/api/auth/login", h.Login)
	return r
}

func authRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func parseAuthBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("parse body: %v | raw: %s", err, w.Body.String())
	}
	return m
}

// ─────────────────────────────────────────
// POST /api/auth/register
// ─────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	r := newAuthTestRouter(setupAuthDB(t))
	w := authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d | %s", w.Code, w.Body.String())
	}
	body := parseAuthBody(t, w)
	if body["message"] != "registered successfully" {
		t.Errorf("unexpected message: %v", body["message"])
	}
}

func TestRegister_ResponseHasUserObject(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	r := newAuthTestRouter(setupAuthDB(t))
	w := authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)

	body := parseAuthBody(t, w)
	user, ok := body["user"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'user' object in response")
	}
	if user["email"] != "sattawut@cs367.com" {
		t.Errorf("expected email, got %v", user["email"])
	}
	if user["role"] != "user" {
		t.Errorf("expected role=user, got %v", user["role"])
	}
}

func TestRegister_PasswordNotExposedInResponse(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	r := newAuthTestRouter(setupAuthDB(t))
	w := authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)

	if strings.Contains(w.Body.String(), "password") {
		t.Errorf("password must not appear in response: %s", w.Body.String())
	}
}

func TestRegister_EmailNormalized(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	r := newAuthTestRouter(setupAuthDB(t))
	w := authRequest(r, "POST", "/api/auth/register",
		`{"email":"SATTAWUT@CS367.COM","password":"password123","name":"เสฎฐวุฒิ"}`)

	body := parseAuthBody(t, w)
	user := body["user"].(map[string]interface{})
	if user["email"] != "sattawut@cs367.com" {
		t.Errorf("expected lowercase email, got %v", user["email"])
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	db := setupAuthDB(t)
	r := newAuthTestRouter(db)

	authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)
	w := authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	cases := []struct {
		name string
		body string
	}{
		{"no email", `{"password":"password123","name":"test"}`},
		{"no password", `{"email":"a@b.com","name":"test"}`},
		{"no name", `{"email":"a@b.com","password":"password123"}`},
		{"short password", `{"email":"a@b.com","password":"123","name":"test"}`},
		{"invalid email", `{"email":"not-email","password":"password123","name":"test"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newAuthTestRouter(setupAuthDB(t))
			w := authRequest(r, "POST", "/api/auth/register", tc.body)
			if w.Code != http.StatusBadRequest {
				t.Errorf("[%s] expected 400, got %d", tc.name, w.Code)
			}
		})
	}
}

// ─────────────────────────────────────────
// POST /api/auth/login
// ─────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	db := setupAuthDB(t)
	r := newAuthTestRouter(db)

	authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)
	w := authRequest(r, "POST", "/api/auth/login",
		`{"email":"sattawut@cs367.com","password":"password123"}`)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d | %s", w.Code, w.Body.String())
	}
	body := parseAuthBody(t, w)
	if body["token"] == nil || body["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	db := setupAuthDB(t)
	r := newAuthTestRouter(db)

	authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)
	w := authRequest(r, "POST", "/api/auth/login",
		`{"email":"sattawut@cs367.com","password":"wrongpass"}`)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_EmailNotFound(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	r := newAuthTestRouter(setupAuthDB(t))
	w := authRequest(r, "POST", "/api/auth/login",
		`{"email":"ghost@cs367.com","password":"password123"}`)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_SameErrorMessage(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	db := setupAuthDB(t)
	r := newAuthTestRouter(db)

	authRequest(r, "POST", "/api/auth/register",
		`{"email":"sattawut@cs367.com","password":"password123","name":"เสฎฐวุฒิ"}`)

	w1 := authRequest(r, "POST", "/api/auth/login",
		`{"email":"sattawut@cs367.com","password":"wrongpass"}`)
	w2 := authRequest(r, "POST", "/api/auth/login",
		`{"email":"nobody@cs367.com","password":"password123"}`)

	b1 := parseAuthBody(t, w1)
	b2 := parseAuthBody(t, w2)
	if b1["error"] != b2["error"] {
		t.Errorf("error messages must match to prevent user enumeration: %v vs %v", b1["error"], b2["error"])
	}
}

func TestLogin_MissingFields(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	cases := []struct {
		name string
		body string
	}{
		{"no email", `{"password":"password123"}`},
		{"no password", `{"email":"a@b.com"}`},
		{"invalid email", `{"email":"not-email","password":"password123"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newAuthTestRouter(setupAuthDB(t))
			w := authRequest(r, "POST", "/api/auth/login", tc.body)
			if w.Code != http.StatusBadRequest {
				t.Errorf("[%s] expected 400, got %d", tc.name, w.Code)
			}
		})
	}
}
