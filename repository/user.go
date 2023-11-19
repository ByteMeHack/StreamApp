package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/folklinoff/hack_and_change/pkg/hashing"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type User struct {
	ID             int64 `gorm:"primaryKey,autoIncrement"`
	Name           string
	HashedPassword string
	Email          string
	Rooms          []Room `gorm:"many2many:user_rooms"`
}

func (u *UserRepository) Save(ctx context.Context, user models.User) (models.User, error) {
	ent, err := userModelToEntity(user)
	if err != nil {
		return models.User{}, fmt.Errorf("UserRepository::Save: couldn't save user: %w", err)
	}

	err = u.db.WithContext(ctx).Save(&ent).Error
	if err != nil {
		return models.User{}, fmt.Errorf("UserRepository::Save: couldn't save user: %w", err)
	}
	user.ID = ent.ID
	return user, nil
}

func (u *UserRepository) GetByMail(ctx context.Context, mail string) (models.User, error) {
	var ent User
	err := u.db.WithContext(ctx).Where(&User{Email: mail}).First(&ent).Error
	if err != nil {
		return models.User{}, fmt.Errorf("UserRepository::GetByMail: couldn't get user by mail: %w", err)
	}
	user := userEntityToModel(ent)
	return user, nil
}

func (u *UserRepository) GetByID(ctx context.Context, id int64) (models.User, error) {
	var ent User
	err := u.db.WithContext(ctx).Where(&User{ID: id}).First(&ent).Error
	if err != nil {
		return models.User{}, fmt.Errorf("UserRepository::GetByID: couldn't get user id: %w", err)
	}
	user := userEntityToModel(ent)
	return user, nil
}

func (u *UserRepository) LoginUserByMail(ctx context.Context, mail string, password string) (models.User, error) {
	var ent User
	err := u.db.WithContext(ctx).Where(&User{Email: mail}).First(&ent).Error
	if err != nil {
		return models.User{}, fmt.Errorf("UserRepository::LoginUserByMail: couldn't get user by mail: %w", err)
	}
	err = hashing.CompareHashAndPassword(password, ent.HashedPassword)
	if err != nil {
		log.Println("Given password: ", password, "; User: ", mail)
		return models.User{}, fmt.Errorf("UserRepository::LoginUserByMail: couldn't login user: password is incorrect: %w", err)
	}
	user := userEntityToModel(ent)
	return user, nil
}

func userModelToEntity(u models.User) (User, error) {
	hashedPassword, err := hashing.GeneratePasswordHash(u.Password)
	log.Println("Hashed password: ", hashedPassword, u.Password)
	if err != nil {
		return User{}, fmt.Errorf("userModelToEntity: couldn't convert user model to entity: %w", err)
	}
	return User{
		ID:             u.ID,
		Name:           u.Name,
		HashedPassword: hashedPassword,
		Email:          u.Email,
	}, nil
}

func userModelsToEntities(u []models.User) ([]User, error) {
	ans := make([]User, len(u))
	for i := range u {
		var err error
		ans[i], err = userModelToEntity(u[i])
		if err != nil {
			return nil, err
		}
	}
	return ans, nil
}

func userEntityToModel(u User) models.User {
	return models.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func userEntitiesToModels(u []User) []models.User {
	ans := make([]models.User, len(u))
	for i := range u {
		ans[i] = userEntityToModel(u[i])
	}
	return ans
}
