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

package wishlist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/roma-glushko/frens/internal/log"
	"github.com/roma-glushko/frens/internal/version"
)

var DefaultUserAgent = fmt.Sprintf("Mozilla/5.0 (Frens/%s)", version.Version)

type ProductInfo struct {
	Name          string `json:"name"`
	PriceAmount   string `json:"price_amount"`
	PriceCurrency string `json:"price_currency"`
}

type URLManager struct {
	UserAgent  string
	httpClient *http.Client
}

func NewURLManager() *URLManager {
	return &URLManager{
		UserAgent:  DefaultUserAgent,
		httpClient: &http.Client{},
	}
}

func (m *URLManager) Scrape(ctx context.Context, url string) (*ProductInfo, error) {
	if url == "" {
		return nil, errors.New("URL cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create product info scrape request: %w", err)
	}

	req.Header.Set("User-Agent", m.UserAgent)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape product info: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"failed to scrape product info: received status code %d",
			resp.StatusCode,
		)
	}

	p, err := m.extractProductInfo(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to extract product info: %w", err)
	}

	return p, nil
}

func (m *URLManager) extractProductInfo(r io.Reader) (*ProductInfo, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var p ProductInfo

	p.Name, _ = doc.Find("meta[property='og:title']").Attr("content")

	if p.Name == "" {
		p.Name = doc.Find("title").Text()
	}

	doc.Find("script[type='application/ld+json']").Each(func(_ int, s *goquery.Selection) {
		var obj map[string]interface{}
		err := json.Unmarshal([]byte(s.Text()), &obj)
		if err != nil {
			log.Infof("Failed to parse JSON-LD: %v", err)
			return
		}

		if obj["@type"] != "Product" {
			return
		}

		log.Debug("Found JSON-LD Product type")
		log.Debugf("%+v", obj)

		if offers, ok := obj["offers"].([]interface{}); ok && len(offers) > 0 {
			for _, offer := range offers {
				offerMap, ok := offer.(map[string]interface{})
				if !ok {
					log.Debugf("Skipping non-map offer: %v", offer)
					continue
				}

				if price, ok := offerMap["highPrice"].(string); ok {
					priceFloat, err := strconv.ParseFloat(price, 64)
					if err == nil {
						p.PriceAmount = fmt.Sprintf("%.2f", priceFloat)
						log.Debugf(
							"Extracted price amount from offers[].highPrice: %s",
							p.PriceAmount,
						)
					} else {
						log.Debugf("Failed to parse price amount from offers[].highPrice: %v", err)
					}
				}

				if price, ok := offerMap["price"].(string); ok {
					priceFloat, err := strconv.ParseFloat(price, 64)
					if err == nil {
						p.PriceAmount = fmt.Sprintf("%.2f", priceFloat)
						log.Debugf("Extracted price amount from offers[].price: %s", p.PriceAmount)
					} else {
						log.Debugf("Failed to parse price amount from offers[].price: %v", err)
					}
				}

				if price, ok := offerMap["price"].(float64); ok {
					p.PriceAmount = fmt.Sprintf("%.2f", price)
					log.Debugf("Extracted price amount from offers[].price: %s", p.PriceAmount)
				}

				if curr, ok := offerMap["priceCurrency"].(string); ok {
					p.PriceCurrency = curr
					log.Debugf(
						"Extracted price currency from offers[].priceCurrency: %s",
						p.PriceCurrency,
					)
				}
			}
		}

		if offers, ok := obj["offers"].(map[string]interface{}); ok {
			if price, ok := offers["highPrice"].(string); ok {
				priceFloat, err := strconv.ParseFloat(price, 64)

				if err == nil {
					p.PriceAmount = fmt.Sprintf("%.2f", priceFloat)
					log.Debugf("Extracted price amount from offers.highPrice: %s", p.PriceAmount)
				} else {
					log.Debugf("Failed to parse price amount from offers.highPrice: %v", err)
				}
			}

			if price, ok := offers["price"].(string); ok {
				priceFloat, err := strconv.ParseFloat(price, 64)

				if err == nil {
					p.PriceAmount = fmt.Sprintf("%.2f", priceFloat)
					log.Debugf("Extracted price amount from offers.price: %s", p.PriceAmount)
				} else {
					log.Debugf("Failed to parse price amount from offers.price: %v", err)
				}
			}

			if price, ok := offers["price"].(float64); ok {
				p.PriceAmount = fmt.Sprintf("%.2f", price)
				log.Debugf("Extracted price amount from offers.price: %s", p.PriceAmount)
			}

			if curr, ok := offers["priceCurrency"].(string); ok {
				p.PriceCurrency = curr
				log.Debugf("Extracted price amount from offers.priceCurrency: %s", p.PriceCurrency)
			}
		}
	})

	if p.PriceAmount == "" {
		p.PriceAmount, _ = doc.Find("meta[property='product:price:amount']").Attr("content")
		p.PriceCurrency, _ = doc.Find("meta[property='product:price:currency']").Attr("content")
	}

	if p.PriceAmount == "" {
		p.PriceAmount, _ = doc.Find("meta[itemprop='price']").Attr("content")
		p.PriceCurrency, _ = doc.Find("meta[itemprop='priceCurrency']").Attr("content")
	}

	if p.PriceAmount == "" {
		// 3. Heuristic search for currency strings
		text := doc.Text()
		re := regexp.MustCompile(`\d[\d\s,]*\.?\d*\s*(грн|uah|usd|€|₴)`)
		if m := re.FindString(text); m != "" {
			cleaned := regexp.MustCompile(`[^\d.,]`).ReplaceAllString(m, "")
			cleaned = strings.ReplaceAll(cleaned, " ", "")
			cleaned = strings.ReplaceAll(cleaned, " ", "")

			price, err := strconv.ParseFloat(cleaned, 64)

			if err == nil {
				p.PriceAmount = fmt.Sprintf("%.2f", price)
			} else {
				log.Infof("Failed to parse price amount: %v", err)
			}

			if strings.Contains(strings.ToLower(m), "грн") ||
				strings.Contains(strings.ToLower(m), "₴") {
				p.PriceCurrency = "UAH"
			} else if strings.Contains(strings.ToLower(m), "usd") {
				p.PriceCurrency = "USD"
			}
		}
	}

	p.Name = strings.TrimSpace(p.Name)
	p.PriceAmount = strings.TrimSpace(p.PriceAmount)
	p.PriceCurrency = strings.TrimSpace(p.PriceCurrency)

	return &p, nil
}
