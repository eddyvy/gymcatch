package backend

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

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

func GetCreds() (string, string, error) {
	sessionID, err := getSessionID()
	if err != nil {
		return "", "", err
	}
	csrfToken, err := getCSRFToken(sessionID)
	if err != nil {
		return "", "", err
	}
	return sessionID, csrfToken, nil
}
