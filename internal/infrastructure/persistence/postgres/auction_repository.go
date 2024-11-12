package postgres

import (
    "context"
    "database/sql"
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/errors"
    "fmt"
    "strings"
)

type AuctionRepository struct {
    db *sql.DB
}

func NewAuctionRepository(db *sql.DB) *AuctionRepository {
    return &AuctionRepository{db: db}
}

func (r *AuctionRepository) Create(ctx context.Context, auction *entity.Auction) error {
    query := `
        INSERT INTO auctions (
            lot_id, start_price, min_step, current_price,
            start_time, end_time, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

    err := r.db.QueryRowContext(
        ctx,
        query,
        auction.LotID,
        auction.StartPrice,
        auction.MinStep,
        auction.CurrentPrice,
        auction.StartTime,
        auction.EndTime,
        auction.Status,
    ).Scan(&auction.ID, &auction.CreatedAt, &auction.UpdatedAt)

    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to create auction", err)
    }

    return nil
}

func (r *AuctionRepository) GetByID(ctx context.Context, id int64) (*entity.Auction, error) {
    query := `
        SELECT id, lot_id, start_price, min_step, current_price,
               start_time, end_time, status, winner_id, winner_bid_id,
               created_at, updated_at
        FROM auctions
        WHERE id = $1`

    auction := &entity.Auction{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &auction.ID,
        &auction.LotID,
        &auction.StartPrice,
        &auction.MinStep,
        &auction.CurrentPrice,
        &auction.StartTime,
        &auction.EndTime,
        &auction.Status,
        &auction.WinnerID,
        &auction.WinnerBidID,
        &auction.CreatedAt,
        &auction.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New(errors.ErrorTypeNotFound, "auction not found", nil)
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get auction", err)
    }

    return auction, nil
}

func (r *AuctionRepository) Update(ctx context.Context, id int64, auction *entity.Auction) (*entity.Auction, error) {
    var queryParts []string
    var args []interface{}
    argPosition := 1

    if auction.StartPrice > 0 {
        queryParts = append(queryParts, fmt.Sprintf("start_price = $%d", argPosition))
        args = append(args, auction.StartPrice)
        argPosition++
    }
    if auction.MinStep > 0 {
        queryParts = append(queryParts, fmt.Sprintf("min_step = $%d", argPosition))
        args = append(args, auction.MinStep)
        argPosition++
    }
    if auction.CurrentPrice > 0 {
        queryParts = append(queryParts, fmt.Sprintf("current_price = $%d", argPosition))
        args = append(args, auction.CurrentPrice)
        argPosition++
    }
    if !auction.StartTime.IsZero() {
        queryParts = append(queryParts, fmt.Sprintf("start_time = $%d", argPosition))
        args = append(args, auction.StartTime)
        argPosition++
    }
    if !auction.EndTime.IsZero() {
        queryParts = append(queryParts, fmt.Sprintf("end_time = $%d", argPosition))
        args = append(args, auction.EndTime)
        argPosition++
    }
    if auction.Status != "" {
        queryParts = append(queryParts, fmt.Sprintf("status = $%d", argPosition))
        args = append(args, auction.Status)
        argPosition++
    }
    if auction.WinnerID != nil {
        queryParts = append(queryParts, fmt.Sprintf("winner_id = $%d", argPosition))
        args = append(args, auction.WinnerID)
        argPosition++
    }
    if auction.WinnerBidID != nil {
        queryParts = append(queryParts, fmt.Sprintf("winner_bid_id = $%d", argPosition))
        args = append(args, auction.WinnerBidID)
        argPosition++
    }

    queryParts = append(queryParts, "updated_at = CURRENT_TIMESTAMP")

    args = append(args, id)

    query := fmt.Sprintf(`
        UPDATE auctions
        SET %s
        WHERE id = $%d
        RETURNING created_at, updated_at`,
        strings.Join(queryParts, ", "),
        argPosition,
    )

    err := r.db.QueryRowContext(
        ctx,
        query,
        args...,
    ).Scan(&auction.CreatedAt, &auction.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, errors.New(errors.ErrorTypeNotFound, "auction not found", nil)
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to update auction", err)
    }

    auction.ID = id
    return auction, nil
}

func (r *AuctionRepository) List(ctx context.Context, offset, limit int, status *entity.AuctionStatus) ([]*entity.Auction, int64, error) {
    var args []interface{}
    var whereClause string

    if status != nil {
        whereClause = "WHERE status = $1"
        args = append(args, *status)
    }

    countQuery := `SELECT COUNT(*) FROM auctions ` + whereClause
    var total int64
    err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to count auctions", err)
    }

    query := `
        SELECT id, lot_id, start_price, min_step, current_price,
               start_time, end_time, status, winner_id, winner_bid_id,
               created_at, updated_at
        FROM auctions
        ` + whereClause + `
        ORDER BY created_at DESC
        LIMIT $` + fmt.Sprintf("%d", len(args)+1) + `
        OFFSET $` + fmt.Sprintf("%d", len(args)+2)

    args = append(args, limit, offset)
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to list auctions", err)
    }
    defer rows.Close()

    var auctions []*entity.Auction
    for rows.Next() {
        auction := &entity.Auction{}
        err := rows.Scan(
            &auction.ID,
            &auction.LotID,
            &auction.StartPrice,
            &auction.MinStep,
            &auction.CurrentPrice,
            &auction.StartTime,
            &auction.EndTime,
            &auction.Status,
            &auction.WinnerID,
            &auction.WinnerBidID,
            &auction.CreatedAt,
            &auction.UpdatedAt,
        )
        if err != nil {
            return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to scan auction", err)
        }
        auctions = append(auctions, auction)
    }

    return auctions, total, nil
}

func (r *AuctionRepository) GetByLotID(ctx context.Context, lotID int64) (*entity.Auction, error) {
    query := `
        SELECT id, lot_id, start_price, min_step, current_price,
               start_time, end_time, status, winner_id, winner_bid_id,
               created_at, updated_at
        FROM auctions
        WHERE lot_id = $1`

    auction := &entity.Auction{}
    err := r.db.QueryRowContext(ctx, query, lotID).Scan(
        &auction.ID,
        &auction.LotID,
        &auction.StartPrice,
        &auction.MinStep,
        &auction.CurrentPrice,
        &auction.StartTime,
        &auction.EndTime,
        &auction.Status,
        &auction.WinnerID,
        &auction.WinnerBidID,
        &auction.CreatedAt,
        &auction.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get auction by lot id", err)
    }

    return auction, nil
}

func (r *AuctionRepository) UpdateStatus(ctx context.Context, id int64, status entity.AuctionStatus) error {
    query := `
        UPDATE auctions
        SET status = $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

    result, err := r.db.ExecContext(ctx, query, status, id)
    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to update auction status", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to get rows affected", err)
    }

    if rows == 0 {
        return errors.New(errors.ErrorTypeNotFound, "auction not found", nil)
    }

    return nil
}

