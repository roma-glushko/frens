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
	"fmt"

	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/reminder"
)

func init() {
	f := ReminderSyncFormatter{}
	log.RegisterFormatter(log.FormatText, reminder.SyncResult{}, f)
	log.RegisterFormatter(log.FormatJSON, reminder.SyncResult{}, f)
	log.RegisterFormatter(log.FormatMarkdown, reminder.SyncResult{}, f)
}

type ReminderSyncFormatter struct{}

var _ log.Formatter = (*ReminderSyncFormatter)(nil)

func (f ReminderSyncFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	r, ok := e.(reminder.SyncResult)
	if !ok {
		return "", ErrInvalidEntity
	}

	if r.Err != nil {
		return log.WarnStyle.Render(fmt.Sprintf("reminder: %v", r.Err)) + "\n", nil
	}

	switch r.Action {
	case reminder.SyncCreated:
		return log.SuccessStyle.Render(" ✔") + " " +
			fmt.Sprintf("Reminder set (triggers %s)", r.Reminder.TriggerAt.Format("2006-01-02")) +
			"\n", nil
	case reminder.SyncRemoved:
		return log.SuccessStyle.Render(" ✔") + " Reminder removed.\n", nil
	default:
		return "", nil
	}
}

func (f ReminderSyncFormatter) FormatList(_ log.FormatterContext, _ any) (string, error) {
	return "", ErrInvalidEntity
}
