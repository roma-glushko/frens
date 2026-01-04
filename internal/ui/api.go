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

package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/roma-glushko/frens/internal/journal"
	"github.com/roma-glushko/frens/internal/store"
	"github.com/roma-glushko/frens/internal/sync"
)

// ComprehensiveStats contains all statistics for the Stats page.
type ComprehensiveStats struct {
	Counts           journal.Stats       `json:"counts"`
	TopFriends       []RankedItem        `json:"topFriends"`
	TopLocations     []RankedItem        `json:"topLocations"`
	TopTags          []RankedItem        `json:"topTags"`
	ActivityTimeline []TimelineDataPoint `json:"activityTimeline"`
	Insights         []Insight           `json:"insights"`
}

// RankedItem represents a ranked item (friend, location, or tag).
type RankedItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Count        int    `json:"count"`
	LastActivity string `json:"lastActivity,omitempty"`
}

// TimelineDataPoint represents activity data for a specific month.
type TimelineDataPoint struct {
	Month      string `json:"month"`
	Activities int    `json:"activities"`
	Notes      int    `json:"notes"`
}

// Insight represents an actionable insight or recommendation.
type Insight struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FriendID    string `json:"friendId,omitempty"`
}

// SyncStatus represents the current sync status of the journal.
type SyncStatus struct {
	GitInstalled bool   `json:"gitInstalled"`
	GitInited    bool   `json:"gitInited"`
	Branch       string `json:"branch,omitempty"`
	HasChanges   bool   `json:"hasChanges"`
	ChangeCount  int    `json:"changeCount"`
}

