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
	"encoding/json"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/log"
)

func init() {
	log.RegisterFormatter(log.FormatJSON, friend.Person{}, PersonJSONFormatter{})
	log.RegisterFormatter(log.FormatJSON, friend.Contact{}, ContactJSONFormatter{})
	log.RegisterFormatter(log.FormatJSON, friend.Event{}, EventJSONFormatter{})
	log.RegisterFormatter(log.FormatJSON, friend.Location{}, LocationJSONFormatter{})
	log.RegisterFormatter(log.FormatJSON, friend.Date{}, DateJSONFormatter{})
	log.RegisterFormatter(log.FormatJSON, friend.WishlistItem{}, WishlistItemJSONFormatter{})
}

// ============================================================================
// Person JSON Formatter
// ============================================================================

type PersonJSONFormatter struct{}

var _ log.Formatter = (*PersonJSONFormatter)(nil)

func (p PersonJSONFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	person, ok := e.(friend.Person)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func (p PersonJSONFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	persons, ok := el.([]friend.Person)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(persons, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

// ============================================================================
// Contact JSON Formatter
// ============================================================================

type ContactJSONFormatter struct{}

var _ log.Formatter = (*ContactJSONFormatter)(nil)

func (f ContactJSONFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	var c friend.Contact

	switch v := e.(type) {
	case friend.Contact:
		c = v
	case *friend.Contact:
		c = *v
	default:
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func (f ContactJSONFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	contacts, ok := el.([]friend.Contact)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

// ============================================================================
// Event JSON Formatter
// ============================================================================

type EventJSONFormatter struct{}

var _ log.Formatter = (*EventJSONFormatter)(nil)

func (f EventJSONFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	event, ok := e.(friend.Event)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func (f EventJSONFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	events, ok := el.([]friend.Event)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

// ============================================================================
// Location JSON Formatter
// ============================================================================

type LocationJSONFormatter struct{}

var _ log.Formatter = (*LocationJSONFormatter)(nil)

func (l LocationJSONFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	location, ok := e.(friend.Location)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func (l LocationJSONFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	locations, ok := el.([]friend.Location)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(locations, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

// ============================================================================
// Date JSON Formatter
// ============================================================================

type DateJSONFormatter struct{}

var _ log.Formatter = (*DateJSONFormatter)(nil)

func (f DateJSONFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	var dt friend.Date

	switch v := e.(type) {
	case friend.Date:
		dt = v
	case *friend.Date:
		dt = *v
	default:
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(dt, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func (f DateJSONFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	dates, ok := el.([]*friend.Date)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(dates, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

// ============================================================================
// WishlistItem JSON Formatter
// ============================================================================

type WishlistItemJSONFormatter struct{}

var _ log.Formatter = (*WishlistItemJSONFormatter)(nil)

func (f WishlistItemJSONFormatter) FormatSingle(_ log.FormatterContext, e any) (string, error) {
	var w friend.WishlistItem

	switch v := e.(type) {
	case friend.WishlistItem:
		w = v
	case *friend.WishlistItem:
		w = *v
	default:
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func (f WishlistItemJSONFormatter) FormatList(_ log.FormatterContext, el any) (string, error) {
	items, ok := el.([]friend.WishlistItem)
	if !ok {
		return "", ErrInvalidEntity
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}
