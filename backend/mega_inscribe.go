package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type EventInfoType struct {
	Texts struct {
		SoldOut   string `json:"sold_out"`
		Unstarted string `json:"unstarted"`
		Finished  string `json:"finished"`
	} `json:"texts"`
	Data []struct {
		ID          int `json:"id"`
		BookingInfo struct {
			Available   bool `json:"available"`
			TooSoon     bool `json:"too_soon"`
			TooLate     bool `json:"too_late"`
			SoldOut     bool `json:"sold_out"`
			IHaveBooked bool `json:"i_have_booked"`
			Places      struct {
				Booked int `json:"booked"`
				Total  int `json:"total"`
			} `json:"places"`
			LoginRequired bool `json:"login_required"`
			ProductsInfo  struct {
			} `json:"products_info"`
			PassRequired bool `json:"pass_required"`
		} `json:"booking_info"`
		BookingWaitingListInfo struct {
			Available   bool `json:"available"`
			TooSoon     bool `json:"too_soon"`
			TooLate     bool `json:"too_late"`
			SoldOut     bool `json:"sold_out"`
			IHaveBooked bool `json:"i_have_booked"`
			Places      struct {
				Booked            int `json:"booked"`
				NotifiedAvailable int `json:"notified_available"`
				Total             int `json:"total"`
			} `json:"places"`
			LoginRequired bool `json:"login_required"`
			ProductsInfo  struct {
			} `json:"products_info"`
		} `json:"booking_waiting_list_info"`
	} `json:"data"`
}

func MegaInscribe(classId int, megaCreds *MegaCreds) bool {
	err := megaCreds.LoadCreds()
	if err != nil {
		fmt.Println(classId, "--->Error loading credentials", err)
		megaCreds.RemoveCreds()
		return false
	}
	eventInfo, err := GetClassInfo(classId, megaCreds)
	if err != nil {
		fmt.Println(classId, "--->Error getting class info", err)
		megaCreds.RemoveCreds()
		return false
	}

	if eventInfo.Data[0].BookingInfo.IHaveBooked {
		return true
	}

	if eventInfo.Data[0].BookingInfo.Available {
		if err := InscribeToClass(classId, megaCreds); err != nil {
			fmt.Println(classId, "--->Error inscribing to class", err)
			megaCreds.RemoveCreds()
			return false
		}
		return true
	}

	return false
}

func GetClassInfo(classId int, megaCreds *MegaCreds) (*EventInfoType, error) {
	url := "https://app.gym-up.com/ws/v2/event_sessions_public/" + megaCreds.publicSession + "/booking_info"
	method := "POST"

	payload := strings.NewReader("event_session_ids=" + strconv.Itoa(classId))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Csrf-Token", megaCreds.csrfToken)
	req.Header.Add("Cookie", "_gymtoken="+megaCreds.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var eventInfo EventInfoType
	if err := json.Unmarshal(body, &eventInfo); err != nil {
		return nil, err
	}

	return &eventInfo, nil
}

type ResponseInscribe struct {
	Success bool `json:"success"`
}

func InscribeToClass(classId int, megaCreds *MegaCreds) error {
	fmt.Println("INSC", classId, megaCreds)
	url := "https://app.gym-up.com/api/v1/bookings"
	method := "POST"

	payload := strings.NewReader("gym_token=" + megaCreds.publicSession + "&booking%5Bevent_session_id%5D=" + strconv.Itoa(classId))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("X-Csrf-Token", megaCreds.csrfToken)
	req.Header.Add("Cookie", "_gymtoken="+megaCreds.authToken+"; _gymapp="+megaCreds.gymAppToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("INSC RESPONSE status", res.Status)
	fmt.Println("INSC RESPONSE", string(body))

	var response ResponseInscribe
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("failed to inscribe to class")
	}

	return nil
}