// FeedItem represents a unified feed item (activity, note, friend added, location added).
type FeedItem struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"` // "activity", "note", "friend_added", "location_added"
	Date        string   `json:"date"`
	Description string   `json:"description"`
	FriendIDs   []string `json:"friendIds,omitempty"`
	LocationIDs []string `json:"locationIds,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	// For friend_added/location_added
	EntityID   string `json:"entityId,omitempty"`
	EntityName string `json:"entityName,omitempty"`
}

// API holds the dependencies for API handlers.
type API struct {
	store store.Store
}

// NewAPI creates a new API instance.
func NewAPI(s store.Store) *API {
	return &API{store: s}
}

// RegisterRoutes registers all API routes on the given mux.
func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/friends", a.handleListFriends)
	mux.HandleFunc("GET /api/friends/{id}", a.handleGetFriend)
	mux.HandleFunc("GET /api/friends/{id}/activities", a.handleGetFriendActivities)
	mux.HandleFunc("GET /api/friends/{id}/notes", a.handleGetFriendNotes)
	mux.HandleFunc("GET /api/locations", a.handleListLocations)
	mux.HandleFunc("GET /api/locations/{id}", a.handleGetLocation)
	mux.HandleFunc("GET /api/locations/{id}/activities", a.handleGetLocationActivities)
	mux.HandleFunc("GET /api/locations/{id}/notes", a.handleGetLocationNotes)
	mux.HandleFunc("GET /api/notes", a.handleListNotes)
	mux.HandleFunc("GET /api/activities", a.handleListActivities)
	mux.HandleFunc("GET /api/stats", a.handleGetStats)
	mux.HandleFunc("GET /api/stats/comprehensive", a.handleGetComprehensiveStats)
	mux.HandleFunc("GET /api/sync/status", a.handleGetSyncStatus)
	mux.HandleFunc("GET /api/feed", a.handleGetFeed)
}

// handleListFriends returns all friends.
func (a *API) handleListFriends(w http.ResponseWriter, r *http.Request) {
	var friends []friend.Person

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		friends = j.ListFriends(friend.ListFriendQuery{
			SortBy:    friend.SortAlpha,
			SortOrder: friend.SortOrderDirect,
		})

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(friends); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetFriend returns a single friend by ID.
func (a *API) handleGetFriend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "friend id required", http.StatusBadRequest)
		return
	}

	var foundFriend friend.Person

	var found bool

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		f, err := j.GetFriend(id)
		if err != nil {
			return err
		}

		foundFriend = f

		found = true

		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !found {
		http.Error(w, "friend not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(foundFriend); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetFriendActivities returns activities for a specific friend.
func (a *API) handleGetFriendActivities(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "friend id required", http.StatusBadRequest)
		return
	}

	var activities []friend.Event

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		for _, event := range j.Activities {
			for _, friendID := range event.FriendIDs {
				if friendID == id {
					activities = append(activities, *event)
					break
				}
			}
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if activities == nil {
		activities = []friend.Event{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(activities); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetFriendNotes returns notes for a specific friend.
func (a *API) handleGetFriendNotes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "friend id required", http.StatusBadRequest)
		return
	}

	var notes []friend.Event

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		for _, event := range j.Notes {
			for _, friendID := range event.FriendIDs {
				if friendID == id {
					notes = append(notes, *event)
					break
				}
			}
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []friend.Event{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleListLocations returns all locations.
func (a *API) handleListLocations(w http.ResponseWriter, r *http.Request) {
	var locations []friend.Location

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		locations = j.ListLocations(friend.ListLocationQuery{
			SortBy:    friend.SortAlpha,
			SortOrder: friend.SortOrderDirect,
		})

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if locations == nil {
		locations = []friend.Location{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetLocation returns a single location by ID.
func (a *API) handleGetLocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "location id required", http.StatusBadRequest)
		return
	}

	var foundLocation friend.Location

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		loc, err := j.GetLocation(id)
		if err != nil {
			return err
		}

		foundLocation = loc

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(foundLocation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetLocationActivities returns activities for a specific location.
func (a *API) handleGetLocationActivities(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "location id required", http.StatusBadRequest)
		return
	}

	var activities []friend.Event

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		for _, event := range j.Activities {
			for _, locationID := range event.LocationIDs {
				if locationID == id {
					activities = append(activities, *event)
					break
				}
			}
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if activities == nil {
		activities = []friend.Event{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(activities); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetLocationNotes returns notes for a specific location.
func (a *API) handleGetLocationNotes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "location id required", http.StatusBadRequest)
		return
	}

	var notes []friend.Event

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		for _, event := range j.Notes {
			for _, locationID := range event.LocationIDs {
				if locationID == id {
					notes = append(notes, *event)
					break
				}
			}
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []friend.Event{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleListNotes returns all notes.
func (a *API) handleListNotes(w http.ResponseWriter, r *http.Request) {
	var notes []friend.Event

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		for _, note := range j.Notes {
			notes = append(notes, *note)
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []friend.Event{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleListActivities returns all activities.
func (a *API) handleListActivities(w http.ResponseWriter, r *http.Request) {
	var activities []friend.Event

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		for _, activity := range j.Activities {
			activities = append(activities, *activity)
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if activities == nil {
		activities = []friend.Event{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(activities); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetStats returns journal statistics.
func (a *API) handleGetStats(w http.ResponseWriter, r *http.Request) {
	var stats journal.Stats

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		stats = j.Stats()
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetComprehensiveStats returns comprehensive statistics for the Stats page.
func (a *API) handleGetComprehensiveStats(w http.ResponseWriter, r *http.Request) { //nolint:cyclop
	var result ComprehensiveStats

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		// Basic counts
		result.Counts = j.Stats()

		// Top friends by activity count
		friends := j.ListFriends(friend.ListFriendQuery{
			SortBy:    friend.SortActivities,
			SortOrder: friend.SortOrderReverse,
		})

		for i, f := range friends {
			if i >= 5 {
				break
			}

			item := RankedItem{
				ID:    f.ID,
				Name:  f.Name,
				Count: f.Activities,
			}

			if !f.MostRecentActivity.IsZero() {
				item.LastActivity = f.MostRecentActivity.Format(time.RFC3339)
			}

			result.TopFriends = append(result.TopFriends, item)
		}

		// Top locations by activity count
		locations := j.ListLocations(friend.ListLocationQuery{
			SortBy:    friend.SortActivities,
			SortOrder: friend.SortOrderReverse,
		})

		for i, l := range locations {
			if i >= 5 {
				break
			}

			item := RankedItem{
				ID:    l.ID,
				Name:  l.Name,
				Count: l.Activities,
			}

			if !l.MostRecentActivity.IsZero() {
				item.LastActivity = l.MostRecentActivity.Format(time.RFC3339)
			}

			result.TopLocations = append(result.TopLocations, item)
		}

		// Top tags by occurrence count
		tagCounts := make(map[string]int)

		for _, event := range j.Activities {
			for _, tag := range event.Tags {
				tagCounts[tag]++
			}
		}

		for _, event := range j.Notes {
			for _, tag := range event.Tags {
				tagCounts[tag]++
			}
		}

		type tagCount struct {
			tag   string
			count int
		}

		var tags []tagCount

		for tag, count := range tagCounts {
			tags = append(tags, tagCount{tag, count})
		}

		sort.Slice(tags, func(i, j int) bool {
			return tags[i].count > tags[j].count
		})

		for i, t := range tags {
			if i >= 10 {
				break
			}

			result.TopTags = append(result.TopTags, RankedItem{
				ID:    t.tag,
				Name:  t.tag,
				Count: t.count,
			})
		}

		// Activity timeline (last 12 months)
		now := time.Now()
		monthData := make(map[string]*TimelineDataPoint)

		// Initialize last 12 months
		for i := 11; i >= 0; i-- {
			m := now.AddDate(0, -i, 0)
			key := m.Format("Jan 2006")
			monthData[key] = &TimelineDataPoint{Month: key}
		}

		// Count activities per month
		for _, event := range j.Activities {
			key := event.Date.Format("Jan 2006")
			if dp, ok := monthData[key]; ok {
				dp.Activities++
			}
		}

		for _, event := range j.Notes {
			key := event.Date.Format("Jan 2006")
			if dp, ok := monthData[key]; ok {
				dp.Notes++
			}
		}

		// Convert to sorted slice
		for i := 11; i >= 0; i-- {
			m := now.AddDate(0, -i, 0)
			key := m.Format("Jan 2006")
			result.ActivityTimeline = append(result.ActivityTimeline, *monthData[key])
		}

		// Insights: Friends to reconnect with (no activity in 30+ days)
		thirtyDaysAgo := now.AddDate(0, 0, -30)

		for _, f := range friends {
			if f.Activities > 0 && !f.MostRecentActivity.IsZero() &&
				f.MostRecentActivity.Before(thirtyDaysAgo) {
				days := int(now.Sub(f.MostRecentActivity).Hours() / 24)
				result.Insights = append(result.Insights, Insight{
					Type:        "reconnect",
					Title:       f.Name,
					Description: formatDaysAgo(days),
					FriendID:    f.ID,
				})
			}
		}

		// Limit insights to 5
		if len(result.Insights) > 5 {
			result.Insights = result.Insights[:5]
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure empty arrays instead of null
	if result.TopFriends == nil {
		result.TopFriends = []RankedItem{}
	}

	if result.TopLocations == nil {
		result.TopLocations = []RankedItem{}
	}

	if result.TopTags == nil {
		result.TopTags = []RankedItem{}
	}

	if result.ActivityTimeline == nil {
		result.ActivityTimeline = []TimelineDataPoint{}
	}

	if result.Insights == nil {
		result.Insights = []Insight{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetFeed returns a unified feed of recent activity.
func (a *API) handleGetFeed(w http.ResponseWriter, r *http.Request) {
	var feed []FeedItem

	err := a.store.Tx(r.Context(), func(j *journal.Journal) error {
		// Add activities
		for _, event := range j.Activities {
			feed = append(feed, FeedItem{
				ID:          event.ID,
				Type:        "activity",
				Date:        event.Date.Format(time.RFC3339),
				Description: event.Desc,
				FriendIDs:   event.FriendIDs,
				LocationIDs: event.LocationIDs,
				Tags:        event.Tags,
			})
		}

		// Add notes
		for _, event := range j.Notes {
			feed = append(feed, FeedItem{
				ID:          event.ID,
				Type:        "note",
				Date:        event.Date.Format(time.RFC3339),
				Description: event.Desc,
				FriendIDs:   event.FriendIDs,
				LocationIDs: event.LocationIDs,
				Tags:        event.Tags,
			})
		}

		// Add friend additions
		for _, f := range j.Friends {
			if !f.CreatedAt.IsZero() {
				feed = append(feed, FeedItem{
					ID:          "friend-" + f.ID,
					Type:        "friend_added",
					Date:        f.CreatedAt.Format(time.RFC3339),
					Description: "Added " + f.Name + " as a friend",
					EntityID:    f.ID,
					EntityName:  f.Name,
				})
			}
		}

		// Add location additions
		for _, l := range j.Locations {
			if !l.CreatedAt.IsZero() {
				feed = append(feed, FeedItem{
					ID:          "location-" + l.ID,
					Type:        "location_added",
					Date:        l.CreatedAt.Format(time.RFC3339),
					Description: "Added " + l.Name + " as a location",
					EntityID:    l.ID,
					EntityName:  l.Name,
				})
			}
		}

		// Sort by date descending
		sort.Slice(feed, func(i, j int) bool {
			ti, _ := time.Parse(time.RFC3339, feed[i].Date)
			tj, _ := time.Parse(time.RFC3339, feed[j].Date)

			return ti.After(tj)
		})

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if feed == nil {
		feed = []FeedItem{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGetSyncStatus returns the current sync status of the journal.
func (a *API) handleGetSyncStatus(w http.ResponseWriter, r *http.Request) { //nolint:cyclop
	status := SyncStatus{}

	git := sync.NewGit(a.store.Path())

	// Check if git is installed
	if err := git.Installed(); err != nil {
		status.GitInstalled = false

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	status.GitInstalled = true

	// Check if git repository is initialized
	if err := git.Inited(); err != nil {
		status.GitInited = false

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	status.GitInited = true

	// Get current branch
	if branch, err := git.GetBranchName(r.Context()); err == nil {
		status.Branch = branch
	}

	// Get uncommitted changes
	if gitStatus, err := git.GetStatus(r.Context()); err == nil {
		if gitStatus != "" {
			status.HasChanges = true
			// Count number of changed files (each line in porcelain output is a file)
			lines := 0

			for _, c := range gitStatus {
				if c == '\n' {
					lines++
				}
			}

			status.ChangeCount = lines + 1 // +1 for the last line without newline
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func formatDaysAgo(days int) string {
	if days == 1 {
		return "1 day ago"
	}

	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}

	if days < 30 {
		weeks := days / 7

		if weeks == 1 {
			return "1 week ago"
		}

		return fmt.Sprintf("%d weeks ago", weeks)
	}

	months := days / 30

	if months == 1 {
		return "1 month ago"
	}

	return fmt.Sprintf("%d months ago", months)
}
