package postgres

import (
    "context"
    "database/sql"
    "time"
    "auction-system/internal/domain/entity"
    "auction-system/internal/domain/errors"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{
        db: db,
    }
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	if err := ctx.Err(); err != nil {
        return err
    }

	query := `
        INSERT INTO users (username, email, balance, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

    user.CreatedAt = time.Now()
    
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Username,
        user.Email,
        user.Balance,
        user.CreatedAt,
    ).Scan(&user.ID)

    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to create user", err)
    }

    return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	if err := ctx.Err(); err != nil {
        return nil, err
    }

	query := `SELECT id, username, email, balance, created_at FROM users`
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get all users", err)
    }
    defer rows.Close()

    var users []*entity.User
    for rows.Next() {
        user := &entity.User{}
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Balance, &user.CreatedAt)
        if err != nil {
            return nil, errors.New(errors.ErrorTypeInternal, "failed to scan user", err)
        }
        users = append(users, user)
    }

    return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
        return err
    }

	query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

func (r *UserRepository) Update(ctx context.Context, id int64, user *entity.User) (*entity.User, error) {
	if err := ctx.Err(); err != nil {
        return nil, err
    }

	query := `
        UPDATE users 
        SET username = $1, email = $2, balance = $3 
        WHERE id = $4 
        RETURNING id, username, email, balance, created_at`
        
    updatedUser := &entity.User{}
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Username,
        user.Email,
        user.Balance,
        id,
    ).Scan(
        &updatedUser.ID,
        &updatedUser.Username,
        &updatedUser.Email,
        &updatedUser.Balance,
        &updatedUser.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.NewNotFoundError("user not found")
    }
    
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to update user", err)
    }

    return updatedUser, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	if err := ctx.Err(); err != nil {
        return nil, err
    }

	query := `
        SELECT id, username, email, balance, created_at
        FROM users
        WHERE id = $1`

    user := &entity.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Balance,
        &user.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, errors.NewNotFoundError("user not found")
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get user", err)
    }

    return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if err := ctx.Err(); err != nil {
        return nil, err
    }
	
	query := `
        SELECT id, username, email, balance, created_at
        FROM users
        WHERE email = $1`

    user := &entity.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Balance,
        &user.CreatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, errors.New(errors.ErrorTypeInternal, "failed to get user by email", err)
    }

    return user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id int64, balance float64) error {
	if err := ctx.Err(); err != nil {
        return err
    }

	query := `
        UPDATE users
        SET balance = $1
        WHERE id = $2`

    result, err := r.db.ExecContext(ctx, query, balance, id)
    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to update balance", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.New(errors.ErrorTypeInternal, "failed to get rows affected", err)
    }

    if rowsAffected == 0 {
        return errors.NewNotFoundError("user not found")
    }

    return nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	if err := ctx.Err(); err != nil {
        return nil, 0, err
    }

	query := `
        SELECT id, username, email, balance, created_at
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2`

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to list users", err)
    }
    defer rows.Close()

    var users []*entity.User
    for rows.Next() {
        user := &entity.User{}
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.Balance,
            &user.CreatedAt,
        )
        if err != nil {
            return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to scan user", err)
        }
        users = append(users, user)
    }

    var total int64
    err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
    if err != nil {
        return nil, 0, errors.New(errors.ErrorTypeInternal, "failed to count users", err)
    }

    return users, total, nil
}
