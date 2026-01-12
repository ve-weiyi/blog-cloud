package gitee

import (
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

type Userinfo struct {
	Id                int         `json:"id"`
	Login             string      `json:"login"`
	Name              string      `json:"name"`
	AvatarUrl         string      `json:"avatar_url"`
	Url               string      `json:"url"`
	HtmlUrl           string      `json:"html_url"`
	Remark            string      `json:"remark"`
	FollowersUrl      string      `json:"followers_url"`
	FollowingUrl      string      `json:"following_url"`
	GistsUrl          string      `json:"gists_url"`
	StarredUrl        string      `json:"starred_url"`
	SubscriptionsUrl  string      `json:"subscriptions_url"`
	OrganizationsUrl  string      `json:"organizations_url"`
	ReposUrl          string      `json:"repos_url"`
	EventsUrl         string      `json:"events_url"`
	ReceivedEventsUrl string      `json:"received_events_url"`
	Type              string      `json:"type"`
	Blog              interface{} `json:"blog"`
	Weibo             interface{} `json:"weibo"`
	Bio               string      `json:"bio"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	Stared            int         `json:"stared"`
	Watched           int         `json:"watched"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	Email             *string     `json:"email"`
}
