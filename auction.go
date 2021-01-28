package auction

import "github.com/google/uuid"

type Auction struct {
	ID        string
	AssetName string // company bought or sold
	Seller    string
	Bidder    string
	IntelUrl  string
	Status    string
}

func New(name string, seller string, bidder string, url string, status string) Auction {
	return Auction{
		ID:        uuid.New().String(),
		AssetName: name,
		Seller:    seller,
		Bidder:    bidder,
		IntelUrl:  url,
		Status:    status,
	}
}
