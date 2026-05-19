package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"CS367-Finance-Management-System/internal/middleware"
	"CS367-Finance-Management-System/pkg/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setup() {
	os.Setenv("JWT_SECRET", "test_secret_key_for_unit_tests")
}

func newTestRouter() *gin.Engine {
	r := gin.New()
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	protected.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		})
	})
	return r
}

// Middleware Tests — กรณีที่ต้องถูกบล็อก (401)

func TestMiddleware_NoAuthHeader(t *testing.T) {
	setup()
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for no Authorization header, got %d", w.Code)
	}
}

func TestMiddleware_MissingBearerPrefix(t *testing.T) {
	setup()
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Token somefaketoken")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong Authorization format, got %d", w.Code)
	}
}

func TestMiddleware_BearerWithNoToken(t *testing.T) {
	setup()
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for 'Bearer ' with no token, got %d", w.Code)
	}
}

func TestMiddleware_FakeToken(t *testing.T) {
	setup()
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer fake.token.string")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for fake token string, got %d", w.Code)
	}
}

func TestMiddleware_TamperedToken(t *testing.T) {
	setup()
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	tampered := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJ1c2VyX2lkIjo5OTksInJvbGUiOiJhZG1pbiJ9" +
		".badsignature"
	req.Header.Set("Authorization", "Bearer "+tampered)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for tampered token, got %d", w.Code)
	}
}

func TestMiddleware_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "attacker_secret")
	token, _ := utils.GenerateToken(999, "hacker@inwza.com", "admin")

	os.Setenv("JWT_SECRET", "test_secret_key_for_unit_tests")
	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for token signed with wrong secret, got %d", w.Code)
	}
}

// Middleware Tests — กรณีที่ต้องผ่าน (200)

func TestMiddleware_ValidToken_Passes(t *testing.T) {
	setup()
	token, _ := utils.GenerateToken(1, "test@dome.tu.ac.th.com", "user")

	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for valid token, got %d", w.Code)
	}
}

func TestMiddleware_ValidToken_InjectsContext(t *testing.T) {
	setup()
	token, _ := utils.GenerateToken(7, "MeMyself@gmail.com", "user")

	r := newTestRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Fatal("expected response body, got empty")
	}
	if !contains(body, "MeMyself@gmail.com") {
		t.Errorf("expected email injected into context, got: %s", body)
	}
}

// String Helper

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
