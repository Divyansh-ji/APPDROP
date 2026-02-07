package main

import (
	"APPDROP/db"
	"APPDROP/routes"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	routes.RegisterRoutes(r)
	return r
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
	req := httptest.NewRequest(http.MethodGet, "/pages/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("GET /pages/invalid: got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreatePage_ValidationErrors(t *testing.T) {
	r := testRouter()

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
	body := `{"type": "banner", "position": 0}`
	req := httptest.NewRequest(http.MethodPost, "/pages/not-a-uuid/widgets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /pages/invalid/widgets: got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAddWidget_InvalidType(t *testing.T) {
	if db.DB == nil {
		t.Skip("DATABASE_URL not set, skipping DB test")
	}
	// Create a page first so we have a valid page ID (we need DB)
	r := testRouter()
	pageBody := `{"name": "Test Page", "route": "/test-widget-validation", "is_home": false}`
	createReq := httptest.NewRequest(http.MethodPost, "/pages", bytes.NewBufferString(pageBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Skipf("could not create test page: %d", createW.Code)
	}
	var page struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(createW.Body.Bytes(), &page); err != nil {
		t.Fatalf("parse page: %v", err)
	}

	// Invalid widget type
	widgetBody := `{"type": "invalid_type", "position": 0}`
	req := httptest.NewRequest(http.MethodPost, "/pages/"+page.ID+"/widgets", bytes.NewBufferString(widgetBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST .../widgets with invalid type: got status %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetPages_RequiresDB(t *testing.T) {
	if db.DB == nil {
		t.Skip("DATABASE_URL not set, skipping DB test")
	}
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/pages", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /pages: got status %d, want %d", w.Code, http.StatusOK)
	}
	// Response should be a JSON array
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
