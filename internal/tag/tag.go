package tag

import (
	"strings"

	"github.com/roma-glushko/frens/internal/utils"
)

type Tag struct {
	Name string
}

type Tagged interface {
	SetTags(tags []string)
	GetTags() []string
}

func Add(e Tagged, t string) {
	tags := utils.Unique(append(e.GetTags(), t))

	e.SetTags(tags)
}

func Remove(e Tagged, t string) {
	var tags []string

	for _, tag := range e.GetTags() {
		if strings.EqualFold(tag, t) {
			tags = append(tags, tag)
		}
	}

	e.SetTags(tags)
}

func HasTag(e Tagged, t string) bool {
	for _, tag := range e.GetTags() {
		if strings.EqualFold(tag, t) {
			return true
		}
	}

	return false
}
