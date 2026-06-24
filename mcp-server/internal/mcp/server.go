package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
)

const protocolVersion = "2024-11-05"

type Registry interface {
	Definitions() []ToolDefinition
	Call(call ToolCall) (ToolResult, error)
}

type Server struct {
	name     string
	version  string
	registry Registry
	logger   *slog.Logger
}

func NewServer(name string, version string, registry Registry, logger *slog.Logger) *Server {
	return &Server{name: name, version: version, registry: registry, logger: logger}
}

func (s *Server) Run(input io.Reader, output io.Writer) {
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

func (s *Server) Handle(req Request) *Response {
	if req.ID == nil && strings.HasPrefix(req.Method, "notifications/") {
		return nil
	}

	res := &Response{JSONRPC: "2.0", ID: req.ID}

	switch req.Method {
	case "initialize":
		res.Result = map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities":     map[string]any{"tools": map[string]any{}},
			"serverInfo":       map[string]any{"name": s.name, "version": s.version},
		}
	case "ping":
		res.Result = map[string]any{}
	case "tools/list":
		res.Result = map[string]any{"tools": s.registry.Definitions()}
	case "tools/call":
		result, err := s.callTool(req.Params)
		if err != nil {
			res.Result = ToolResult{Content: []TextContent{{Type: "text", Text: err.Error()}}, IsError: true}
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

func (s *Server) callTool(params json.RawMessage) (ToolResult, error) {
	var call ToolCall
	if err := json.Unmarshal(params, &call); err != nil {
		return ToolResult{}, err
	}
	if call.Arguments == nil {
		call.Arguments = map[string]any{}
	}
	return s.registry.Call(call)
}

func (s *Server) write(writer *bufio.Writer, res Response) {
	payload, err := json.Marshal(res)
	if err != nil {
		s.logger.Error("failed to encode mcp response", "error", err)
		return
	}
	_, _ = writer.Write(payload)
	_, _ = writer.WriteString("\n")
	_ = writer.Flush()
}
