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

package notify

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/roma-glushko/frens/internal/config"
	"github.com/roma-glushko/frens/internal/friend"
)

// TemplateContext provides data for template rendering
type TemplateContext struct {
	EntityType   string // "date", "wishlist", "activity", "note"
	Reminder     *friend.Reminder
	Event        *friend.Event
	Date         *friend.Date
	WishlistItem *friend.WishlistItem
	Friend       *friend.Person
	DaysUntil    int // Days until trigger date
	Now          time.Time
}

// NewTemplateContext creates a new template context for a reminder
func NewTemplateContext(
	rc *ReminderContext,
	now time.Time,
) *TemplateContext {
	return &TemplateContext{
		Reminder:     rc.Reminder,
		Date:         rc.Date,
		WishlistItem: rc.WishlistItem,
		Event:        rc.Event,
		Friend:       rc.Friend,
		Now:          now,
	}
}

// Default templates for different entity types
var defaultTemplates = map[string]string{
	"date": `{{ if .Friend }}{{ .Friend.Name }}'s {{ end }}{{ if .LinkedEntity.Desc }}{{ .LinkedEntity.Desc }}{{ else }}important date{{ end }}
{{ if eq .DaysUntil 0 }}Today!{{ else if eq .DaysUntil 1 }}Tomorrow{{ else }}In {{ .DaysUntil }} days{{ end }}
Date: {{ .LinkedEntity.DateExpr }}`,

	"wishlist": `{{ if .Friend }}Gift idea for {{ .Friend.Name }}{{ else }}Wishlist reminder{{ end }}
{{ .LinkedEntity.Desc }}{{ if .LinkedEntity.Link }}
Link: {{ .LinkedEntity.Link }}{{ end }}{{ if .LinkedEntity.Price }}
Price: {{ .LinkedEntity.Price }}{{ end }}`,

	"activity": `Activity reminder{{ if .Friend }} ({{ .Friend.Name }}){{ end }}
{{ .LinkedEntity.Desc }}`,

	"note": `Note reminder{{ if .Friend }} ({{ .Friend.Name }}){{ end }}
{{ .LinkedEntity.Desc }}`,
}

// RenderTemplate renders a notification message using the provided template
func RenderTemplate(tmpl *config.NotificationTemplate, ctx *TemplateContext) (string, error) {
	var templateBody string

	if tmpl != nil && tmpl.Body != "" {
		templateBody = tmpl.Body
	} else {
		// Fall back to default template based on entity type
		defaultTmpl, ok := defaultTemplates[ctx.EntityType]
		if !ok {
			defaultTmpl = defaultTemplates["note"]
		}

		templateBody = defaultTmpl
	}

	return executeTemplate(templateBody, ctx)
}

// RenderSubject renders the subject line for email notifications
func RenderSubject(tmpl *config.NotificationTemplate, ctx *TemplateContext) (string, error) {
	if tmpl == nil || tmpl.Subject == "" {
		return generateDefaultSubject(ctx), nil
	}

	return executeTemplate(tmpl.Subject, ctx)
}

func generateDefaultSubject(ctx *TemplateContext) string {
	switch ctx.EntityType {
	case "date":
		if ctx.Friend != nil {
			return "Reminder: " + ctx.Friend.Name + "'s date"
		}

		return "Date reminder"
	case "wishlist":
		if ctx.Friend != nil {
			return "Gift reminder for " + ctx.Friend.Name
		}

		return "Wishlist reminder"
	default:
		return "Frens reminder"
	}
}

func executeTemplate(templateStr string, ctx *TemplateContext) (string, error) {
	funcMap := template.FuncMap{
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"title":    strings.Title, //nolint:staticcheck
		"trim":     strings.TrimSpace,
		"contains": strings.Contains,
		"join":     strings.Join,
		"formatDate": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
		"now": func() time.Time {
			return ctx.Now
		},
	}

	tmpl, err := template.New("notification").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// GetTemplateForReminder selects the appropriate template for a reminder
func GetTemplateForReminder(
	notifications *config.Notifications,
	rule *config.RoutingRule,
	entityType string,
) *config.NotificationTemplate {
	// First check if rule specifies a template
	if rule != nil && rule.TemplateID != "" {
		if tmpl := notifications.GetTemplate(rule.TemplateID); tmpl != nil {
			return tmpl
		}
	}

	// Then check for entity-type specific template
	entityTemplateID := entityType + "-default"
	if tmpl := notifications.GetTemplate(entityTemplateID); tmpl != nil {
		return tmpl
	}

	// Fall back to nil (will use built-in default)
	return nil
}
