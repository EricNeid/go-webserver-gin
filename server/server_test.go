package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EricNeid/go-webserver-gin/internal/verify"
	"github.com/gin-gonic/gin"
)

func TestNormalizePath(t *testing.T) {
	// action
	result := normalizePath("")
	// verify
	verify.Equals(t, "", result)

	// action
	result = normalizePath("/")
	// verify
	verify.Equals(t, "", result)

	// action
	result = normalizePath("/test")
	// verify
	verify.Equals(t, "/test", result)

	// action
	result = normalizePath("test")
	// verify
	verify.Equals(t, "/test", result)

	// action
	result = normalizePath("test/")
	// verify
	verify.Equals(t, "/test", result)
}

func TestWelcome(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()
	unit := NewApplicationServer(":5001", "", nil)
	// action
	unit.Router.ServeHTTP(rec, req)
	// verify
	verify.Assert(t, rec.Code == 200, fmt.Sprintf("Status code is %d\n", rec.Code))
	verify.Assert(t, rec.Body.String() == "Hello, World!", fmt.Sprintf("Body is %s\n", rec.Body.String()))
}

func TestInstallGETHandler(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/foobar", http.NoBody)
	rec := httptest.NewRecorder()
	unit := NewApplicationServer(":5001", "", nil)
	// action
	unit.InstallGETHandler("/foobar", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })
	unit.Router.ServeHTTP(rec, req)
	// verify
	verify.Assert(t, rec.Code == 204, fmt.Sprintf("Status code is %d\n", rec.Code))
}

func TestInstallGETHandler_shouldEnsurePath(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/foobar", http.NoBody)
	rec := httptest.NewRecorder()
	unit := NewApplicationServer(":5001", "", nil)
	// action
	unit.InstallGETHandler("foobar/", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })
	unit.Router.ServeHTTP(rec, req)
	// verify
	verify.Assert(t, rec.Code == 204, fmt.Sprintf("Status code is %d\n", rec.Code))
}

func TestInstallGETHandler_withBasePath(t *testing.T) {
	// arrange
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/base/foobar", http.NoBody)
	rec := httptest.NewRecorder()
	unit := NewApplicationServer(":5001", "base", nil)
	// action
	unit.InstallGETHandler("foobar", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })
	unit.Router.ServeHTTP(rec, req)
	// verify
	verify.Assert(t, rec.Code == 204, fmt.Sprintf("Status code is %d\n", rec.Code))
}
