package model

type BoardType string

const (
	BoardTypeOpen    BoardType = "O"
	BoardTypePrivate BoardType = "P"
)

// ToDo: add swagger docs, check block.go
type Board struct {
    ID        string    `json:"id"`
    TeamID    string    `json:"team_id"`
    ChannelID string    `json:"channel_id"`
    CreatorID string    `json:"creator_id"`
    Type      BoardType `json:"type"`

    Title           string `json:"title"`
    Description     string `json:"description"`
    Icon            string `json:"icon"`
    ShowDescription bool   `json:"show_description"`
    IsTemplate      bool   `json:"is_template"`

    Properties         map[string]interface{} `json:"properties" db:"-"`
    CardProperties     map[string]interface{} `json:"card_properties" db:"-"`
    ColumnCalculations map[string]interface{} `json:"column_calculations" db:"-"`

    CreateAt int64 `json:"create_at"`
    UpdateAt int64 `json:"update_at"`
    DeleteAt int64 `json:"delete_at"`
}
