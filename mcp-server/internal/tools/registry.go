package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/0xzcarlet/sregep/mcp-server/internal/mcp"
)

type BackendClient interface {
	Get(ctx context.Context, path string) ([]byte, error)
	Post(ctx context.Context, path string, payload any) ([]byte, error)
}

type Registry struct {
	backend       BackendClient
	defaultUserID string
}

func NewRegistry(backend BackendClient, defaultUserID string) *Registry {
	return &Registry{backend: backend, defaultUserID: defaultUserID}
}

func (r *Registry) Definitions() []mcp.ToolDefinition {
	return []mcp.ToolDefinition{
		definition("finance_add_transaction", "Record income or expense into Sregep.", map[string]any{"type": enum([]string{"income", "expense"}), "amount": num(), "category": str(), "note": str(), "currency": str(), "user_id": str()}, []string{"type", "amount", "category"}),
		definition("finance_list_transactions", "List finance transactions.", map[string]any{"user_id": str()}, []string{}),
		definition("finance_summary", "Get total income, expense, and balance.", map[string]any{"user_id": str()}, []string{}),
		definition("pomodoro_start", "Start a Pomodoro focus session.", map[string]any{"task_name": str(), "duration_minutes": num(), "user_id": str()}, []string{}),
		definition("pomodoro_stop", "Stop a Pomodoro session by ID.", map[string]any{"session_id": str(), "user_id": str()}, []string{"session_id"}),
		definition("pomodoro_current", "Get current running Pomodoro session.", map[string]any{"user_id": str()}, []string{}),
	}
}

func (r *Registry) Call(call mcp.ToolCall) (mcp.ToolResult, error) {
	if call.Arguments == nil {
		call.Arguments = map[string]any{}
	}
	switch call.Name {
	case "finance_add_transaction":
		return r.financeAdd(call.Arguments)
	case "finance_list_transactions":
		return r.financeList(call.Arguments)
	case "finance_summary":
		return r.financeSummary(call.Arguments)
	case "pomodoro_start":
		return r.pomodoroStart(call.Arguments)
	case "pomodoro_stop":
		return r.pomodoroStop(call.Arguments)
	case "pomodoro_current":
		return r.pomodoroCurrent(call.Arguments)
	default:
		return mcp.ToolResult{}, fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (r *Registry) financeAdd(args map[string]any) (mcp.ToolResult, error) {
	userID, err := r.userID(args)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	transactionType, err := stringArg(args, "type")
	if err != nil {
		return mcp.ToolResult{}, err
	}
	amount, err := numberArg(args, "amount")
	if err != nil {
		return mcp.ToolResult{}, err
	}
	category, err := stringArg(args, "category")
	if err != nil {
		return mcp.ToolResult{}, err
	}

	payload := map[string]any{"user_id": userID, "type": transactionType, "amount": amount, "currency": valueOrDefault(args, "currency", "IDR"), "category": category, "note": valueOrDefault(args, "note", ""), "source": "mcp"}
	body, err := r.backend.Post(context.Background(), "/api/transactions", payload)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	return textResult("Transaction saved.\n\n"+pretty(body), false), nil
}

func (r *Registry) financeList(args map[string]any) (mcp.ToolResult, error) {
	userID, err := r.userID(args)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	body, err := r.backend.Get(context.Background(), "/api/transactions?user_id="+url.QueryEscape(userID))
	if err != nil {
		return mcp.ToolResult{}, err
	}
	return textResult(pretty(body), false), nil
}

func (r *Registry) financeSummary(args map[string]any) (mcp.ToolResult, error) {
	userID, err := r.userID(args)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	body, err := r.backend.Get(context.Background(), "/api/summary?user_id="+url.QueryEscape(userID))
	if err != nil {
		return mcp.ToolResult{}, err
	}
	return textResult(pretty(body), false), nil
}

func (r *Registry) pomodoroStart(args map[string]any) (mcp.ToolResult, error) {
	userID, err := r.userID(args)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	duration := 25.0
	if value, ok := args["duration_minutes"]; ok {
		duration, err = numberValue(value)
		if err != nil {
			return mcp.ToolResult{}, err
		}
	}
	payload := map[string]any{"user_id": userID, "task_name": valueOrDefault(args, "task_name", "Focus session"), "duration_minutes": int(duration)}
	body, err := r.backend.Post(context.Background(), "/api/pomodoro/start", payload)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	return textResult("Pomodoro started.\n\n"+pretty(body), false), nil
}

func (r *Registry) pomodoroStop(args map[string]any) (mcp.ToolResult, error) {
	userID, err := r.userID(args)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	sessionID, err := stringArg(args, "session_id")
	if err != nil {
		return mcp.ToolResult{}, err
	}
	body, err := r.backend.Post(context.Background(), "/api/pomodoro/stop", map[string]any{"user_id": userID, "session_id": sessionID})
	if err != nil {
		return mcp.ToolResult{}, err
	}
	return textResult("Pomodoro stopped.\n\n"+pretty(body), false), nil
}

func (r *Registry) pomodoroCurrent(args map[string]any) (mcp.ToolResult, error) {
	userID, err := r.userID(args)
	if err != nil {
		return mcp.ToolResult{}, err
	}
	body, err := r.backend.Get(context.Background(), "/api/pomodoro/current?user_id="+url.QueryEscape(userID))
	if err != nil {
		return mcp.ToolResult{}, err
	}
	return textResult(pretty(body), false), nil
}

func (r *Registry) userID(args map[string]any) (string, error) {
	if userID, ok := args["user_id"].(string); ok && userID != "" {
		return userID, nil
	}
	if r.defaultUserID != "" {
		return r.defaultUserID, nil
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

func textResult(text string, isError bool) mcp.ToolResult {
	return mcp.ToolResult{Content: []mcp.TextContent{{Type: "text", Text: text}}, IsError: isError}
}

func pretty(payload []byte) string {
	var out bytes.Buffer
	if err := json.Indent(&out, payload, "", "  "); err != nil {
		return string(payload)
	}
	return out.String()
}

func valueOrDefault(args map[string]any, key string, fallback string) any {
	value, ok := args[key]
	if !ok || value == nil || value == "" {
		return fallback
	}
	return value
}

func definition(name string, description string, properties map[string]any, required []string) mcp.ToolDefinition {
	return mcp.ToolDefinition{Name: name, Description: description, InputSchema: map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}}
}

func str() map[string]any { return map[string]any{"type": "string"} }
func num() map[string]any { return map[string]any{"type": "number"} }
func enum(values []string) map[string]any { return map[string]any{"type": "string", "enum": values} }
