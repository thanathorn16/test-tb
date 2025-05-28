package stores

import (
	"be/services/login"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ login.UserRepo = (*userDB)(nil)

type userDB struct {
	db *gorm.DB
}

type User struct {
	ID       string `gorm:"column:id;primaryKey"`
	UserName string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func NewUserDB(db *gorm.DB) *userDB {
	return &userDB{db: db}
}

func (r *userDB) CheckUserExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userDB) FindByUsernameAndPassword(ctx context.Context, username, password string) (string, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return user.ID, nil
}

func (r *userDB) CreateUser(ctx context.Context, userID, username, password string) error {
	user := User{
		ID:       userID,
		UserName: username,
		Password: password,
	}
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userDB) GetUserProfile(ctx context.Context, userID string) (string, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return user.UserName, nil
}
