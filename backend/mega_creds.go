package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type MegaCreds struct {
	mu            sync.RWMutex
	email         string
	publicSession string
	csrfToken     string
	authToken     string
	gymAppToken   string
}

func NewMegaCreds(email string) *MegaCreds {
	return &MegaCreds{
		email:         email,
		publicSession: "",
		csrfToken:     "",
		authToken:     "",
		gymAppToken:   "",
	}
}

func (m *MegaCreds) LoadCreds() error {
	if m.publicSession != "" && m.csrfToken != "" && m.authToken != "" {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID, err := getSessionID()
	if err != nil {
		return err
	}
	csrfToken, err := getCSRFToken(sessionID)
	if err != nil {
		return err
	}
	authToken, gymAppToken, err := getAuthTokens(sessionID, csrfToken, m.email)
	if err != nil {
		return err
	}

	m.publicSession = sessionID
	m.csrfToken = csrfToken
	m.authToken = authToken
	m.gymAppToken = gymAppToken

	return nil
}

func (m *MegaCreds) RemoveCreds() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.publicSession = ""
	m.csrfToken = ""
}

func (m *MegaCreds) GetCreds() (string, string, string, string, error) {
	err := m.LoadCreds()
	if err != nil {
		return "", "", "", "", err
	}
	return m.publicSession, m.csrfToken, m.authToken, m.gymAppToken, nil
}

// extractSessionID extracts the session ID from the HTML
func extractSessionID(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var sessionID string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "iframe" {
			for _, a := range n.Attr {
				if a.Key == "src" && strings.Contains(a.Val, "https://app.gym-up.com/ws/v2/event_sessions_public/index/") {
					parts := strings.Split(a.Val, "https://app.gym-up.com/ws/v2/event_sessions_public/index/")
					parts = strings.Split(parts[1], "?")
					sessionID = parts[0]
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if sessionID == "" {
		return "", fmt.Errorf("session ID not found")
	}
	return sessionID, nil
}

func extractCSRFToken(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var csrfToken string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "csrf-token" {
					name = a.Val
				}
				if a.Key == "content" {
					content = a.Val
				}
			}
			if name == "csrf-token" {
				csrfToken = content
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if csrfToken == "" {
		return "", fmt.Errorf("CSRF token not found")
	}
	return csrfToken, nil
}

func getSessionID() (string, error) {
	resp, err := http.Get("https://megasportcentre.com/horario-clases")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return extractSessionID(string(body))
}

func getCSRFToken(sessionID string) (string, error) {
	resp, err := http.Get("https://app.gym-up.com/ws/v2/event_sessions_public/index/" + sessionID)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return extractCSRFToken(string(body))
}

type MegaCredsLoginResponse struct {
	Success bool `json:"success"`
	User    struct {
		Avatar      string `json:"avatar"`
		ShortName   string `json:"short_name"`
		LogoutText  string `json:"logout_text"`
		AccessToken string `json:"access_token"`
	} `json:"user"`
}

func getAuthTokens(publicSession, csrfToken, email string) (string, string, error) {
	url := "https://app.gym-up.com/ws/v2/event_sessions_public/" + publicSession + "/login"
	payload := strings.NewReader("email=" + email + "&password=" + os.Getenv("MEGA_PASSWORD"))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Csrf-Token", csrfToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch URL: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var loginResponse MegaCredsLoginResponse
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return "", "", err
	}

	// Get the Set-Cookie headers
	cookies := resp.Header["Set-Cookie"]
	var gymAppToken string
	for _, cookie := range cookies {
		if strings.Contains(cookie, "_gymapp") {
			parts := strings.Split(cookie, ";")
			gymAppToken = strings.Split(parts[0], "_gymapp=")[1]
			break
		}
	}

	return loginResponse.User.AccessToken, gymAppToken, nil
}
