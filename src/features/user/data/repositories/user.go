package repositories

import (
	"context"
	"main/src/core/db"
	"main/src/features/user/domain/entities"
	"main/src/features/user/domain/repositories"
)

type UsersRepository struct {
	db *db.Db
}

func (r *UsersRepository) GetUserById(ctx context.Context, id int64) (*entities.User, error) {
	row := r.db.GetClient().QueryRowContext(ctx, "SELECT id, first_name, last_name, user_name, "+
		"email, password_hash, created_at, updated_at, deleted_at FROM users "+
		"WHERE id = ? LIMIT 1", id)

	user := &entities.User{}

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUsersRepository(
	db *db.Db,
) repositories.UserRepository {
	return &UsersRepository{
		db: db,
	}
}
