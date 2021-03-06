package user_repository

import (
	"errors"

	"github.com/Valeriy-Totubalin/test_project/db/orm"
	"github.com/Valeriy-Totubalin/test_project/internal/app/interfaces/repository_interfaces"
	"github.com/Valeriy-Totubalin/test_project/internal/domain"
)

type UserRepository struct {
	Gorm repository_interfaces.GetterGormDB
}

func NewUserRepository(gorm repository_interfaces.GetterGormDB) *UserRepository {
	return &UserRepository{
		Gorm: gorm,
	}
}

func (repo *UserRepository) Create(user *domain.User) error {
	db, err := repo.Gorm.GetDB()
	if nil != err {
		return err
	}

	err = db.Create(&orm.User{
		Login:    user.Login,
		Password: user.Password,
	}).Error

	if nil != err {
		return err
	}

	return nil
}

func (repo *UserRepository) GetByLogin(login string) (*domain.User, error) {
	db, err := repo.Gorm.GetDB()
	if nil != err {
		return nil, err
	}

	user := orm.User{}
	err = db.Where("login = ?", login).Find(&user).Error
	if nil != err {
		return nil, err
	}

	if 0 == user.Id {
		return &domain.User{}, errors.New("user does not exist")
	}
	return &domain.User{
		Id:       user.Id,
		Login:    user.Login,
		Password: user.Password,
	}, nil
}

func (repo *UserRepository) GetById(userId int) (*domain.User, error) {
	db, err := repo.Gorm.GetDB()
	if nil != err {
		return nil, err
	}

	user := orm.User{}
	err = db.First(&user, userId).Error
	if nil != err {
		return nil, err
	}

	if 0 == user.Id {
		return nil, errors.New("user does not exist")
	}

	return &domain.User{
		Id:       user.Id,
		Login:    user.Login,
		Password: user.Password,
	}, nil
}

func (repo *UserRepository) IsExists(user *domain.User) bool {
	domainUser, _ := repo.GetByLogin(user.Login)
	if 0 == domainUser.Id {
		return false
	}

	return true
}
