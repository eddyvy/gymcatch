package backend

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

var interval = 30 * time.Second
var tryingToInscribeMap = make(map[string]map[int]bool)

func StartCronInscribe(classId int, megaCreds *MegaCreds, startDate, endDate time.Time) {
	if _, exists := tryingToInscribeMap[megaCreds.authToken]; !exists {
		tryingToInscribeMap[megaCreds.authToken] = make(map[int]bool)
	}

	if tryingToInscribeMap[megaCreds.authToken][classId] {
		return
	}

	c := cron.New()

	// Define the job
	job := func() {
		fmt.Println(classId, "--->Executing job")
		now := time.Now()
		stopIt := false
		if now.After(startDate) && now.Before(endDate) {
			fmt.Println(classId, "--->Trying to Inscribe")
			success := MegaInscribe(classId, megaCreds)
			if success {
				fmt.Println(classId, "--->Successfully inscribed")
				stopIt = true
			} else {
				fmt.Println(classId, "--->Failed to inscribe")
			}
		} else if now.After(endDate) {
			fmt.Println(classId, "--->Too late to inscribe")
			stopIt = true
		} else if now.Before(startDate) {
			fmt.Println(classId, "--->Too early to inscribe")
		}

		if stopIt {
			fmt.Println(classId, "--->Stopping job")
			tryingToInscribeMap[megaCreds.authToken][classId] = false
			c.Stop()
		}
	}

	tryingToInscribeMap[megaCreds.authToken][classId] = true

	// Add the job to the cron scheduler
	c.AddFunc(fmt.Sprintf("@every %s", interval), job)

	// Start the cron scheduler
	c.Start()
}

// RequestBody represents the structure of the request body
type RequestBody struct {
	ClassId   int    `json:"classId"`
	ClassDate string `json:"classDate"`
}

func HandleInscribe(c *fiber.Ctx) error {
	sessionID := GetSessionID(c)
	megaCreds, _ := Sessions.Get(sessionID)

	megaCreds.LoadCreds()

	classIdParam := c.Params("classId")
	classId, err := strconv.Atoi(classIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid classId",
		})
	}

	classDate, err := GetClassDate(classId, megaCreds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set startDate to 2 days before classDate at 20:59
	startDate := classDate.AddDate(0, 0, -2).Add(time.Hour*20 + time.Minute*59)
	endDate := classDate.Add(-15 * time.Minute) // End 15 minutes before class starts

	if time.Now().After(endDate) {
		fmt.Println(classId, "--->Too late to inscribe")
		return c.JSON(fiber.Map{
			"success": false,
		})
	}

	// Start the cron job
	StartCronInscribe(classId, megaCreds, startDate, endDate)

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func HandleGetInscribedClasses(c *fiber.Ctx) error {
	// Return tryingToInscribeMap keys that are true
	sessionID := GetSessionID(c)
	megaCreds, _ := Sessions.Get(sessionID)

	if _, exists := tryingToInscribeMap[megaCreds.authToken]; !exists {
		tryingToInscribeMap[megaCreds.authToken] = make(map[int]bool)
	}

	inscribedClasses := []int{}
	for classId, inscribing := range tryingToInscribeMap[megaCreds.authToken] {
		if inscribing {
			inscribedClasses = append(inscribedClasses, classId)
		}
	}
	return c.JSON(fiber.Map{
		"inscribedClasses": inscribedClasses,
	})
}
