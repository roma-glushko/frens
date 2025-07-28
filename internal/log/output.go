package log

import (
	"encoding/json"
	"fmt"
	"io"
)

type OutputHandler = func(w io.Writer, data any)

func TextOutputHandler(w io.Writer, data any) {
	_, _ = fmt.Fprintln(w, data)
}

func JSONOutputHandler(w io.Writer, data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error formatting JSON output: %v\n", err)
	}

	_, _ = fmt.Fprintln(w, string(jsonData))
}
