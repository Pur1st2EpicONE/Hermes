package v1

// CreateCommentV1 represents the request payload for creating a comment.
// ParentID is optional and used for nested comments.
type CreateCommentV1 struct {
	ParentID *int64 `json:"parent_id,omitempty"`
	Content  string `json:"content"`
	Author   string `json:"author"`
}
