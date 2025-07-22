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

package journal

import (
	"testing"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
	"github.com/stretchr/testify/require"
)

func TestJournal_AddEvent(t *testing.T) {
	frID := "alice"
	locID := "wonderland"

	jr := Journal{
		Friends: []*friend.Person{
			{
				ID:   frID,
				Name: "Alice",
			},
		},
		Locations: []*friend.Location{
			{
				ID:   locID,
				Name: "Wonderland",
			},
		},
	}

	jr.Init()

	event, err := jr.AddEvent(friend.Event{
		Type: friend.EventTypeActivity,
		Date: time.Now().UTC(),
		Desc: "Alice went to Wonderland",
	})

	require.NoError(t, err)
	require.NotEmpty(t, event.ID)
	require.Contains(t, event.Friends, frID)
	// require.Contains(t, event.Locations, locID)

	f, err := jr.GetFriend(frID)
	require.NoError(t, err)

	require.Equal(t, 1, f.Activities)
}
