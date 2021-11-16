package model

// ToDo: add swagger docs, check block.go
type BoardMember struct {
    BoardID         string                 `json:"board_id"`
    UserID          string                 `json:"user_id"`
    Roles           string                 `json:"roles"`
	SchemeAdmin     bool                   `json:"scheme_admin"`
	SchemeEditor    bool                   `json:"scheme_editor"`
    SchemeCommenter bool                   `json:"scheme_commenter"`
    SchemeViewer    bool                   `json:"scheme_viewer"`
    LastViewedAt    int64                  `json:"last_viewed_at"`
    MentionCount    int64                  `json:"mention_count"`
    NotifyProps     map[string]interface{} `json:"notify_props"`
	// ToDo: this should probably go
    // ExplicitRoles   string                 `json:"explicit_roles"`
}
