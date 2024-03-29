package response

import (
	"time"
)

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`

	Status       int    `json:"status"`
	RegisterType string `json:"register_type"`
	IpAddress    string `json:"ip_address"` // ip host
	IpSource     string `json:"ip_source"`  // ip 源

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Roles     []*RoleDTO `json:"roles"`
}

type UserAreaDTO struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type UserMenuDTO struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	Component string        `json:"component"`
	Icon      string        `json:"icon"`
	Rank      int           `json:"rank"`
	IsHidden  int           `json:"is_hidden"`
	Children  []UserMenuDTO `json:"children"`
}
