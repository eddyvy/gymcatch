package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ResponseInscribe struct {
	Success bool `json:"success"`
}

type ResponseInscribeError struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func MegaInscribe(classId int, megaCreds *MegaCreds) bool {
	err := megaCreds.LoadCreds()
	if err != nil {
		fmt.Println(classId, "--->Error loading credentials", err)
		megaCreds.RemoveCreds()
		return false
	}
	eventInfo, err := GetClassInfo([]int{classId}, megaCreds)
	if err != nil {
		fmt.Println(classId, "--->Error getting class info", err)
		megaCreds.RemoveCreds()
		return false
	}

	if eventInfo.Data[0].BookingInfo.IHaveBooked {
		fmt.Println(classId, "--->Already booked it")
		return true
	}

	if eventInfo.Data[0].BookingInfo.Available {
		if err := InscribeToClass(classId, megaCreds); err != nil {
			fmt.Println(classId, "--->Error inscribing to class", err)
			megaCreds.RemoveCreds()
			return false
		}
		return true
	} else {
		fmt.Println(classId, "--->Class not available")
		return false
	}
}

func InscribeToClass(classId int, megaCreds *MegaCreds) error {
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

	var response ResponseInscribe
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Success {
		var responseError ResponseInscribeError
		if err := json.Unmarshal(body, &responseError); err != nil {
			return fmt.Errorf("failed to inscribe to class: %v", responseError.Errors)
		} else {
			return fmt.Errorf("failed to inscribe to class")
		}
	}

	return nil
}
