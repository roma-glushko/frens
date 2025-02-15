package friend

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNameEmpty = errors.New("friend name must be provided")
)

type Friend struct {
	Name      string    `toml:"name"`
	Birthday  time.Time `toml:"birthday,omitempty"`
	Nicknames []string  `toml:"nicknames,omitempty"`
	Tags      []string  `toml:"tags,omitempty"`
	Locations []string  `toml:"locations,omitempty"`
	Reminders []string  `toml:"reminders,omitempty"`
}

func (f *Friend) Validate() error {
	if f.Name == "" {
		return ErrNameEmpty
	}

	return nil
}

func (f *Friend) AddNickname(n string) {
	f.Nicknames = append(f.Nicknames, n)
}

func (f *Friend) AddTag(t string) {
	f.Tags = append(f.Tags, t)
}

func (f *Friend) AddLocation(l string) {
	f.Locations = append(f.Locations, l)
}

func (f *Friend) String() string {
	var sb strings.Builder

	sb.WriteString(f.Name)

	if len(f.Nicknames) > 0 {
		sb.WriteString(fmt.Sprintf(" (a.k.a %s)", strings.Join(f.Nicknames, ", ")))
	}

	return sb.String()
}
