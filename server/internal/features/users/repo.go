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
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id int64, user *User) error
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
		SELECT id, username, display_name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
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
		SELECT id, username, display_name, email, created_at, updated_at, deleted_at
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
			&u.Username,
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
		SET username = $1,
			display_name = $2,
			email = $3,
			updated_at = NOW()
		where id = $4 AND deleted_at IS NULL
	`

	commandTag, err := r.pool.Exec(ctx, query, u.Username, u.DisplayName, u.Email, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", u.ID)
	}

	return nil
}

func (r *userRepo) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			username, 
			display_name, 
			email, 
			created_at
		) 
		VALUES ($1, $2, $3, NOW()) 
		RETURNING id`

	err := r.pool.QueryRow(ctx, query, user.Username, user.DisplayName, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
