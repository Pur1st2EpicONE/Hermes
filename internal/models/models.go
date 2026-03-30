// Package models defines core domain entities and data transfer structures
// used across the application layers.
package models

import "time"

// Comment represents a user comment entity.
type Comment struct {
	ID        int64      `json:"id"`                  // ID is the unique identifier of the comment.
	ParentID  *int64     `json:"parent_id,omitempty"` // ParentID references the parent comment (nil for root comments).
	Content   string     `json:"content"`             // Content is the textual body of the comment.
	Author    string     `json:"author"`              // Author is the name or identifier of the comment creator.
	CreatedAt time.Time  `json:"created_at"`          // CreatedAt is the timestamp when the comment was created.
	UpdatedAt time.Time  `json:"updated_at"`          // UpdatedAt is the timestamp of the last update.
	Children  []*Comment `json:"children,omitempty"`  // Children contains nested replies to this comment.
}

// QueryParams encapsulates parameters for querying comments.
type QueryParams struct {
	ParentID *int64 // ParentID filters comments by parent (nil for top-level).
	Page     int    // Page is the page number (starting from 1).
	Limit    int    // Limit is the number of items per page.
	Sort     string // Sort defines ordering (e.g., created_at_desc or created_at_asc).
	Offset   int    // Offset is the calculated value for database queries.
}
