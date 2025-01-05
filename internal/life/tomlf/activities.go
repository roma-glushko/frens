package tomlf

import (
	"github.com/roma-glushko/frens/internal/activity"
)

const FileNameActivities = "activities.toml"

type ActivitiesFile struct {
	Activities []activity.Activity `toml:"activities"`
}