func (r *AuctionRepository) GetActiveAuctions(ctx context.Context) ([]*entity.Auction, error) {
    query := `
        SELECT id, lot_id, start_price, min_step, current_price,
               start_time, end_time, status, winner_id, winner_bid_id,
               created_at, updated_at
        FROM auctions
        WHERE status = $1`

    return r.queryAuctions(ctx, query, entity.AuctionStatusActive)
}

func (r *AuctionRepository) GetEndedAuctions(ctx context.Context) ([]*entity.Auction, error) {
    query := `
        SELECT id, lot_id, start_price, min_step, current_price,
               start_time, end_time, status, winner_id, winner_bid_id,
               created_at, updated_at
        FROM auctions
        WHERE status = $1 AND end_time <= CURRENT_TIMESTAMP`

    return r.queryAuctions(ctx, query, entity.AuctionStatusActive)
}

func (r *AuctionRepository) GetPendingAuctionsToStart(ctx context.Context) ([]*entity.Auction, error) {
    query := `
        SELECT id, lot_id, start_price, min_step, current_price,
               start_time, end_time, status, winner_id, winner_bid_id,
               created_at, updated_at
        FROM auctions
        WHERE status = $1 
        AND start_time <= CURRENT_TIMESTAMP`

    return r.queryAuctions(ctx, query, entity.AuctionStatusPending)
}

func (r *AuctionRepository) queryAuctions(ctx context.Context, query string, args ...interface{}) ([]*entity.Auction, error) {
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to query auctions", err)
    }
    defer rows.Close()

    var auctions []*entity.Auction
    for rows.Next() {
        auction := &entity.Auction{}
        err := rows.Scan(
            &auction.ID,
            &auction.LotID,
            &auction.StartPrice,
            &auction.MinStep,
            &auction.CurrentPrice,
            &auction.StartTime,
            &auction.EndTime,
            &auction.Status,
            &auction.WinnerID,
            &auction.WinnerBidID,
            &auction.CreatedAt,
            &auction.UpdatedAt,
        )
        if err != nil {
            return nil, errors.New(errors.ErrorTypeInternal, "failed to scan auction", err)
        }
        auctions = append(auctions, auction)
    }

    return auctions, nil
}
