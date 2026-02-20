package models

import "time"

type Deal struct {
	DealID      int        `json:"deal_id" db:"deal_id"`
	RequestID   int        `json:"request_id" db:"request_id"`
	ExecutorID  int        `json:"executor_id" db:"executor_id"`
	AgreedPrice float64    `json:"agreed_price" db:"agreed_price"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ClosedAt    *time.Time `json:"closed_at" db:"closed_at"`
}
