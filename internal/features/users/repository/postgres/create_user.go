package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/artlink52/go-todoapp/internal/core/domain"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.users (full_name, phone_number)
		VALUES ($1, $2)
		RETURNING id, full_name, phone_number, version
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.Version,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
