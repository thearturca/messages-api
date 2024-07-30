package db

import (
	"time"
)

type Message struct {
	Id          string     `db:"id" json:"id" goqu:"pk,skipupdate,skipinsert,defaultifempty"`
	Text        string     `db:"text" json:"text" goqu:"defaultifempty"`
	IsProcessed bool       `db:"is_processed" json:"isProcessed" goqu:"skipinsert,defaultifempty"`
	ProcessedAt *time.Time `db:"processed_at" json:"processedAt" goqu:"defaultifempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt" goqu:"skipupdate,skipinsert,defaultifempty"`
}
