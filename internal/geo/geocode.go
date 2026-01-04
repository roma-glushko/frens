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

package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Coordinates represents a geographic location
type Coordinates struct {
	Lat float64
	Lng float64
}

// NominatimResult represents a single result from Nominatim API
type NominatimResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// Geocoder provides geocoding functionality using OpenStreetMap's Nominatim API
type Geocoder struct {
	client  *http.Client
	baseURL string
}

// NewGeocoder creates a new Geocoder instance
func NewGeocoder() *Geocoder {
	return &Geocoder{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://nominatim.openstreetmap.org/search",
	}
}

// Geocode resolves a location query (e.g., "Berlin, Germany") to coordinates
func (g *Geocoder) Geocode(ctx context.Context, query string) (*Coordinates, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", "1")

	reqURL := fmt.Sprintf("%s?%s", g.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Nominatim requires a User-Agent header
	req.Header.Set("User-Agent", "frens-cli/1.0 (https://github.com/roma-glushko/frens)")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var results []NominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(results) == 0 {
		return nil, nil // No results found, not an error
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return &Coordinates{
		Lat: lat,
		Lng: lng,
	}, nil
}

// GeocodeLocation builds a query from name and country and resolves coordinates
func (g *Geocoder) GeocodeLocation(
	ctx context.Context,
	name, country string,
) (*Coordinates, error) {
	query := name

	if country != "" {
		query = fmt.Sprintf("%s, %s", name, country)
	}

	return g.Geocode(ctx, query)
}
