// Copyright 2025 Roma Hlushko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		return
	}

	_, _ = fmt.Fprintln(w, string(jsonData))
}
