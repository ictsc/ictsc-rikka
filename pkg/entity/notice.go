package entity

import (
	"time"
)

type Notice struct {
	Base

	SourceId string `json:"source_id"`
	Title    string `json:"title"`
	Body     string `json:"body,omitempty"`
	Draft    bool   `json:"draft"`
}

type NoticeWithSyncTime struct {
	Notice
	UpdatedAt time.Time
}
