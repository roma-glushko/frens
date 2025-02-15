package location

import (
	"fmt"
	"strings"
)

type Location struct {
	Name    string   `toml:"name"`
	Country string   `toml:"country,omitempty"`
	Alias   []string `toml:"alias,omitempty"`
	Tags    []string `toml:"tags,omitempty"`
}

func (l Location) String() string {
	var sb strings.Builder

	sb.WriteString(l.Name)

	if len(l.Country) > 0 {
		sb.WriteString(fmt.Sprintf(", %s", l.Country))
	}

	if len(l.Alias) > 0 {
		sb.WriteString(fmt.Sprintf(" (%s)", strings.Join(l.Alias, ", ")))
	}

	return sb.String()
}
