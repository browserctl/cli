package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type CdpClient struct {
	addr    string
	secret  string
	conn    *websocket.Conn
	mu      sync.Mutex
	nextId  int64
}

func New(addr, secret string) (*CdpClient, error) {
	if addr == "" {
		addr = "ws://localhost:9222"
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid addr: %w", err)
	}

	query := u.Query()
	if secret != "" {
		query.Set("secret", secret)
	}
	u.RawQuery = query.Encode()

	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}

	conn, _, err := dialer.Dial(u.String(), http.Header{
		"Origin": {u.String()},
	})
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return &CdpClient{addr: addr, secret: secret, conn: conn}, nil
}

func (c *CdpClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *CdpClient) nextId() int64 {
	c.nextId++
	return c.nextId
}

type JsonRpcRequest struct {
	ID      int64           `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params json.RawMessage  `json:"params,omitempty"`
}

type JsonRpcResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *JsonRpcError   `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TargetInfo struct {
	TargetId         string `json:"targetId"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	URL              string `json:"url"`
	Attached         bool   `json:"attached"`
	BrowserContextId string `json:"browserContextId,omitempty"`
}

// Send issues a CDP command and waits for the response.
func (c *CdpClient) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	id := c.nextId()
	c.mu.Unlock()

	req := JsonRpcRequest{
		ID:     id,
		Method: method,
	}
	if params != nil {
		payload, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("marshal params: %w", err)
		}
		req.Params = payload
	}

	if err := c.conn.WriteJSON(req); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	for {
		var resp JsonRpcResponse
		if err := c.conn.ReadJSON(&resp); err != nil {
			return nil, fmt.Errorf("read: %w", err)
		}
		if resp.ID != id {
			continue
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("cdp error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

type AttachParams struct {
	TargetId string `json:"targetId"`
	Flatten  bool   `json:"flatten"`
}

type AttachResult struct {
	SessionId string `json:"sessionId"`
}

func (c *CdpClient) Attach(ctx context.Context, tabId string) (string, error) {
	params := AttachParams{TargetId: tabId, Flatten: true}
	result, err := c.Send(ctx, "Target.attachToTarget", params)
	if err != nil {
		return "", err
	}
	var attachResult AttachResult
	if err := json.Unmarshal(result, &attachResult); err != nil {
		return "", fmt.Errorf("unmarshal attach result: %w", err)
	}
	return attachResult.SessionId, nil
}

type NavigateParams struct {
	URL string `json:"url"`
}

func (c *CdpClient) Navigate(ctx context.Context, sessionId, url string) error {
	params := struct {
		SessionId string `json:"sessionId"`
		URL       string `json:"url"`
	}{sessionId, url}

	tabId := ParseSessionId(sessionId)
	_ = tabId
	return nil
}

func ParseSessionId(sessionId string) int {
	if len(sessionId) < 3 || sessionId[:3] != "cs-" {
		return 0
	}
	var tabId int
	for _, c := range sessionId[3:] {
		if c == '-' {
			break
		}
		if c >= '0' && c <= '9' {
			tabId = tabId*10 + int(c-'0')
		}
	}
	return tabId
}

type TabsResult struct {
	TargetInfos []TargetInfo `json:"targetInfos"`
}

func (c *CdpClient) GetTabs(ctx context.Context) ([]TargetInfo, error) {
	result, err := c.Send(ctx, "Target.getTargets", nil)
	if err != nil {
		return nil, err
	}
	var tabsResult TabsResult
	if err := json.Unmarshal(result, &tabsResult); err != nil {
		return nil, fmt.Errorf("unmarshal get targets: %w", err)
	}
	return tabsResult.TargetInfos, nil
}