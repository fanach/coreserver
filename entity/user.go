package entity

import "time"

// User is the struct for registered user
type User struct {
	// required
	Username string   `json:"username"`
	Password Password `json:"password"`
	// optional
	WeChatID string `json:"wechat_id"`
	Type     string `json:"type"`
	Email    string `json:"email"`
	// generated
	ID         string    `json:"id"`
	RegTime    time.Time `json:"reg_time"`
	UpdateTime time.Time `json:"update_time"`
}
