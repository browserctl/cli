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
	addr  string
	conn  *websocket.Conn
	mu    sync.Mutex
	seq   int64
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

	return &CdpClient{addr: addr, conn: conn}, nil
}

func (c *CdpClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *CdpClient) nextSeq() int64 {
	c.seq++
	return c.seq
}

type JsonRpcRequest struct {
	ID        int64          `json:"id,omitempty"`
	Method    string         `json:"method"`
	SessionId string        `json:"sessionId,omitempty"`
	Params    json.RawMessage `json:"params,omitempty"`
}

type JsonRpcResponse struct {
	ID     int64          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *JsonRpcError  `json:"error,omitempty"`
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
func (c *CdpClient) Send(ctx context.Context, method string, sessionId string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	id := c.nextSeq()
	c.mu.Unlock()

	req := JsonRpcRequest{
		ID:        id,
		Method:    method,
		SessionId: sessionId,
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

type TabsResult struct {
	TargetInfos []TargetInfo `json:"targetInfos"`
}

func (c *CdpClient) GetTabs(ctx context.Context) ([]TargetInfo, error) {
	result, err := c.Send(ctx, "Target.getTargets", "", nil)
	if err != nil {
		return nil, err
	}
	var tabsResult TabsResult
	if err := json.Unmarshal(result, &tabsResult); err != nil {
		return nil, fmt.Errorf("unmarshal get targets: %w", err)
	}
	return tabsResult.TargetInfos, nil
}

func (c *CdpClient) Attach(ctx context.Context, tabId string) (string, error) {
	params := map[string]interface{}{
		"targetId": tabId,
		"flatten":  true,
	}
	result, err := c.Send(ctx, "Target.attachToTarget", "", params)
	if err != nil {
		return "", err
	}
	var attachResult struct {
		SessionId string `json:"sessionId"`
	}
	if err := json.Unmarshal(result, &attachResult); err != nil {
		return "", fmt.Errorf("unmarshal attach result: %w", err)
	}
	return attachResult.SessionId, nil
}

func (c *CdpClient) Navigate(ctx context.Context, sessionId, url string) error {
	tabId := ParseSessionId(sessionId)
	if tabId == 0 {
		return fmt.Errorf("invalid sessionId: %s", sessionId)
	}
	attachParams := map[string]interface{}{
		"targetId": fmt.Sprintf("tab-%d", tabId),
		"flatten":  true,
	}
	attachResult, err := c.Send(ctx, "Target.attachToTarget", "", attachParams)
	if err != nil {
		return fmt.Errorf("attach: %w", err)
	}
	var attach struct {
		SessionId string `json:"sessionId"`
	}
	if err := json.Unmarshal(attachResult, &attach); err != nil {
		return fmt.Errorf("unmarshal attach: %w", err)
	}

	navParams := map[string]interface{}{"url": url}
	_, err = c.Send(ctx, "Page.navigate", attach.SessionId, navParams)
	return err
}

func (c *CdpClient) Eval(ctx context.Context, sessionId, expr string) (interface{}, error) {
	params := map[string]interface{}{"expression": expr}
	result, err := c.Send(ctx, "Runtime.evaluate", sessionId, params)
	if err != nil {
		return nil, err
	}
	var evalResult struct {
		Result struct {
			Type  string      `json:"type"`
			Value interface{} `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(result, &evalResult); err != nil {
		return nil, fmt.Errorf("unmarshal eval result: %w", err)
	}
	return evalResult.Result.Value, nil
}

func (c *CdpClient) Click(ctx context.Context, sessionId, selector string) error {
	script := fmt.Sprintf(`(function(){var el=document.querySelector('%s'); if(!el) throw new Error('element not found: '+'%s'); el.click();})()`, selector, selector)
	_, err := c.Eval(ctx, sessionId, script)
	return err
}

func (c *CdpClient) Fill(ctx context.Context, sessionId, selector, value string) error {
	script := fmt.Sprintf(`(function(){var el=document.querySelector('%s'); if(!el) throw new Error('element not found'); el.value='%s'; el.dispatchEvent(new Event('input',{bubbles:true})); el.dispatchEvent(new Event('change',{bubbles:true}));})()`, selector, value)
	_, err := c.Eval(ctx, sessionId, script)
	return err
}

func (c *CdpClient) Find(ctx context.Context, sessionId, selector string) error {
	script := fmt.Sprintf(`(function(){if(document.querySelector('%s')) return true; return new Promise(resolve=>{var obs=new MutationObserver(ms=>{if(document.querySelector('%s')){obs.disconnect();resolve(true);}});obs.observe(document.body,{childList:true,subtree:true});});})()`, selector, selector)
	_, err := c.Eval(ctx, sessionId, script)
	return err
}

func (c *CdpClient) Hover(ctx context.Context, sessionId, selector string) error {
	script := fmt.Sprintf(`(function(){var el=document.querySelector('%s'); if(!el) throw new Error('element not found'); el.dispatchEvent(new MouseEvent('mouseover',{bubbles:true}));})()`, selector)
	_, err := c.Eval(ctx, sessionId, script)
	return err
}

func (c *CdpClient) Type(ctx context.Context, sessionId, selector, text string) error {
	script := fmt.Sprintf(`(function(){var el=document.querySelector('%s'); if(!el) throw new Error('element not found'); el.focus(); el.value=''; el.dispatchEvent(new Event('input',{bubbles:true}));})()`, selector)
	_, err := c.Eval(ctx, sessionId, script)
	if err != nil {
		return err
	}
	script2 := fmt.Sprintf(`(function(){var el=document.querySelector('%s'); el.value='%s'; el.dispatchEvent(new Event('input',{bubbles:true})); el.dispatchEvent(new Event('change',{bubbles:true}));})()`, selector, text)
	_, err = c.Eval(ctx, sessionId, script2)
	return err
}

func (c *CdpClient) Scroll(ctx context.Context, sessionId string, px int) error {
	script := fmt.Sprintf("window.scrollBy(0, %d)", px)
	_, err := c.Eval(ctx, sessionId, script)
	return err
}

func (c *CdpClient) Screenshot(ctx context.Context, sessionId string) ([]byte, error) {
	result, err := c.Send(ctx, "Page.captureScreenshot", sessionId, map[string]interface{}{"format": "png"})
	if err != nil {
		return nil, err
	}
	var screenshot struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(result, &screenshot); err != nil {
		return nil, fmt.Errorf("unmarshal screenshot: %w", err)
	}
	raw, err := decodeBase64(screenshot.Data)
	if err != nil {
		return []byte(screenshot.Data), nil
	}
	return raw, nil
}

func decodeBase64(s string) ([]byte, error) {
	return []byte(s), nil
}

func (c *CdpClient) CloseTab(ctx context.Context, sessionId string) error {
	tabId := ParseSessionId(sessionId)
	if tabId == 0 {
		return fmt.Errorf("invalid sessionId: %s", sessionId)
	}
	params := map[string]interface{}{
		"targetId": fmt.Sprintf("tab-%d", tabId),
	}
	_, err := c.Send(ctx, "Target.closeTarget", "", params)
	return err
}

func (c *CdpClient) NewTab(ctx context.Context, url string) (string, error) {
	params := map[string]interface{}{"url": url}
	result, err := c.Send(ctx, "Target.createTarget", "", params)
	if err != nil {
		return "", err
	}
	var createResult struct {
		TargetId string `json:"targetId"`
	}
	if err := json.Unmarshal(result, &createResult); err != nil {
		return "", fmt.Errorf("unmarshal createTarget: %w", err)
	}
	return createResult.TargetId, nil
}

func (c *CdpClient) GetUrl(ctx context.Context, sessionId string) (string, error) {
	val, err := c.Eval(ctx, sessionId, "window.location.href")
	if err != nil {
		return "", err
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("unexpected result type: %T", val)
}

func (c *CdpClient) Back(ctx context.Context, sessionId string) error {
	_, err := c.Eval(ctx, sessionId, "history.back()")
	return err
}

func (c *CdpClient) Forward(ctx context.Context, sessionId string) error {
	_, err := c.Eval(ctx, sessionId, "history.forward()")
	return err
}

func (c *CdpClient) Reload(ctx context.Context, sessionId string) error {
	_, err := c.Send(ctx, "Page.reload", sessionId, map[string]interface{}{"ignoreCache": false})
	return err
}

func (c *CdpClient) Html(ctx context.Context, sessionId string) (string, error) {
	val, err := c.Eval(ctx, sessionId, "document.documentElement.outerHTML")
	if err != nil {
		return "", err
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("unexpected result type: %T", val)
}

func (c *CdpClient) Cookies(ctx context.Context, sessionId string) (string, error) {
	val, err := c.Eval(ctx, sessionId, "document.cookie")
	if err != nil {
		return "", err
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("unexpected result type: %T", val)
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