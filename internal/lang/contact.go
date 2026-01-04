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

package lang

import (
	"strings"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/tag"
)

// Contact type aliases for shorthand prefixes
var contactTypeAliases = map[string]friend.ContactType{
	"email":     friend.ContactTypeEmail,
	"mail":      friend.ContactTypeEmail,
	"phone":     friend.ContactTypePhone,
	"tel":       friend.ContactTypePhone,
	"tg":        friend.ContactTypeTelegram,
	"telegram":  friend.ContactTypeTelegram,
	"wa":        friend.ContactTypeWhatsApp,
	"whatsapp":  friend.ContactTypeWhatsApp,
	"signal":    friend.ContactTypeSignal,
	"x":         friend.ContactTypeTwitter,
	"twitter":   friend.ContactTypeTwitter,
	"li":        friend.ContactTypeLinkedIn,
	"linkedin":  friend.ContactTypeLinkedIn,
	"gh":        friend.ContactTypeGitHub,
	"github":    friend.ContactTypeGitHub,
	"ig":        friend.ContactTypeInstagram,
	"instagram": friend.ContactTypeInstagram,
	"fb":        friend.ContactTypeFacebook,
	"facebook":  friend.ContactTypeFacebook,
	"discord":   friend.ContactTypeDiscord,
	"slack":     friend.ContactTypeSlack,
}

// FormatContactInfo describes the expected input format for contacts
var FormatContactInfo = "[TYPE:]VALUE [...] [#tags] - e.g., ig:@handle +48123456789 x:@user name@example.com #work"

// ExtractContacts parses multiple contacts from a single input string.
// Supported formats:
//   - type:value (e.g., ig:@instagram_handle, x:@twitter_user, tg:@telegram)
//   - email addresses are auto-detected (contains @ and .)
//   - phone numbers are auto-detected (starts with + or digits)
//   - bare values default to "other" type
//   - tags at the end apply to all contacts (e.g., #work #personal)
func ExtractContacts(s string) ([]friend.Contact, error) {
	s = strings.TrimSpace(s)

	if s == "" {
		return nil, ErrNoInfo
	}

	// Extract and remove tags from input - they apply to all contacts
	tags := tag.Tags(ExtractTags(s)).ToNames()
	s = RemoveTags(s)

	parts := strings.Fields(s)
	contacts := make([]friend.Contact, 0, len(parts))

	for _, part := range parts {
		c := parseContactPart(part)
		if c.Value != "" {
			c.Tags = tags
			contacts = append(contacts, c)
		}
	}

	if len(contacts) == 0 {
		return nil, ErrNoInfo
	}

	return contacts, nil
}

// ExtractContact parses a single contact (for editing)
func ExtractContact(s string) (friend.Contact, error) {
	contacts, err := ExtractContacts(s)
	if err != nil {
		return friend.Contact{}, err
	}

	if len(contacts) == 0 {
		return friend.Contact{}, ErrNoInfo
	}

	return contacts[0], nil
}

func parseContactPart(part string) friend.Contact {
	part = strings.TrimSpace(part)

	// Check for type:value format
	if idx := strings.Index(part, ":"); idx > 0 && idx < len(part)-1 {
		prefix := strings.ToLower(part[:idx])
		value := part[idx+1:]

		if ct, ok := contactTypeAliases[prefix]; ok {
			return friend.Contact{
				Type:  ct,
				Value: value,
			}
		}
	}

	// Auto-detect type
	ct := detectContactType(part)

	return friend.Contact{
		Type:  ct,
		Value: part,
	}
}

func detectContactType(value string) friend.ContactType {
	// Email detection
	if strings.Contains(value, "@") && strings.Contains(value, ".") {
		return friend.ContactTypeEmail
	}

	// Phone detection (starts with + or contains mostly digits)
	if strings.HasPrefix(value, "+") {
		return friend.ContactTypePhone
	}

	// Check if it's mostly digits (phone number)
	digitCount := 0

	for _, r := range value {
		if r >= '0' && r <= '9' {
			digitCount++
		}
	}

	if digitCount > len(value)/2 && digitCount >= 7 {
		return friend.ContactTypePhone
	}

	return friend.ContactTypeOther
}

// RenderContact formats a single contact for display/editing
func RenderContact(c friend.Contact) string {
	var sb strings.Builder

	// Find shortest alias for the type
	alias := string(c.Type)
	for a, t := range contactTypeAliases {
		if t == c.Type && len(a) < len(alias) {
			alias = a
		}
	}

	if alias != string(friend.ContactTypeOther) && alias != "" {
		sb.WriteString(alias)
		sb.WriteString(":")
	}

	sb.WriteString(c.Value)

	if len(c.Tags) > 0 {
		sb.WriteString(" ")
		sb.WriteString(RenderTags(c.Tags))
	}

	return sb.String()
}

// RenderContacts formats multiple contacts for display/editing
func RenderContacts(contacts []friend.Contact) string { //nolint:cyclop
	parts := make([]string, 0, len(contacts))

	// Collect unique tags across all contacts
	allTags := make(map[string]bool)

	for _, c := range contacts {
		for _, t := range c.Tags {
			allTags[t] = true
		}
	}

	for _, c := range contacts {
		// Render contact without individual tags
		var sb strings.Builder

		alias := string(c.Type)
		for a, t := range contactTypeAliases {
			if t == c.Type && len(a) < len(alias) {
				alias = a
			}
		}

		if alias != string(friend.ContactTypeOther) && alias != "" {
			sb.WriteString(alias)
			sb.WriteString(":")
		}

		sb.WriteString(c.Value)
		parts = append(parts, sb.String())
	}

	result := strings.Join(parts, " ")

	// Append shared tags at the end
	if len(allTags) > 0 {
		tags := make([]string, 0, len(allTags))

		for t := range allTags {
			tags = append(tags, t)
		}

		result += " " + RenderTags(tags)
	}

	return result
}
