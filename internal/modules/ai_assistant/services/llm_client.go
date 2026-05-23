package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/ai_assistant/models"
)

type BookingExtraction struct {
	Intent          string `json:"intent"`
	ServiceQuery    string `json:"service_query"`
	ServiceID       string `json:"service_id"`
	DoctorID        string `json:"doctor_id"`
	ClinicAddressID string `json:"clinic_address_id"`
	Date            string `json:"date"`
	Time            string `json:"time"`
}

type LLMClient interface {
	ExtractBookingInfo(message string, state models.BookingState, history []models.ChatMessage) (BookingExtraction, error)
}

type OpenAIClient struct {
	cfg config.Config
}

func NewOpenAIClient(cfg config.Config) *OpenAIClient {
	return &OpenAIClient{cfg: cfg}
}

func (c *OpenAIClient) ExtractBookingInfo(message string, state models.BookingState, history []models.ChatMessage) (BookingExtraction, error) {
	if c.cfg.OpenAIAPIKey == "" {
		return fallbackExtract(message, state), nil
	}

	prompt := buildExtractionPrompt(message, state, history)
	body := map[string]interface{}{
		"model": c.cfg.OpenAIModel,
		"input": prompt,
		"text": map[string]interface{}{
			"format": map[string]interface{}{
				"type":   "json_schema",
				"name":   "booking_extraction",
				"strict": true,
				"schema": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": false,
					"properties": map[string]interface{}{
						"intent":            map[string]interface{}{"type": "string"},
						"service_query":     map[string]interface{}{"type": "string"},
						"service_id":        map[string]interface{}{"type": "string"},
						"doctor_id":         map[string]interface{}{"type": "string"},
						"clinic_address_id": map[string]interface{}{"type": "string"},
						"date":              map[string]interface{}{"type": "string"},
						"time":              map[string]interface{}{"type": "string"},
					},
					"required": []string{"intent", "service_query", "service_id", "doctor_id", "clinic_address_id", "date", "time"},
				},
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return BookingExtraction{}, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(payload))
	if err != nil {
		return BookingExtraction{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.cfg.OpenAIAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fallbackExtract(message, state), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fallbackExtract(message, state), nil
	}

	var openAIResp struct {
		Output []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return BookingExtraction{}, err
	}

	for _, output := range openAIResp.Output {
		for _, content := range output.Content {
			if content.Text == "" {
				continue
			}
			var extraction BookingExtraction
			if err := json.Unmarshal([]byte(content.Text), &extraction); err != nil {
				return BookingExtraction{}, err
			}
			return extraction, nil
		}
	}

	return BookingExtraction{}, errors.New("empty ai response")
}

func buildExtractionPrompt(message string, state models.BookingState, history []models.ChatMessage) string {
	historyLines := make([]string, 0, len(history))
	for i := len(history) - 1; i >= 0; i-- {
		historyLines = append(historyLines, history[i].Role+": "+history[i].Content)
	}
	stateJSON, _ := json.Marshal(state)
	return "Extract dental appointment booking information from the latest user message.\n" +
		"Return only JSON matching the schema. Use empty string for unknown fields.\n" +
		"Intent must be create_appointment unless the user is clearly not booking.\n" +
		"Resolve relative dates like tomorrow using today's date: " + time.Now().Format("2006-01-02") + ".\n" +
		"Current booking state: " + string(stateJSON) + "\n" +
		"Recent messages:\n" + strings.Join(historyLines, "\n") + "\n" +
		"Latest user message: " + message
}

func fallbackExtract(message string, state models.BookingState) BookingExtraction {
	text := strings.ToLower(strings.TrimSpace(message))
	extraction := BookingExtraction{Intent: "create_appointment"}

	uuidRe := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	ids := uuidRe.FindAllString(message, -1)
	if len(ids) > 0 {
		switch state.Step {
		case "collect_service":
			extraction.ServiceID = ids[0]
		case "collect_clinic":
			extraction.ClinicAddressID = ids[0]
		case "collect_doctor":
			extraction.DoctorID = ids[0]
		default:
			extraction.ServiceID = ids[0]
		}
	}

	dateRe := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	if date := dateRe.FindString(message); date != "" {
		extraction.Date = date
	}
	if strings.Contains(text, "tomorrow") {
		extraction.Date = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	timeRe := regexp.MustCompile(`\b([01]?\d|2[0-3]):[0-5]\d\b`)
	if value := timeRe.FindString(message); value != "" {
		extraction.Time = value
	}

	if extraction.ServiceID == "" && state.Step == "collect_service" {
		extraction.ServiceQuery = fallbackServiceQuery(text)
	}

	return extraction
}

func fallbackServiceQuery(text string) string {
	replacer := strings.NewReplacer(
		"i want to", "",
		"i want", "",
		"book me", "",
		"book", "",
		"appointment", "",
		"dentist", "",
		"dental", "",
		"tomorrow", "",
		"today", "",
	)
	text = replacer.Replace(text)
	dateRe := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	text = dateRe.ReplaceAllString(text, "")
	timeRe := regexp.MustCompile(`\b([01]?\d|2[0-3]):[0-5]\d\b`)
	text = timeRe.ReplaceAllString(text, "")
	return strings.TrimSpace(text)
}
