package repository

import (
	"fmt"
	commonTime "github.com/lactobasilusprotectus/go-template/pkg/common/time"
	"github.com/lactobasilusprotectus/go-template/pkg/domain"
	"github.com/lactobasilusprotectus/go-template/pkg/util/db"
)

type UserRepository struct {
	dbClient *db.DatabaseConnection
	time     commonTime.TimeInterface
}

func NewUserRepository(dbClient *db.DatabaseConnection, time commonTime.TimeInterface) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
		time:     time,
	}
}

func (u *UserRepository) InsertUser(user domain.User) (err error) {
	result := u.dbClient.Master.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (u *UserRepository) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User

	result := u.dbClient.Slave.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (u *UserRepository) FindUserByID(id int64) (domain.User, error) {
	var user domain.User

	result := u.dbClient.Slave.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}
