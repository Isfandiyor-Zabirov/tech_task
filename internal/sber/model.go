package sber

type SendRequestBody struct {
	Search string `json:"search"`
}

type Response struct {
	Keyword  string    `json:"keyword"`
	Auctions []Auction `json:"auctions"`
}

type Auction struct {
	AuctionName   string `json:"auction_name"`
	AuctionNumber string `json:"auction_number"`
}
