package postgres

import (
    "context"
    "database/sql"
    "time"
    
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/errors"
)

type LotRepository struct {
    db *sql.DB
}

func NewLotRepository(db *sql.DB) *LotRepository {
    return &LotRepository{
        db: db,
    }
}

func (r *LotRepository) Create(ctx context.Context, lot *entity.Lot) error {
    query := `
        INSERT INTO lots (title, description, start_price, creator_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at`

    now := time.Now()
    lot.CreatedAt = now
    lot.UpdatedAt = now

    err := r.db.QueryRowContext(
        ctx,
        query,
        lot.Title,
        lot.Description,
        lot.StartPrice,
        lot.CreatorID,
        lot.CreatedAt,
        lot.UpdatedAt,
    ).Scan(&lot.ID, &lot.CreatedAt, &lot.UpdatedAt)

    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to create lot", err)
    }

    return nil
}

func (r *LotRepository) GetByID(ctx context.Context, id int64) (*entity.Lot, error) {
    query := `
        SELECT id, title, description, start_price, creator_id, created_at, updated_at
        FROM lots
        WHERE id = $1`

    lot := &entity.Lot{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &lot.ID,
        &lot.Title,
        &lot.Description,
        &lot.StartPrice,
        &lot.CreatorID,
        &lot.CreatedAt,
        &lot.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New(errors.ErrorTypeNotFound, "lot not found", nil)
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get lot", err)
    }

    return lot, nil
}

func (r *LotRepository) Update(ctx context.Context, id int64, lot *entity.Lot) (*entity.Lot, error) {
    query := `
        UPDATE lots 
        SET title = $1, description = $2, start_price = $3, updated_at = $4
        WHERE id = $5
        RETURNING id, title, description, start_price, creator_id, created_at, updated_at`

    now := time.Now()
    updatedLot := &entity.Lot{}
    err := r.db.QueryRowContext(
        ctx,
        query,
        lot.Title,
        lot.Description,
        lot.StartPrice,
        now,
        id,
    ).Scan(
        &updatedLot.ID,
        &updatedLot.Title,
        &updatedLot.Description,
        &updatedLot.StartPrice,
        &updatedLot.CreatorID,
        &updatedLot.CreatedAt,
        &updatedLot.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.New(errors.ErrorTypeNotFound, "lot not found", nil)
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to update lot", err)
    }

    return updatedLot, nil
}

func (r *LotRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM lots WHERE id = $1`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to delete lot", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to get rows affected", err)
    }

    if rowsAffected == 0 {
        return errors.New(errors.ErrorTypeNotFound, "lot not found", nil)
    }

    return nil
}

func (r *LotRepository) List(ctx context.Context, offset, limit int) ([]*entity.Lot, int64, error) {
    query := `
        SELECT id, title, description, start_price, creator_id, created_at, updated_at
        FROM lots
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to list lots", err)
    }
    defer rows.Close()

    var lots []*entity.Lot
    for rows.Next() {
        lot := &entity.Lot{}
        err := rows.Scan(
            &lot.ID,
            &lot.Title,
            &lot.Description,
            &lot.StartPrice,
            &lot.CreatorID,
            &lot.CreatedAt,
            &lot.UpdatedAt,
        )
        if err != nil {
            return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to scan lot", err)
        }
        lots = append(lots, lot)
    }

    var total int64
    err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM lots").Scan(&total)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to count lots", err)
    }

    return lots, total, nil
}
