package location

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/tag"
	"github.com/roma-glushko/frens/internal/utils"
)

var ErrNameEmpty = errors.New("location name must be provided")

type Location struct {
	Name    string   `toml:"name"`
	Country string   `toml:"country,omitempty"`
	Alias   []string `toml:"alias,omitempty"`
	Tags    []string `toml:"tags,omitempty"`
}

var _ tag.Tagged = (*Location)(nil)

func (l *Location) Validate() error {
	if l.Name == "" {
		return ErrNameEmpty
	}

	return nil
}

func (l *Location) Match(q string) bool {
	q = strings.ToLower(q)

	if strings.Contains(strings.ToLower(l.Name), q) {
		return true
	}

	for _, a := range l.Alias {
		if strings.Contains(strings.ToLower(a), q) {
			return true
		}
	}

	return false
}

func (l *Location) AddAlias(a string) {
	l.Alias = utils.Unique(append(l.Alias, a))
}

func (l *Location) RemoveAlias(a string) {
	var aliases []string

	for _, alias := range l.Alias {
		if strings.EqualFold(alias, a) {
			aliases = append(aliases, alias)
		}
	}

	l.Alias = aliases
}

func (l *Location) SetTags(tags []string) {
	l.Tags = tags
}

func (l *Location) GetTags() []string {
	return l.Tags
}

func (l *Location) String() string {
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
