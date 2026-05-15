package output

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Format(v interface{}, style string) string {
	switch style {
	case "json":
		b, _ := json.Marshal(v)
		return string(b)
	case "pretty":
		b, _ := json.MarshalIndent(v, "", "  ")
		return string(b)
	case "text":
		return textFormat(v)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func textFormat(v interface{}) string {
	switch val := v.(type) {
	case []interface{}:
		var buf bytes.Buffer
		for _, item := range val {
			buf.WriteString(textFormat(item))
			buf.WriteString("\n")
		}
		return buf.String()
	case map[string]interface{}:
		var buf bytes.Buffer
		for k, v := range val {
			fmt.Fprintf(&buf, "%-20s %v\n", k+":", v)
		}
		return buf.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}