package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SteakItem struct {
	Name      string
	Price     string
	Available bool
}

type SteakParser struct {
	config *Config
	client *http.Client
}

func NewSteakParser(config *Config) *SteakParser {
	return &SteakParser{
		config: config,
		client: &http.Client{},
	}
}

func (p *SteakParser) ParseAndNotify() error {
	steaks, err := p.fetchSteaks()
	if err != nil {
		return fmt.Errorf("failed to fetch steaks: %w", err)
	}

	trackedSteaks := p.filterTrackedSteaks(steaks)
	if len(trackedSteaks) > 0 {
		return p.sendNotification(trackedSteaks)
	}

	log.Println("No tracked steaks found")
	return nil
}

func (p *SteakParser) fetchSteaks() ([]SteakItem, error) {
	log.Printf("Starting to fetch steaks from URL: %s", p.config.Tracking.URL)

	resp, err := p.client.Get(p.config.Tracking.URL)
	if err != nil {
		log.Printf("Failed to fetch URL: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("HTTP response status: %d %s", resp.StatusCode, resp.Status)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML document: %v", err)
		return nil, err
	}

	var steaks []SteakItem
	log.Println("Starting to parse product items...")

	doc.Find(".product-item").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".product-title").Text())

		// Get the raw price text and clean it up
		rawPrice := strings.TrimSpace(s.Find(".product-new, .product-price").Text())
		price := p.cleanPrice(rawPrice)

		// Check if product is available by looking for "купить" button vs "уведомить" button
		buyButton := s.Find("button, .btn").FilterFunction(func(i int, sel *goquery.Selection) bool {
			text := strings.TrimSpace(strings.ToLower(sel.Text()))
			return strings.Contains(text, "купить")
		})

		notifyButton := s.Find("button, .btn").FilterFunction(func(i int, sel *goquery.Selection) bool {
			text := strings.TrimSpace(strings.ToLower(sel.Text()))
			return strings.Contains(text, "уведомить")
		})

		available := buyButton.Length() > 0
		unavailable := notifyButton.Length() > 0

		// If we found notify button, it's definitely unavailable
		if unavailable {
			available = false
		}

		availabilityStatus := "AVAILABLE"
		if !available {
			availabilityStatus = "OUT OF STOCK"
		}

		log.Printf("Found product #%d: Name='%s', Price='%s', Status='%s'", i+1, name, price, availabilityStatus)

		if name != "" {
			steaks = append(steaks, SteakItem{
				Name:      name,
				Price:     price,
				Available: available,
			})
		} else {
			log.Printf("Skipping product #%d: empty name", i+1)
		}
	})

	log.Printf("Successfully parsed %d products", len(steaks))

	log.Println("=== ALL PRODUCTS ===")
	availableCount := 0
	for i, steak := range steaks {
		status := "OUT OF STOCK"
		if steak.Available {
			status = "AVAILABLE"
			availableCount++
		}
		log.Printf("[%d] %s - %s [%s]", i+1, steak.Name, steak.Price, status)
	}
	log.Printf("=== END OF PRODUCT LIST === (Available: %d/%d)", availableCount, len(steaks))

	return steaks, nil
}

func (p *SteakParser) cleanPrice(rawPrice string) string {
	// Remove excessive whitespace and newlines
	cleaned := regexp.MustCompile(`\s+`).ReplaceAllString(rawPrice, " ")
	cleaned = strings.TrimSpace(cleaned)
	
	// Extract the final price (usually the last occurrence of "руб.")
	priceRegex := regexp.MustCompile(`(\d+)\s*руб\.(?:\s*за\s+.*)?$`)
	matches := priceRegex.FindStringSubmatch(cleaned)
	
	if len(matches) > 1 {
		return matches[1] + " руб."
	}
	
	// Fallback: try to find any number followed by "руб."
	fallbackRegex := regexp.MustCompile(`(\d+)\s*руб\.`)
	allMatches := fallbackRegex.FindAllStringSubmatch(cleaned, -1)
	
	if len(allMatches) > 0 {
		// Take the last match (usually the actual price, not the crossed-out one)
		lastMatch := allMatches[len(allMatches)-1]
		return lastMatch[1] + " руб."
	}
	
	// If no price pattern found, return cleaned text
	return cleaned
}

func (p *SteakParser) filterTrackedSteaks(steaks []SteakItem) []SteakItem {
	var tracked []SteakItem

	for _, steak := range steaks {
		for _, trackName := range p.config.Tracking.SteaksToTrack {
			if strings.Contains(strings.ToLower(steak.Name), strings.ToLower(trackName)) {
				tracked = append(tracked, steak)
				break
			}
		}
	}

	return tracked
}

func (p *SteakParser) sendNotification(steaks []SteakItem) error {
	message := "🥩 Найденные стейки:\n\n"
	for _, steak := range steaks {
		status := "❌ НЕТ В НАЛИЧИИ"
		if steak.Available {
			status = "✅ ДОСТУПЕН"
		}
		message += fmt.Sprintf("• %s - %s [%s]\n", steak.Name, steak.Price, status)
	}

	telegram := NewTelegramBot(p.config.Telegram.Token, p.config.Telegram.ChatID)
	return telegram.SendMessage(message)
}
