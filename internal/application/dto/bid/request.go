package bid

type PlaceBidRequest struct {
    AuctionID int64   `json:"auction_id"`
    UserID    int64   `json:"user_id"`
    Amount    float64 `json:"amount"`
}

type ListBidsRequest struct {
    AuctionID  int64 `json:"auction_id"`
    PageSize   int   `json:"page_size"`
    PageNumber int   `json:"page_number"`
}