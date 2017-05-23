package entity

import "time"

const (
	// StatPaying is the status of order which is waiting for paying
	StatPaying = "paying"
	// StatCanceled is the status of canceled order
	StatCanceled = "canceled"
	// StatFinished is the final status of a completed order
	StatFinished = "finished"
)

// Order is purchase of product
type Order struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	ProductID string    `json:"product_id"`
	BeginTime time.Time `json:"begin_time"`
	EndTime   time.Time `json:"end_time"`
	Stat      string    `json:"stat"`
}
