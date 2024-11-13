package postgres

import (
    "context"
    "database/sql"
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/errors"
)

type BidRepository struct {
    db *sql.DB
}

func NewBidRepository(db *sql.DB) *BidRepository {
    return &BidRepository{db: db}
}

func (r *BidRepository) Create(ctx context.Context, bid *entity.Bid) error {
    query := `
        INSERT INTO bids (auction_id, user_id, amount)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at`

    err := r.db.QueryRowContext(
        ctx,
        query,
        bid.AuctionID,
        bid.UserID,
        bid.Amount,
    ).Scan(&bid.ID, &bid.CreatedAt, &bid.UpdatedAt)

    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to create bid", err)
    }

    return nil
}

func (r *BidRepository) GetByID(ctx context.Context, id int64) (*entity.Bid, error) {
    query := `
        SELECT id, auction_id, user_id, amount, created_at, updated_at
        FROM bids
        WHERE id = $1`

    bid := &entity.Bid{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &bid.ID,
        &bid.AuctionID,
        &bid.UserID,
        &bid.Amount,
        &bid.CreatedAt,
        &bid.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New(errors.ErrorTypeNotFound, "bid not found", nil)
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get bid", err)
    }

    return bid, nil
}

func (r *BidRepository) GetByAuctionID(ctx context.Context, auctionID int64) ([]*entity.Bid, error) {
    query := `
        SELECT id, auction_id, user_id, amount, created_at, updated_at
        FROM bids
        WHERE auction_id = $1
        ORDER BY amount DESC`

    rows, err := r.db.QueryContext(ctx, query, auctionID)
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get bids", err)
    }
    defer rows.Close()

    var bids []*entity.Bid
    for rows.Next() {
        bid := &entity.Bid{}
        err := rows.Scan(
            &bid.ID,
            &bid.AuctionID,
            &bid.UserID,
            &bid.Amount,
            &bid.CreatedAt,
            &bid.UpdatedAt,
        )
        if err != nil {
            return nil, errors.New(errors.ErrorTypeInternal, "failed to scan bid", err)
        }
        bids = append(bids, bid)
    }

    return bids, nil
}

func (r *BidRepository) List(ctx context.Context, auctionID int64, offset, limit int) ([]*entity.Bid, int64, error) {
    query := `
        SELECT id, auction_id, user_id, amount, created_at, updated_at
        FROM bids
        WHERE auction_id = $1
        ORDER BY amount DESC
        LIMIT $2 OFFSET $3`

    rows, err := r.db.QueryContext(ctx, query, auctionID, limit, offset)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to list bids", err)
    }
    defer rows.Close()

    var bids []*entity.Bid
    for rows.Next() {
        bid := &entity.Bid{}
        err := rows.Scan(
            &bid.ID,
            &bid.AuctionID,
            &bid.UserID,
            &bid.Amount,
            &bid.CreatedAt,
            &bid.UpdatedAt,
        )
        if err != nil {
            return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to scan bid", err)
        }
        bids = append(bids, bid)
    }

    var total int64
    err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM bids WHERE auction_id = $1", auctionID).Scan(&total)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to count bids", err)
    }

    return bids, total, nil
}

func (r *BidRepository) GetUniqueParticipantsByAuctionID(ctx context.Context, auctionID int64) ([]int64, error) {
    query := `
        SELECT DISTINCT user_id 
        FROM bids 
        WHERE auction_id = $1`

    rows, err := r.db.QueryContext(ctx, query, auctionID)
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get auction participants", err)
    }
    defer rows.Close()

    var participants []int64
    for rows.Next() {
        var userID int64
        if err := rows.Scan(&userID); err != nil {
            return nil, errors.New(errors.ErrorTypeInternal, "failed to scan participant", err)
        }
        participants = append(participants, userID)
    }

    return participants, nil
}
