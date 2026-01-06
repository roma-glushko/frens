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

package config

import (
	"testing"
)

func TestMatchRuleWithContext_Tags(t *testing.T) {
	config := &NotificationConfig{
		Rules: []RoutingRule{
			{
				ID:         "family-rule",
				Priority:   1,
				MatchTags:  []string{"family"},
				ChannelIDs: []string{"discord-family"},
			},
			{
				ID:         "work-rule",
				Priority:   2,
				MatchTags:  []string{"work"},
				ChannelIDs: []string{"tg-main"},
			},
		},
	}

	tests := []struct {
		name     string
		ctx      MatchRuleContext
		expected string
	}{
		{
			name:     "match family tag",
			ctx:      MatchRuleContext{Tags: []string{"family"}},
			expected: "family-rule",
		},
		{
			name:     "match work tag",
			ctx:      MatchRuleContext{Tags: []string{"work"}},
			expected: "work-rule",
		},
		{
			name:     "match first rule by priority",
			ctx:      MatchRuleContext{Tags: []string{"family", "work"}},
			expected: "family-rule",
		},
		{
			name:     "no match",
			ctx:      MatchRuleContext{Tags: []string{"other"}},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.MatchRuleWithContext(tt.ctx)
			if tt.expected == "" {
				if result != nil {
					t.Errorf("expected no match, got %s", result.ID)
				}
			} else {
				if result == nil {
					t.Errorf("expected %s, got nil", tt.expected)
				} else if result.ID != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.ID)
				}
			}
		})
	}
}

func TestMatchRuleWithContext_Keywords(t *testing.T) {
	config := &NotificationConfig{
		Rules: []RoutingRule{
			{
				ID:            "birthday-rule",
				Priority:      1,
				MatchKeywords: []string{"birthday", "bday"},
				ChannelIDs:    []string{"discord-family"},
			},
			{
				ID:            "gift-rule",
				Priority:      2,
				MatchKeywords: []string{"gift", "present"},
				ChannelIDs:    []string{"tg-main"},
			},
		},
	}

	tests := []struct {
		name     string
		ctx      MatchRuleContext
		expected string
	}{
		{
			name:     "match birthday keyword",
			ctx:      MatchRuleContext{Content: "John's birthday is coming up"},
			expected: "birthday-rule",
		},
		{
			name:     "match bday keyword",
			ctx:      MatchRuleContext{Content: "Don't forget bday party"},
			expected: "birthday-rule",
		},
		{
			name:     "match gift keyword",
			ctx:      MatchRuleContext{Content: "Buy a gift for Jane"},
			expected: "gift-rule",
		},
		{
			name:     "case insensitive match",
			ctx:      MatchRuleContext{Content: "BIRTHDAY reminder"},
			expected: "birthday-rule",
		},
		{
			name:     "no match",
			ctx:      MatchRuleContext{Content: "Meeting reminder"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.MatchRuleWithContext(tt.ctx)
			if tt.expected == "" {
				if result != nil {
					t.Errorf("expected no match, got %s", result.ID)
				}
			} else {
				if result == nil {
					t.Errorf("expected %s, got nil", tt.expected)
				} else if result.ID != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result.ID)
				}
			}
		})
	}
}

func TestMatchRuleWithContext_TagsOrKeywords(t *testing.T) {
	config := &NotificationConfig{
		Rules: []RoutingRule{
			{
				ID:            "combined-rule",
				Priority:      1,
				MatchTags:     []string{"family"},
				MatchKeywords: []string{"birthday"},
				MatchAll:      false, // OR logic
				ChannelIDs:    []string{"discord-family"},
			},
		},
	}

	tests := []struct {
		name     string
		ctx      MatchRuleContext
		expected bool
	}{
		{
			name:     "match by tag only",
			ctx:      MatchRuleContext{Tags: []string{"family"}, Content: "some content"},
			expected: true,
		},
		{
			name:     "match by keyword only",
			ctx:      MatchRuleContext{Tags: []string{"work"}, Content: "birthday party"},
			expected: true,
		},
		{
			name:     "match by both",
			ctx:      MatchRuleContext{Tags: []string{"family"}, Content: "birthday party"},
			expected: true,
		},
		{
			name:     "no match",
			ctx:      MatchRuleContext{Tags: []string{"work"}, Content: "meeting"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.MatchRuleWithContext(tt.ctx)
			if tt.expected && result == nil {
				t.Errorf("expected match, got nil")
			}
			if !tt.expected && result != nil {
				t.Errorf("expected no match, got %s", result.ID)
			}
		})
	}
}

func TestMatchRuleWithContext_MatchAll(t *testing.T) {
	config := &NotificationConfig{
		Rules: []RoutingRule{
			{
				ID:            "strict-rule",
				Priority:      1,
				MatchTags:     []string{"family", "urgent"},
				MatchKeywords: []string{"birthday"},
				MatchAll:      true, // AND logic
				ChannelIDs:    []string{"discord-family"},
			},
		},
	}

	tests := []struct {
		name     string
		ctx      MatchRuleContext
		expected bool
	}{
		{
			name:     "match all conditions",
			ctx:      MatchRuleContext{Tags: []string{"family", "urgent"}, Content: "birthday party"},
			expected: true,
		},
		{
			name:     "missing one tag",
			ctx:      MatchRuleContext{Tags: []string{"family"}, Content: "birthday party"},
			expected: false,
		},
		{
			name:     "missing keyword",
			ctx:      MatchRuleContext{Tags: []string{"family", "urgent"}, Content: "party"},
			expected: false,
		},
		{
			name:     "missing all",
			ctx:      MatchRuleContext{Tags: []string{"work"}, Content: "meeting"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.MatchRuleWithContext(tt.ctx)
			if tt.expected && result == nil {
				t.Errorf("expected match, got nil")
			}
			if !tt.expected && result != nil {
				t.Errorf("expected no match, got %s", result.ID)
			}
		})
	}
}

func TestMatchRuleWithContext_EmptyRuleMatchesAll(t *testing.T) {
	config := &NotificationConfig{
		Rules: []RoutingRule{
			{
				ID:         "catch-all",
				Priority:   1,
				ChannelIDs: []string{"default"},
			},
		},
	}

	result := config.MatchRuleWithContext(MatchRuleContext{
		Tags:    []string{"any"},
		Content: "anything",
	})

	if result == nil {
		t.Error("expected catch-all rule to match")
	}
	if result.ID != "catch-all" {
		t.Errorf("expected catch-all, got %s", result.ID)
	}
}
