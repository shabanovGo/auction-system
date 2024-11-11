package handler

type Handlers struct {
    UserHandler    *UserHandler
    // AuctionHandler *AuctionHandler
    LotHandler     *LotHandler
}

func NewHandlers(
    userHandler *UserHandler,
    // auctionHandler *AuctionHandler,
    lotHandler *LotHandler,
) *Handlers {
    return &Handlers{
        UserHandler: userHandler,
        // AuctionHandler: auctionHandler,
        LotHandler:     lotHandler,
    }
}
