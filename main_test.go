package main

import (
	"APPDROP/db"
	"APPDROP/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.RegisterRoutes(r)
	return r
}

func testBrandAndCookie(t *testing.T, r *gin.Engine) (domain, cookie string) {
	t.Helper()
	if db.DB == nil {
		t.Skip("DATABASE_URL not set, skipping test that requires brand and auth")
	}
	domain = "testbrand"

	body := `{"name":"Test Brand","domain":"testbrand","email":"test@testbrand.com","password":"secret","office_address":"","logo":""}`
	req := httptest.NewRequest(http.MethodPost, "/brands", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated && w.Code != http.StatusConflict {
		t.Fatalf("create brand: got status %d", w.Code)
	}

	loginBody := `{"email":"test@testbrand.com","password":"secret"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("X-Brand-Domain", domain)
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Fatalf("login: got status %d, body %s", loginW.Code, loginW.Body.String())
	}
	setCookie := loginW.Header().Get("Set-Cookie")
	if setCookie == "" {
		t.Fatal("login did not return Set-Cookie")
	}
	if idx := strings.Index(setCookie, ";"); idx != -1 {
		cookie = strings.TrimSpace(setCookie[:idx])
	} else {
		cookie = setCookie
	}
	return domain, cookie
}

func TestHealth(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /health: got status %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("GET /health: got status %q, want %q", body["status"], "ok")
	}
}

func TestGetPageByID_InvalidUUID(t *testing.T) {
	r := testRouter()
	domain, cookie := testBrandAndCookie(t, r)
	req := httptest.NewRequest(http.MethodGet, "/pages/not-a-uuid", nil)
	req.Header.Set("X-Brand-Domain", domain)
	req.Header.Set("Cookie", cookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("GET /pages/invalid: got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreatePage_ValidationErrors(t *testing.T) {
	r := testRouter()
	domain, cookie := testBrandAndCookie(t, r)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"empty body", "{}", http.StatusBadRequest},
		{"missing name", `{"route": "/about"}`, http.StatusBadRequest},
		{"missing route", `{"name": "About"}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Brand-Domain", domain)
			req.Header.Set("Cookie", cookie)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Errorf("POST /pages: got status %d, want %d", w.Code, tt.want)
			}
		})
	}
}

func TestAddWidget_InvalidPageID(t *testing.T) {
	r := testRouter()
	domain, cookie := testBrandAndCookie(t, r)
	body := `{"type": "banner", "position": 0}`
	req := httptest.NewRequest(http.MethodPost, "/pages/not-a-uuid/widgets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Brand-Domain", domain)
	req.Header.Set("Cookie", cookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /pages/invalid/widgets: got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAddWidget_InvalidType(t *testing.T) {
	r := testRouter()
	domain, cookie := testBrandAndCookie(t, r)
	// Create a page first so we have a valid page ID (unique route to avoid 409 on re-runs)
	route := fmt.Sprintf("/test-widget-%d", time.Now().UnixNano())
	pageBody := fmt.Sprintf(`{"name": "Test Page", "route": "%s", "is_home": false}`, route)
	createReq := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBufferString(pageBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-Brand-Domain", domain)
	createReq.Header.Set("Cookie", cookie)
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Fatalf("could not create test page: got %d", createW.Code)
	}
	var page struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(createW.Body.Bytes(), &page); err != nil {
		t.Fatalf("parse page: %v", err)
	}

	widgetBody := `{"type": "invalid_type", "position": 0}`
	req := httptest.NewRequest(http.MethodPost, "/pages/"+page.ID+"/widgets", bytes.NewBufferString(widgetBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Brand-Domain", domain)
	req.Header.Set("Cookie", cookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST .../widgets with invalid type: got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetPages_RequiresDB(t *testing.T) {
	r := testRouter()
	domain, cookie := testBrandAndCookie(t, r)
	req := httptest.NewRequest(http.MethodGet, "/pages", nil)
	req.Header.Set("X-Brand-Domain", domain)
	req.Header.Set("Cookie", cookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /pages: got status %d, want %d", w.Code, http.StatusOK)
	}

	var pages []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &pages); err != nil {
		t.Errorf("GET /pages: invalid JSON array: %v", err)
	}
}

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		db.Connect()
	}
	code := m.Run()
	os.Exit(code)
}
