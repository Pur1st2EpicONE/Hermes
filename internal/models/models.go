package models

import "time"

type Comment struct {
	ID        int64      `json:"id"`
	ParentID  *int64     `json:"parent_id,omitempty"`
	Content   string     `json:"content"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Children  []*Comment `json:"children,omitempty"`
}

type QueryParams struct {
	ParentID *int64
	Page     int
	Limit    int
	Sort     string
	Offset   int
}
