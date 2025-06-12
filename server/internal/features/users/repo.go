package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *CreateInput) error
	Update(ctx context.Context, id int64, user *User) error

	FindUserByEmail(ctx context.Context, email string) (*UserRecord, error)
}

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	query := `
		SELECT id, display_name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.DisplayName,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, display_name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		err := rows.Scan(
			&u.ID,
			&u.DisplayName,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) Update(ctx context.Context, id int64, u *User) error {
	query := `
		UPDATE users
		SET display_name = $2,
			email = $3,
			updated_at = NOW()
		where id = $4 AND deleted_at IS NULL
	`

	commandTag, err := r.pool.Exec(ctx, query, u.DisplayName, u.Email, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", u.ID)
	}

	return nil
}

type CreateInput struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (r *userRepo) Create(ctx context.Context, input *CreateInput) error {
	query := `
		INSERT INTO users (
			display_name, 
			email, 
			password,
			created_at
		) 
		VALUES ($1, $2, $3, NOW()) 
		RETURNING id`

	err := r.pool.QueryRow(ctx, query, input.DisplayName, input.Email, input.Password).Scan(&input.ID)
	if err != nil {
		return err
	}

	return nil
}

type UserRecord struct {
	User
	HashedPassword string
}

func (r *userRepo) FindUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	u := &UserRecord{}
	query := `
		SELECT id, display_name, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1
	`

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.DisplayName,
		&u.Email,
		&u.HashedPassword,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}
