package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	ID         string     `json:"id,omitempty"`
	UserID     string     `json:"user_id" binding:"required"`
	Type       string     `json:"type" binding:"required"`
	Amount     float64    `json:"amount" binding:"required"`
	Currency   string     `json:"currency"`
	Category   string     `json:"category" binding:"required"`
	Note       string     `json:"note"`
	Source     string     `json:"source"`
	OccurredAt *time.Time `json:"occurred_at,omitempty"`
	CreatedAt  string     `json:"created_at,omitempty"`
}

type PomodoroSession struct {
	ID              string     `json:"id,omitempty"`
	UserID          string     `json:"user_id" binding:"required"`
	TaskName        string     `json:"task_name"`
	Status          string     `json:"status"`
	DurationMinutes int        `json:"duration_minutes"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	CreatedAt       string     `json:"created_at,omitempty"`
}

type StopPomodoroRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
}

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type App struct {
	SupabaseURL string
	SupabaseKey string
	HTTPClient  *http.Client
}

func main() {
	app := App{
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_API_KEY"),
		HTTPClient:  &http.Client{Timeout: 15 * time.Second},
	}

	if app.SupabaseURL == "" || app.SupabaseKey == "" {
		panic("SUPABASE_URL and SUPABASE_API_KEY are required")
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/api/transactions", app.createTransaction)
	r.GET("/api/transactions", app.listTransactions)
	r.GET("/api/summary", app.getSummary)

	r.POST("/api/pomodoro/start", app.startPomodoro)
	r.POST("/api/pomodoro/stop", app.stopPomodoro)
	r.GET("/api/pomodoro/current", app.currentPomodoro)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	_ = r.Run(":" + port)
}

func (app App) createTransaction(c *gin.Context) {
	var trx Transaction
	if err := c.ShouldBindJSON(&trx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if trx.Type != "income" && trx.Type != "expense" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "type must be income or expense"})
		return
	}

	if trx.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "amount must be greater than 0"})
		return
	}

	if trx.Currency == "" {
		trx.Currency = "IDR"
	}

	if trx.Source == "" {
		trx.Source = "api"
	}

	payload, err := json.Marshal(trx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	body, status, err := app.supabaseRequest(http.MethodPost, "/rest/v1/transactions", bytes.NewReader(payload), map[string]string{
		"Content-Type": "application/json",
		"Prefer":       "return=representation",
	})
	if err != nil {
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusCreated, "application/json", body)
}

func (app App) listTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "user_id is required"})
		return
	}

	path := fmt.Sprintf("/rest/v1/transactions?user_id=eq.%s&select=*&order=occurred_at.desc", url.QueryEscape(userID))
	body, status, err := app.supabaseRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

func (app App) getSummary(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "user_id is required"})
		return
	}

	path := fmt.Sprintf("/rest/v1/transactions?user_id=eq.%s&select=type,amount", url.QueryEscape(userID))
	body, status, err := app.supabaseRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	var transactions []Transaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	var summary Summary
	for _, trx := range transactions {
		amount, _ := strconv.ParseFloat(fmt.Sprintf("%v", trx.Amount), 64)
		if trx.Type == "income" {
			summary.TotalIncome += amount
		}
		if trx.Type == "expense" {
			summary.TotalExpense += amount
		}
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense

	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

func (app App) startPomodoro(c *gin.Context) {
	var session PomodoroSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if session.DurationMinutes <= 0 {
		session.DurationMinutes = 25
	}

	if session.Status == "" {
		session.Status = "running"
	}

	if session.StartedAt == nil {
		now := time.Now()
		session.StartedAt = &now
	}

	payload, err := json.Marshal(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	body, status, err := app.supabaseRequest(http.MethodPost, "/rest/v1/pomodoro_sessions", bytes.NewReader(payload), map[string]string{
		"Content-Type": "application/json",
		"Prefer":       "return=representation",
	})
	if err != nil {
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusCreated, "application/json", body)
}

func (app App) stopPomodoro(c *gin.Context) {
	var req StopPomodoroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	now := time.Now()
	payload, err := json.Marshal(map[string]any{
		"status":   "completed",
		"ended_at": now,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	path := fmt.Sprintf("/rest/v1/pomodoro_sessions?id=eq.%s&user_id=eq.%s", url.QueryEscape(req.SessionID), url.QueryEscape(req.UserID))
	body, status, err := app.supabaseRequest(http.MethodPatch, path, bytes.NewReader(payload), map[string]string{
		"Content-Type": "application/json",
		"Prefer":       "return=representation",
	})
	if err != nil {
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

func (app App) currentPomodoro(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "user_id is required"})
		return
	}

	path := fmt.Sprintf("/rest/v1/pomodoro_sessions?user_id=eq.%s&status=eq.running&select=*&order=started_at.desc&limit=1", url.QueryEscape(userID))
	body, status, err := app.supabaseRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

func (app App) supabaseRequest(method string, path string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(method, app.SupabaseURL+path, body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("apikey", app.SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+app.SupabaseKey)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := app.HTTPClient.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if res.StatusCode >= 400 {
		return resBody, res.StatusCode, errors.New(string(resBody))
	}

	return resBody, res.StatusCode, nil
}
