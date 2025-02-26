package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EventResponse represents the structure of the JSON response
type EventResponse struct {
	Events []interface{} `json:"events"`
}

// HandleMegaEvents handles the API call to fetch event sessions and returns the concatenated JSON response
func HandleMegaEvents(c *fiber.Ctx) error {
	sessionID := GetSessionID(c)
	megaCreds, _ := Sessions.Get(sessionID)

	now := time.Now()
	oneWeekLater := now.AddDate(0, 0, 7)
	end := oneWeekLater.Unix()

	url1 := fmt.Sprintf("https://app.gym-up.com/ws/v2/event_sessions_public/%s/timetable", megaCreds.publicSession)
	url2 := fmt.Sprintf("https://app.gym-up.com/ws/v2/event_sessions_public/%s/timetable?start=%d", megaCreds.publicSession, end)

	// First API call
	req1, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create first API request",
		})
	}
	req1.Header.Add("X-Csrf-Token", megaCreds.csrfToken)
	req1.Header.Add("Cookie", fmt.Sprintf("_gymtoken=%s", megaCreds.authToken))

	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to make first API call",
		})
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to fetch URL1: %s", resp1.Status),
		})
	}

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read response body from first API call",
		})
	}

	var eventResponse1 EventResponse
	if err := json.Unmarshal(body1, &eventResponse1); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to unmarshal JSON from first API call",
		})
	}

	// Second API call
	req2, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create second API request",
		})
	}
	req2.Header.Add("X-Csrf-Token", megaCreds.csrfToken)
	req2.Header.Add("Cookie", fmt.Sprintf("_gymtoken=%s", megaCreds.authToken))

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to make second API call",
		})
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to fetch URL2: %s", resp2.Status),
		})
	}

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read response body from second API call",
		})
	}

	var eventResponse2 EventResponse
	if err := json.Unmarshal(body2, &eventResponse2); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to unmarshal JSON from second API call",
		})
	}

	// Concatenate events
	concatenatedEvents := append(eventResponse1.Events, eventResponse2.Events...)

	// Create the final response
	finalResponse := EventResponse{
		Events: concatenatedEvents,
	}

	return c.JSON(finalResponse)
}
