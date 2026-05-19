# browserctl CLI

Command-line interface for controlling Chrome browser via CDP (Chrome DevTools Protocol).

## Installation

```bash
go install ./...
```

## Quick Start

```bash
# Launch browser with browserctl service
browserctl launch

# Navigate to a URL
browserctl navigate https://example.com

# Take a screenshot
browserctl screenshot -o image.png

# Click an element
browserctl click "#submit-button"

# Execute JavaScript
browserctl eval "document.title"
```

## Commands

| Command | Description |
|---------|-------------|
| `launch` | Launch Chrome browser |
| `navigate` | Navigate to URL |
| `back` / `forward` | Browser history navigation |
| `click` | Click element by selector |
| `hover` | Hover over element |
| `fill` | Fill input field |
| `typeinput` | Type text character by character |
| `scroll` | Scroll page or element |
| `screenshot` | Take screenshot |
| `html` | Get page HTML |
| `eval` | Execute JavaScript |
| `find` | Find element |
| `switch` | Switch tab or window |
| `tabs` | List all tabs |
| `url` | Get current URL |
| `cookies` | Manage cookies |
| `attach` | Attach to existing Chrome |
| `close` | Close tab or browser |
| `reload` | Reload page |
| `version` | Show version |

## Global Options

```bash
-s, --svc string     browserctl svc address (default "ws://localhost:9222")
--secret string      auth secret
-o, --output string  output format: json, text, pretty (default "json")
-t, --timeout int    timeout in milliseconds (default 30000)
--config string     config file (default "./browserctl.yaml")
```

## Configuration

Create `~/.config/browserctl/browserctl.yaml`:

```yaml
svc: "ws://localhost:9222"
secret: "your-secret"
output: "pretty"
timeout: 30000
```

## Architecture

```
browserctl CLI → browserctl service → Chrome extension → Chrome browser
```

The CLI connects to the browserctl service via WebSocket. The service manages the Chrome extension and forwards commands.

## See Also

- [browserctl/svc](https://github.com/browserctl/svc) - Service documentation
- [browserctl/ext](https://github.com/browserctl/ext) - Chrome extension documentation