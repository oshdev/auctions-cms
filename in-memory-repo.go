package auction

type InMemoryRepo struct {
	auctions []Auction
}

func NewInMemoryRepo() *InMemoryRepo {
	var auctions []Auction
	auctions = append(auctions, New("Super Company", "Riya", "Osh", "https://mergermarket.com/secret", "MIA"))
	return &InMemoryRepo{auctions: auctions}
}

func (i *InMemoryRepo) GetAuction(id string) Auction {
	for _, auction := range i.auctions {
		if auction.ID == id {
			return auction
		}
	}
	return New("wtf this doesnt exist", "seller", "bidder", "url", "status")
}

func (i *InMemoryRepo) EditAuction(id string, assetName string) {
	for index := range i.auctions {
		if i.auctions[index].ID == id {
			i.auctions[index].AssetName = assetName
		}
	}
}

func (i *InMemoryRepo) GetAuctions() []Auction {
	return i.auctions
}

func (i *InMemoryRepo) AddAuction(name string) {
	i.auctions = append(i.auctions, New(name, "seller", "bidder", "url", "status"))
}

func (i *InMemoryRepo) DeleteAuction(id string) {
	var newList []Auction
	for _, auction := range i.auctions {
		if auction.ID != id {
			newList = append(newList, auction)
		}
	}
	i.auctions = newList
}
