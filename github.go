package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
)

const GithubUrl string = "https://api.github.com/notifications"

type GithubPoller struct {
	client *http.Client
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
		client: client,
	}
}

func (p *GithubPoller) getLastModifiedPath() string {
	return filepath.Join(xdg.DataHome, "gh_notify", "last_modified")
}

func (p *GithubPoller) readLastModified() *string {
	path := p.getLastModifiedPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	lastMod := string(data)
	return &lastMod
}

func (p *GithubPoller) writeLastModified(lastModified string) error {
	path := p.getLastModifiedPath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(lastModified), 0644)
}

func (p *GithubPoller) checkIfNotificationsAsChanged() (bool, error) {
	client := p.client
	var url = GithubUrl

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	lastModified := p.readLastModified()
	if lastModified != nil {
		req.Header.Set("If-Modified-Since", *lastModified)
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
	lastModified := p.readLastModified()
	if lastModified != nil {
		// Parse the RFC1123 format and convert to ISO 8601
		parsedTime, err := time.Parse(time.RFC1123, *lastModified)
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
		if err := p.writeLastModified(lastMod); err != nil {
			fmt.Printf("Error writing last modified time: %v\n", err)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var notifications []map[string]any
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
