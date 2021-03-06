package service_interfaces

import "github.com/Valeriy-Totubalin/test_project/internal/domain"

type ItemService interface {
	Create(item *domain.Item) error
	Delete(item *domain.Item) error
	GetAll(userId int) ([]*domain.Item, error)
	GetTempLink(link *domain.Link) (string, error)
	CanConfirm(tempLink string, userId int) (bool, error)
	Confirm(tempLink string, userId int) error
	IsOwner(itemId int, userId int) (bool, error)
	IsDeleted(itemId int) (bool, error)
}
