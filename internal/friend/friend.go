package friend

import "errors"

var (
	ErrNameEmpty = errors.New("friend name must be provided")
)

type Friend struct {
	Name      string
	Nicknames []string
	Tags      []string
	Locations []string
	Reminders []string
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
