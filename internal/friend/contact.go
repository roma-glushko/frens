// Copyright 2025 Roma Hlushko
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

package friend

import (
	"errors"
	"fmt"
	"strings"

	"github.com/roma-glushko/frens/internal/tag"
)

type ContactType string

const (
	ContactTypeEmail     ContactType = "email"
	ContactTypePhone     ContactType = "phone"
	ContactTypeTelegram  ContactType = "telegram"
	ContactTypeWhatsApp  ContactType = "whatsapp"
	ContactTypeSignal    ContactType = "signal"
	ContactTypeTwitter   ContactType = "twitter"
	ContactTypeLinkedIn  ContactType = "linkedin"
	ContactTypeGitHub    ContactType = "github"
	ContactTypeInstagram ContactType = "instagram"
	ContactTypeFacebook  ContactType = "facebook"
	ContactTypeDiscord   ContactType = "discord"
	ContactTypeSlack     ContactType = "slack"
	ContactTypeOther     ContactType = "other"
)

var ContactTypes = []ContactType{
	ContactTypeEmail,
	ContactTypePhone,
	ContactTypeTelegram,
	ContactTypeWhatsApp,
	ContactTypeSignal,
	ContactTypeTwitter,
	ContactTypeLinkedIn,
	ContactTypeGitHub,
	ContactTypeInstagram,
	ContactTypeFacebook,
	ContactTypeDiscord,
	ContactTypeSlack,
	ContactTypeOther,
}

func ValidateContactType(s string) error {
	validTypes := make([]string, 0, len(ContactTypes))

	for _, ct := range ContactTypes {
		t := string(ct)
		validTypes = append(validTypes, t)

		if strings.EqualFold(s, t) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid contact type '%s' (supported: %s)",
		s,
		strings.Join(validTypes, ", "),
	)
}

func ParseContactType(s string) ContactType {
	for _, ct := range ContactTypes {
		if strings.EqualFold(s, string(ct)) {
			return ct
		}
	}

	return ContactTypeOther
}

type Contact struct {
	ID     string      `toml:"id"`
	Type   ContactType `toml:"type"`
	Value  string      `toml:"value"`
	Tags   []string    `toml:"tags,omitempty"`
	Person string      `toml:"-"`
}

func (c *Contact) SetTags(tags []string) {
	c.Tags = tags
}

func (c *Contact) GetTags() []string {
	return c.Tags
}

var _ tag.Tagged = (*Contact)(nil)

func (c *Contact) Validate() error {
	if c.Value == "" {
		return errors.New("contact value cannot be empty")
	}

	if c.Type == "" {
		return errors.New("contact type cannot be empty")
	}

	return nil
}

func (c *Contact) String() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Value)
}
