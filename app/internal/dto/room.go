package dto

import "time"

type Room struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	ItemID    string    `json:"itemId"`
	BuyerID   string    `json:"buyerId"`
	SellerID  string    `json:"sellerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
