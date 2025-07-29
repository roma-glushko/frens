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

package graph

import (
	"fmt"
	"time"

	"github.com/roma-glushko/frens/internal/friend"
)

// Constants
const (
	DateFormat = "Jan 2006"
	ScaledSize = 20
)

type RGB [3]int

var colors = []RGB{
	{153, 255, 0},
	{153, 255, 0},
	{153, 204, 0},
	{153, 204, 0},
	{204, 204, 0},
	{204, 204, 0},
	{255, 153, 0},
	{255, 153, 0},
	{255, 102, 0},
	{255, 102, 0},
	{255, 51, 51},
	{255, 51, 51},
	{255, 0, 102},
	{255, 0, 102},
	{255, 0, 153},
	{255, 0, 153},
	{204, 0, 204},
	{204, 0, 204},
	{153, 0, 255},
	{153, 0, 255},
	{102, 0, 255},
	{102, 0, 255},
	{51, 51, 255},
	{51, 51, 255},
	{0, 102, 255},
	{0, 102, 255},
	{0, 153, 255},
	{0, 153, 255},
	{0, 204, 204},
	{0, 204, 204},
	{0, 255, 153},
	{0, 255, 153},
	{0, 255, 102},
	{0, 255, 102},
	{51, 255, 51},
	{51, 255, 51},
	{102, 255, 0},
	{102, 255, 0},
}

// ActivityGraph holds the data needed to produce a graph.
type ActivityGraph struct {
	filteredActivities []*friend.Event
	allActivities      []*friend.Event
	scaled             bool
	startDate          time.Time
	endDate            time.Time
}

// NewActivityGraph creates a new graph instance.
func NewActivityGraph(filtered, all []*friend.Event, scaled bool) *ActivityGraph {
	g := &ActivityGraph{
		filteredActivities: filtered,
		allActivities:      all,
		scaled:             scaled,
	}

	if len(all) > 0 {
		g.startDate = all[len(all)-1].Date
		g.endDate = all[0].Date
	}

	return g
}

// Output renders the graph lines.
func (g *ActivityGraph) Output() []string { //nolint:cyclop
	graphMap := g.toMap()
	if len(graphMap) == 0 {
		return []string{}
	}

	var globalTotal int

	if g.scaled {
		for _, counts := range graphMap {
			if counts[1] > globalTotal {
				globalTotal = counts[1]
			}
		}
	}

	var result []string

	for i := len(graphMap) - 1; i >= 0; i-- {
		month := g.months()[i]
		counts := graphMap[month]
		filteredCount, totalCount := counts[0], counts[1]

		if g.scaled && globalTotal > 0 {
			filteredCount = int(
				float64(filteredCount)*float64(ScaledSize)/float64(globalTotal) + 0.5,
			)
			totalCount = ScaledSize
		}

		bar := month + " |"
		for j := 0; j < filteredCount; j++ {
			bar += Paint("█", g.color(j))
		}

		for j := filteredCount; j < totalCount; j++ {
			bar += Paint("∙", g.color(j))
		}

		if totalCount > 0 {
			bar += Paint("|", g.color(totalCount))
		}

		result = append(result, bar)
	}

	return result
}

// toMap builds a map of "Month Year" => [filtered, total] counts.
func (g *ActivityGraph) toMap() map[string][2]int {
	graph := g.emptyGraph()
	for _, a := range g.filteredActivities {
		key := a.Date.Format(DateFormat)
		val := graph[key]
		val[0]++
		graph[key] = val
	}

	for _, a := range g.allActivities {
		key := a.Date.Format(DateFormat)
		val := graph[key]
		val[1]++
		graph[key] = val
	}

	return graph
}

// emptyGraph initializes the map with 0s for each month.
func (g *ActivityGraph) emptyGraph() map[string][2]int {
	result := make(map[string][2]int)
	for _, date := range g.monthRange() {
		result[date.Format(DateFormat)] = [2]int{0, 0}
	}

	return result
}

// months returns the list of month labels in order.
func (g *ActivityGraph) months() []string {
	monthRange := g.monthRange()
	labels := make([]string, 0, len(monthRange))

	for _, date := range monthRange {
		labels = append(labels, date.Format(DateFormat))
	}

	return labels
}

// monthRange returns a slice of time.Time representing each month from start to end.
func (g *ActivityGraph) monthRange() []time.Time {
	var months []time.Time
	if g.startDate.IsZero() || g.endDate.IsZero() {
		return months
	}

	start := time.Date(g.startDate.Year(), g.startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(g.endDate.Year(), g.endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for d := start; !d.After(end); d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}

	return months
}

// color returns the RGB color at a given position.
func (g *ActivityGraph) color(x int) RGB {
	return colors[x%len(colors)]
}

// Paint applies a color to a character. (You can customize this to output ANSI codes, etc.)
func Paint(char string, color RGB) string {
	// For real ANSI coloring, replace with proper formatting:
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", color[0], color[1], color[2], char)
}
