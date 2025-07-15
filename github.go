package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const GithubUrl string = "https://api.github.com/notifications"

type GithubPoller struct {
	lastModified *string
	client       *http.Client
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func NewGithubPoller(accessToken string) *GithubPoller {
	client := &http.Client{}
	transport := http.DefaultTransport

	client.Transport = roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
			req.Header.Set("Accept", "application/vnd.github.v3+json")
			return transport.RoundTrip(req)
		},
	)

	return &GithubPoller{
		lastModified: nil,
		client:       client,
	}
}

func (p *GithubPoller) checkIfNotificationsAsChanged() (bool, error) {
	client := p.client
	var url = GithubUrl

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	if p.lastModified != nil {
		req.Header.Set("If-Modified-Since", *p.lastModified)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		fmt.Println("No changes")
		return false, nil
	}

	fmt.Println("New changes")
	return true, nil
}

func (p *GithubPoller) getNotifications() {
	client := p.client
	var url = GithubUrl
	if p.lastModified != nil {
		// Parse the RFC1123 format and convert to ISO 8601
		parsedTime, err := time.Parse(time.RFC1123, *p.lastModified)
		if err != nil {
			panic(err)
		}
		iso8601Time := parsedTime.Format(time.RFC3339)
		url = fmt.Sprintf("%s?since=%s", GithubUrl, iso8601Time)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Store the Last-Modified header for future requests
	if lastMod := resp.Header.Get("Last-Modified"); lastMod != "" {
		p.lastModified = &lastMod
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var notifications []map[string]interface{}
	if err := json.Unmarshal(body, &notifications); err != nil {
		panic(err)
	}

	fmt.Printf("Found %d notifications\n", len(notifications))
	for i, notification := range notifications {
		if subject, ok := notification["subject"].(map[string]any); ok {
			if title, ok := subject["title"].(string); ok {
				fmt.Printf("Notification %d: %s\n", i, title)
			}
		}
	}
}

func getTokenFromGithubCLI() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
