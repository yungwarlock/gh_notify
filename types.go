package main

import "time"

type Notification struct {
	ID              string     `json:"id"`
	Unread          bool       `json:"unread"`
	Reason          string     `json:"reason"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastReadAt      *time.Time `json:"last_read_at"`
	Subject         Subject    `json:"subject"`
	Repository      Repository `json:"repository"`
	URL             string     `json:"url"`
	SubscriptionURL string     `json:"subscription_url"`
}

type Subject struct {
	Title            string  `json:"title"`
	URL              *string `json:"url"`
	LatestCommentURL *string `json:"latest_comment_url"`
	Type             string  `json:"type"`
}

type Repository struct {
	ID          int64   `json:"id"`
	NodeID      string  `json:"node_id"`
	Name        string  `json:"name"`
	FullName    string  `json:"full_name"`
	Private     bool    `json:"private"`
	Owner       User    `json:"owner"`
	HTMLURL     string  `json:"html_url"`
	Description *string `json:"description"`
	Fork        bool    `json:"fork"`
	URL         string  `json:"url"`
}

type User struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}
