package models

import "time"

type PendingASREQ struct {
	Usuario   string
	Realm     string
	EType     int32
	Salt      string
	Timestamp time.Time
	SrcIP     string
	DstIP     string
}
