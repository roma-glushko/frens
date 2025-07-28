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
	log.RegisterFormatter(log.FormatText, friend.Location{}, LocationTextFormatter{})
}

type LocationTextFormatter struct{}

var _ log.Formatter = (*LocationTextFormatter)(nil)

func (l LocationTextFormatter) FormatSingle(e any) (string, error) {
	location, ok := e.(friend.Location)

	if !ok {
		return "", fmt.Errorf("expected 'friend.Location'")
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" %s [%s]", labelStyle.Render(location.String()), location.ID))
	sb.WriteString("\n")

	if len(location.Tags) > 0 {
		sb.WriteString(" " + tagStyle.Render(lang.RenderTags(location.Tags)))
		sb.WriteString(" ")
	}

	if location.Desc != "" {
		sb.WriteString("\n")

		wrapped := wrapText(location.Desc, 80)

		for _, line := range wrapped {
			sb.WriteString(" " + line + "\n")
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (l LocationTextFormatter) FormatList(el any) (string, error) {
	locations, ok := el.([]*friend.Location)

	if !ok {
		return "", fmt.Errorf("expected '[]friend.Location'")
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	for _, l := range locations {
		_, _ = fmt.Fprintf(
			w,
			" %s\t%s\t%s\n",
			l.ID,
			labelStyle.Render(l.String()),
			tagStyle.Render(lang.RenderTags(l.Tags)),
		)
	}

	_ = w.Flush()

	return buf.String(), nil
}
