package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/folklinoff/hack_and_change/models"
	"github.com/folklinoff/hack_and_change/pkg/hashing"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

type Room struct {
	ID             int64 `gorm:"primaryKey,autoIncrement"`
	Name           string
	OwnerId        int64
	Private        bool
	HashedPassword string
	Users          []User `gorm:"many2many:user_rooms"`
}

type UserRoom struct {
	UserID      int64 `gorm:"primaryKey"`
	RoomID      int64 `gorm:"primaryKey"`
	CreatedTime time.Time
	DeletedTime time.Time
	Relation    string
}

func (u *RoomRepository) Save(ctx context.Context, Room models.Room) (models.Room, error) {
	ent, err := RoomModelToEntity(Room)
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::Save: couldn't save room: %w", err)
	}

	err = u.db.WithContext(ctx).Preload("Users").Save(&ent).Error
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::Save: couldn't save room: %w", err)
	}
	Room.ID = ent.ID
	return Room, nil
}

func (u *RoomRepository) GetByName(ctx context.Context, name string) (models.Room, error) {
	var ent Room
	err := u.db.WithContext(ctx).Where(&Room{Name: name}).Preload("Users").First(&ent).Error
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::GetByMail: couldn't get room by mail: %w", err)
	}
	room := RoomEntityToModel(ent)
	return room, nil
}

func (u *RoomRepository) GetByID(ctx context.Context, id int64) (models.Room, error) {
	var ent Room
	err := u.db.WithContext(ctx).Where(&Room{ID: id}).Preload("Users").First(&ent).Error
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::GetByID: couldn't get room by id: %w", err)
	}
	room := RoomEntityToModel(ent)
	return room, nil
}

func (u *RoomRepository) Get(ctx context.Context) ([]models.Room, error) {
	var ents []Room
	err := u.db.WithContext(ctx).Model(&Room{}).Preload("Users").Find(&ents).Error
	if err != nil {
		return []models.Room{}, fmt.Errorf("RoomRepository::Get: couldn't get all rooms: %w", err)
	}
	rooms := RoomEntitiesToModels(ents)
	return rooms, nil
}

func (u *RoomRepository) LogIntoRoom(ctx context.Context, id, userId int64, password string) (models.Room, error) {
	var ent Room
	err := u.db.WithContext(ctx).Where(&Room{ID: id}).Preload("Users").First(&ent).Error
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::LogIntoRoom: couldn't log into room: %w", err)
	}
	if ent.Private {
		err = hashing.CompareHashAndPassword(password, ent.HashedPassword)
		if err != nil {
			return models.Room{}, fmt.Errorf("RoomRepository::LogIntoRoom: couldn't log into room: %w", err)
		}
	}

	room := RoomEntityToModel(ent)
	var user User
	err = u.db.Model(&User{ID: userId}).First(&user).Error
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::LogIntoRoom: couldn't log into room: %w", err)
	}

	room.Users = append(room.Users, userEntityToModel(user))
	err = u.db.WithContext(ctx).Save(room).Error
	if err != nil {
		return models.Room{}, fmt.Errorf("RoomRepository::LogIntoRoom: couldn't log into room: %w", err)
	}
	return room, nil
}

func (u *RoomRepository) AddUser(ctx context.Context, roomId, userId int64, relation string) error {
	return u.db.
		WithContext(ctx).
		Where(&Room{ID: roomId}).
		Association("Users").
		Append(&UserRoom{UserID: userId, RoomID: roomId, Relation: relation, CreatedTime: time.Now()})
}

// func (u *RoomRepository) ChangeUserPermissions(ctx context.Context, roomId, userId int64) (models.Room, error) {
// 	return u.db.
// 		WithContext(ctx).
// 		Where(&Room{ID: roomId}).
// 		Association("Users").
// }

func RoomModelToEntity(r models.Room) (Room, error) {
	hashedPassword, err := hashing.GeneratePasswordHash(r.Password)
	log.Println("Hashed password: ", hashedPassword, r.Password)
	if err != nil {
		return Room{}, fmt.Errorf("RoomModelToEntity: couldn't convert Room model to entity: %w", err)
	}
	users, err := userModelsToEntities(r.Users)
	if err != nil {
		return Room{}, err
	}
	return Room{
		ID:             r.ID,
		Name:           r.Name,
		Private:        r.Private,
		HashedPassword: hashedPassword,
		OwnerId:        r.OwnerId,
		Users:          users,
	}, nil
}

func RoomEntityToModel(r Room) models.Room {
	return models.Room{
		ID:      r.ID,
		Name:    r.Name,
		Private: r.Private,
		OwnerId: r.OwnerId,
		Users:   userEntitiesToModels(r.Users),
	}
}

func RoomEntitiesToModels(r []Room) []models.Room {
	models := make([]models.Room, len(r))
	for i := range r {
		models = append(models, RoomEntityToModel(r[i]))
	}
	return models
}
