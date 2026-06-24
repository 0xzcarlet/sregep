package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	serverName        = "sregep-mcp"
	serverVersion     = "0.1.0"
	protocolVersion   = "2024-11-05"
	defaultAPIBaseURL = "http://localhost:8080"
)

type Server struct {
	APIBaseURL    string
	DefaultUserID string
	HTTPClient    *http.Client
}

type Request struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  any              `json:"result,omitempty"`
	Error   *RPCError        `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolResult struct {
	Content []TextContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	apiBaseURL := strings.TrimRight(os.Getenv("SREGEP_API_BASE_URL"), "/")
	if apiBaseURL == "" {
		apiBaseURL = defaultAPIBaseURL
	}

	server := Server{
		APIBaseURL:    apiBaseURL,
		DefaultUserID: os.Getenv("SREGEP_DEFAULT_USER_ID"),
		HTTPClient:    &http.Client{Timeout: 20 * time.Second},
	}

	server.Run(os.Stdin, os.Stdout)
}

func (s Server) Run(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.write(writer, Response{JSONRPC: "2.0", Error: &RPCError{Code: -32700, Message: err.Error()}})
			continue
		}

		res := s.Handle(req)
		if res != nil {
			s.write(writer, *res)
		}
	}
}

func (s Server) Handle(req Request) *Response {
	if req.ID == nil && strings.HasPrefix(req.Method, "notifications/") {
		return nil
	}

	res := &Response{JSONRPC: "2.0", ID: req.ID}

	switch req.Method {
	case "initialize":
		res.Result = map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities":     map[string]any{"tools": map[string]any{}},
			"serverInfo":       map[string]any{"name": serverName, "version": serverVersion},
		}
	case "ping":
		res.Result = map[string]any{}
	case "tools/list":
		res.Result = map[string]any{"tools": s.Tools()}
	case "tools/call":
		result, err := s.CallTool(req.Params)
		if err != nil {
			res.Result = resultText(err.Error(), true)
		} else {
			res.Result = result
		}
	case "resources/list":
		res.Result = map[string]any{"resources": []any{}}
	case "prompts/list":
		res.Result = map[string]any{"prompts": []any{}}
	default:
		res.Error = &RPCError{Code: -32601, Message: "method not found: " + req.Method}
	}

	return res
}

func (s Server) write(writer *bufio.Writer, res Response) {
	payload, err := json.Marshal(res)
	if err != nil {
		return
	}
	_, _ = writer.Write(payload)
	_, _ = writer.WriteString("\n")
	_ = writer.Flush()
}

func (s Server) CallTool(params json.RawMessage) (ToolResult, error) {
	var call ToolCall
	if err := json.Unmarshal(params, &call); err != nil {
		return ToolResult{}, fmt.Errorf("invalid tools/call params: %w", err)
	}
	if call.Arguments == nil {
		call.Arguments = map[string]any{}
	}

	switch call.Name {
	case "finance_add_transaction":
		return s.financeAdd(call.Arguments)
	case "finance_list_transactions":
		return s.financeList(call.Arguments)
	case "finance_summary":
		return s.financeSummary(call.Arguments)
	case "pomodoro_start":
		return s.pomodoroStart(call.Arguments)
	case "pomodoro_stop":
		return s.pomodoroStop(call.Arguments)
	case "pomodoro_current":
		return s.pomodoroCurrent(call.Arguments)
	default:
		return ToolResult{}, fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (s Server) financeAdd(args map[string]any) (ToolResult, error) {
	userID, err := s.userID(args)
	if err != nil {
		return ToolResult{}, err
	}

	transactionType, err := stringArg(args, "type")
	if err != nil {
		return ToolResult{}, err
	}
	if transactionType != "income" && transactionType != "expense" {
		return ToolResult{}, errors.New("type must be income or expense")
	}

	amount, err := numberArg(args, "amount")
	if err != nil {
		return ToolResult{}, err
	}

	category, err := stringArg(args, "category")
	if err != nil {
		return ToolResult{}, err
	}

	payload := map[string]any{
		"user_id":  userID,
		"type":     transactionType,
		"amount":   amount,
		"currency": fallback(args, "currency", "IDR"),
		"category": category,
		"note":     fallback(args, "note", ""),
		"source":   "mcp",
	}

	if occurredAt, ok := args["occurred_at"].(string); ok && occurredAt != "" {
		payload["occurred_at"] = occurredAt
	}

	body, err := s.api(http.MethodPost, "/api/transactions", payload)
	if err != nil {
		return ToolResult{}, err
	}
	return resultText("Transaction saved.\n\n"+pretty(body), false), nil
}

func (s Server) financeList(args map[string]any) (ToolResult, error) {
	userID, err := s.userID(args)
	if err != nil {
		return ToolResult{}, err
	}
	body, err := s.api(http.MethodGet, "/api/transactions?user_id="+userID, nil)
	if err != nil {
		return ToolResult{}, err
	}
	return resultText(pretty(body), false), nil
}

func (s Server) financeSummary(args map[string]any) (ToolResult, error) {
	userID, err := s.userID(args)
	if err != nil {
		return ToolResult{}, err
	}
	body, err := s.api(http.MethodGet, "/api/summary?user_id="+userID, nil)
	if err != nil {
		return ToolResult{}, err
	}
	return resultText(pretty(body), false), nil
}

func (s Server) pomodoroStart(args map[string]any) (ToolResult, error) {
	userID, err := s.userID(args)
	if err != nil {
		return ToolResult{}, err
	}

	duration := 25.0
	if value, ok := args["duration_minutes"]; ok {
		duration, err = numberValue(value)
		if err != nil {
			return ToolResult{}, errors.New("duration_minutes must be a number")
		}
	}

	payload := map[string]any{
		"user_id":          userID,
		"task_name":        fallback(args, "task_name", "Focus session"),
		"duration_minutes": int(duration),
		"status":           "running",
	}

	body, err := s.api(http.MethodPost, "/api/pomodoro/start", payload)
	if err != nil {
		return ToolResult{}, err
	}
	return resultText("Pomodoro started.\n\n"+pretty(body), false), nil
}

func (s Server) pomodoroStop(args map[string]any) (ToolResult, error) {
	userID, err := s.userID(args)
	if err != nil {
		return ToolResult{}, err
	}
	sessionID, err := stringArg(args, "session_id")
	if err != nil {
		return ToolResult{}, err
	}

	body, err := s.api(http.MethodPost, "/api/pomodoro/stop", map[string]any{"user_id": userID, "session_id": sessionID})
	if err != nil {
		return ToolResult{}, err
	}
	return resultText("Pomodoro stopped.\n\n"+pretty(body), false), nil
}

func (s Server) pomodoroCurrent(args map[string]any) (ToolResult, error) {
	userID, err := s.userID(args)
	if err != nil {
		return ToolResult{}, err
	}
	body, err := s.api(http.MethodGet, "/api/pomodoro/current?user_id="+userID, nil)
	if err != nil {
		return ToolResult{}, err
	}
	return resultText(pretty(body), false), nil
}

func (s Server) api(method, path string, payload any) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequest(method, s.APIBaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("backend returned HTTP %d: %s", res.StatusCode, string(resBody))
	}
	return resBody, nil
}

func (s Server) userID(args map[string]any) (string, error) {
	if userID, ok := args["user_id"].(string); ok && userID != "" {
		return userID, nil
	}
	if s.DefaultUserID != "" {
		return s.DefaultUserID, nil
	}
	return "", errors.New("user_id is required, or set SREGEP_DEFAULT_USER_ID")
}

func stringArg(args map[string]any, key string) (string, error) {
	value, ok := args[key]
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	text, ok := value.(string)
	if !ok || text == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return text, nil
}

func numberArg(args map[string]any, key string) (float64, error) {
	value, ok := args[key]
	if !ok {
		return 0, fmt.Errorf("%s is required", key)
	}
	return numberValue(value)
}

func numberValue(value any) (float64, error) {
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case int:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case json.Number:
		return typed.Float64()
	default:
		return 0, errors.New("value must be a number")
	}
}

func fallback(args map[string]any, key, fallback string) any {
	value, ok := args[key]
	if !ok || value == nil || value == "" {
		return fallback
	}
	return value
}

func pretty(payload []byte) string {
	var out bytes.Buffer
	if err := json.Indent(&out, payload, "", "  "); err != nil {
		return string(payload)
	}
	return out.String()
}

func resultText(text string, isError bool) ToolResult {
	return ToolResult{Content: []TextContent{{Type: "text", Text: text}}, IsError: isError}
}

func (s Server) Tools() []map[string]any {
	return []map[string]any{
		{"name": "finance_add_transaction", "description": "Record income or expense into Sregep.", "inputSchema": schema(map[string]any{"user_id": str("User UUID. Optional when SREGEP_DEFAULT_USER_ID is set."), "type": enum([]string{"income", "expense"}, "Transaction type."), "amount": num("Amount."), "currency": str("Currency. Default: IDR."), "category": str("Category."), "note": str("Optional note."), "occurred_at": str("Optional ISO timestamp.")}, []string{"type", "amount", "category"})},
		{"name": "finance_list_transactions", "description": "List finance transactions.", "inputSchema": schema(map[string]any{"user_id": str("User UUID. Optional when SREGEP_DEFAULT_USER_ID is set.")}, []string{})},
		{"name": "finance_summary", "description": "Get total income, expense, and balance.", "inputSchema": schema(map[string]any{"user_id": str("User UUID. Optional when SREGEP_DEFAULT_USER_ID is set.")}, []string{})},
		{"name": "pomodoro_start", "description": "Start a Pomodoro focus session.", "inputSchema": schema(map[string]any{"user_id": str("User UUID. Optional when SREGEP_DEFAULT_USER_ID is set."), "task_name": str("Task name."), "duration_minutes": num("Duration in minutes. Default: 25.")}, []string{})},
		{"name": "pomodoro_stop", "description": "Stop a Pomodoro session by session id.", "inputSchema": schema(map[string]any{"user_id": str("User UUID. Optional when SREGEP_DEFAULT_USER_ID is set."), "session_id": str("Pomodoro session UUID.")}, []string{"session_id"})},
		{"name": "pomodoro_current", "description": "Get current running Pomodoro session.", "inputSchema": schema(map[string]any{"user_id": str("User UUID. Optional when SREGEP_DEFAULT_USER_ID is set.")}, []string{})},
	}
}

func schema(properties map[string]any, required []string) map[string]any {
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}

func str(description string) map[string]any {
	return map[string]any{"type": "string", "description": description}
}

func num(description string) map[string]any {
	return map[string]any{"type": "number", "description": description}
}

func enum(values []string, description string) map[string]any {
	return map[string]any{"type": "string", "enum": values, "description": description}
}
