package models

import "time"

type Comment struct {
	ID        int64
	ParentID  *int64
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Children  []*Comment
}

type QueryParams struct {
	ParentID *int64
	Page     int
	Limit    int
	Sort     string
	Offset   int
}
