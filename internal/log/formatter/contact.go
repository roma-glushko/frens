// Copyright 2026 Roma Hlushko
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

package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/lang"
	"github.com/roma-glushko/frens/internal/log"
)

func init() {
	log.RegisterFormatter(log.FormatText, friend.Contact{}, ContactTextFormatter{})
}

type ContactTextFormatter struct{}

var _ log.Formatter = (*ContactTextFormatter)(nil)

func (f ContactTextFormatter) FormatSingle(e any) (string, error) {
	var c friend.Contact

	switch v := e.(type) {
	case friend.Contact:
		c = v
	case *friend.Contact:
		c = *v
	default:
		return "", ErrInvalidEntity
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s]", idStyle.Render(c.ID)))
	sb.WriteString(" " + labelStyle.Render(string(c.Type)))
	sb.WriteString(": " + c.Value)

	if c.Person != "" {
		sb.WriteString("\n")
		sb.WriteString(" 👤 " + c.Person)
	}

	if len(c.Tags) > 0 {
		sb.WriteString("\n")
		sb.WriteString(" • " + tagStyle.Render(lang.RenderTags(c.Tags)))
	}

	return sb.String(), nil
}

func (f ContactTextFormatter) FormatList(el any) (string, error) {
	contacts, ok := el.([]friend.Contact)

	if !ok {
		return "", ErrInvalidEntity
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, c := range contacts {
		_, _ = fmt.Fprintf(
			w,
			" %s\t%s\t%s\t%s\t%s\n",
			idStyle.Render(c.ID),
			labelStyle.Render(c.Person),
			c.Type,
			c.Value,
			tagStyle.Render(lang.RenderTags(c.Tags)),
		)
	}

	_ = w.Flush()

	return buf.String(), nil
}
