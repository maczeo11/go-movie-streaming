package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminOnlyAllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Set("role", "ADMIN")

	AdminOnly()(ctx)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for admin, got %d", rec.Code)
	}
}

func TestAdminOnlyRejectsUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Set("role", "USER")

	AdminOnly()(ctx)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 for user, got %d", rec.Code)
	}
}

func TestAdminOnlyRejectsMissingRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	AdminOnly()(ctx)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 when no role present, got %d", rec.Code)
	}
}
