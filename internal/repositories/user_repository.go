package repositories

import (
	"context"
	"go-native-webserver/internal/dal"
	"go-native-webserver/internal/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	// db *sql.DB // Assume there's a database connection here
}

func NewUserRepository() UserRepository {
	return &userRepository{
		// Initialize db connection here
	}
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	// Implement the logic to find a user by ID from the database
	return nil, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := dal.DB.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// Implement the logic to create a new user in the database
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Implement the logic to update an existing user in the database
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	// Implement the logic to delete a user by ID from the database
	return nil
}
