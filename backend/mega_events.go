package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EventResponse struct {
	Events []struct {
		ActivityID             int    `json:"activity_id"`
		ActivityName           string `json:"activity_name"`
		AffectedByReplacements bool   `json:"affected_by_replacements"`
		Attendees              int    `json:"attendees"`
		BookingInfo            struct {
			Available     bool `json:"available"`
			IHaveBooked   bool `json:"i_have_booked"`
			LoginRequired bool `json:"login_required"`
			PassRequired  bool `json:"pass_required"`
			Places        struct {
				Booked int `json:"booked"`
				Total  int `json:"total"`
			} `json:"places"`
			ProductsInfo struct {
			} `json:"products_info"`
			SoldOut bool `json:"sold_out"`
			TooLate bool `json:"too_late"`
			TooSoon bool `json:"too_soon"`
		} `json:"booking_info"`
		BookingWaitingList     bool `json:"booking_waiting_list"`
		BookingWaitingListInfo struct {
		} `json:"booking_waiting_list_info,omitempty"`
		BookingsAppendable bool          `json:"bookings_appendable"`
		BookingsCancelable bool          `json:"bookings_cancelable"`
		BookingsEditable   bool          `json:"bookings_editable"`
		BookingsListable   bool          `json:"bookings_listable"`
		CategoriesIds      []int         `json:"categories_ids"`
		Color              string        `json:"color"`
		Conflict           []interface{} `json:"conflict"`
		Duration           int           `json:"duration"`
		End                string        `json:"end"`
		Hour               time.Time     `json:"hour"`
		ID                 string        `json:"id"`
		Instructors        []struct {
			ActivitiesIds []int  `json:"activities_ids"`
			Avatar        string `json:"avatar"`
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Surname       string `json:"surname"`
		} `json:"instructors"`
		Mobile struct {
			Color        string      `json:"color"`
			Duration     string      `json:"duration"`
			Icon         string      `json:"icon"`
			MonthDay     string      `json:"month_day"`
			RoomOrder    interface{} `json:"roomOrder"`
			StartTime    string      `json:"start_time"`
			Subtitle     string      `json:"subtitle"`
			Title        string      `json:"title"`
			WeekDay      string      `json:"week_day"`
			WeekDayShort string      `json:"week_day_short"`
		} `json:"mobile"`
		Replacer   interface{} `json:"replacer"`
		ResourceID int         `json:"resourceId"`
		Room       string      `json:"room"`
		RoomObj    struct {
			Icon struct {
				CSSClass string `json:"css_class"`
			} `json:"icon"`
			ID int `json:"id"`
		} `json:"room_obj"`
		Rotating                bool          `json:"rotating"`
		SessionID               int           `json:"session_id"`
		SpecialClass            interface{}   `json:"special_class"`
		Start                   string        `json:"start"`
		StartEditable           bool          `json:"startEditable"`
		SubstitutionInstructors []interface{} `json:"substitution_instructors"`
		Target                  int           `json:"target"`
		Title                   string        `json:"title"`
		TypeOfClass             interface{}   `json:"type_of_class"`
		Wday                    int           `json:"wday"`
	} `json:"events"`
}

func FetchMegaEvents(megaCreds *MegaCreds) (*EventResponse, error) {
	now := time.Now()
	oneWeekLater := now.AddDate(0, 0, 7)
	end := oneWeekLater.Unix()

	url1 := fmt.Sprintf("https://app.gym-up.com/ws/v2/event_sessions_public/%s/timetable", megaCreds.publicSession)
	url2 := fmt.Sprintf("https://app.gym-up.com/ws/v2/event_sessions_public/%s/timetable?start=%d", megaCreds.publicSession, end)

	// First API call
	req1, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		return nil, errors.New("failed to create first API request")
	}
	req1.Header.Add("X-Csrf-Token", megaCreds.csrfToken)
	req1.Header.Add("Cookie", fmt.Sprintf("_gymtoken=%s", megaCreds.authToken))

	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		return nil, errors.New("failed to make first API call")
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL1: %s", resp1.Status)
	}

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return nil, errors.New("failed to read response body from first API call")
	}

	var eventResponse1 EventResponse
	if err := json.Unmarshal(body1, &eventResponse1); err != nil {
		return nil, errors.New("failed to unmarshal JSON from first API call")
	}

	// Second API call
	req2, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		return nil, errors.New("failed to create second API request")
	}
	req2.Header.Add("X-Csrf-Token", megaCreds.csrfToken)
	req2.Header.Add("Cookie", fmt.Sprintf("_gymtoken=%s", megaCreds.authToken))

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, errors.New("failed to make second API call")
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL2: %s", resp2.Status)
	}

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, errors.New("failed to read response body from second API call")
	}

	var eventResponse2 EventResponse
	if err := json.Unmarshal(body2, &eventResponse2); err != nil {
		return nil, errors.New("failed to unmarshal JSON from second API call")
	}

	// Concatenate events
	concatenatedEvents := append(eventResponse1.Events, eventResponse2.Events...)

	// Create the final response
	return &EventResponse{
		Events: concatenatedEvents,
	}, nil
}

// HandleMegaEvents handles the API call to fetch event sessions and returns the concatenated JSON response
func HandleMegaEvents(c *fiber.Ctx) error {
	sessionID := GetSessionID(c)
	megaCreds, _ := Sessions.Get(sessionID)

	// Create the final response
	finalResponse, err := FetchMegaEvents(megaCreds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(finalResponse)
}

func GetClassDate(classId int, megaCreds *MegaCreds) (*time.Time, error) {
	events, err := FetchMegaEvents(megaCreds)
	if err != nil {
		return nil, err
	}

	for _, event := range events.Events {
		if event.SessionID == classId {
			return &event.Hour, nil
		}
	}

	return nil, errors.New("class not found")
}
