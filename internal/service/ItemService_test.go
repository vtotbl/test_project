package service

import (
	"testing"
	"time"

	"github.com/Valeriy-Totubalin/test_project/internal/domain"
	"github.com/Valeriy-Totubalin/test_project/pkg/link_manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type LinkManagerMock struct {
	mock.Mock
}

func (m *LinkManagerMock) NewLink(link *link_manager.Link, ttl time.Duration) (string, error) {
	m.Called(link, ttl)

	return "temp_link", nil
}

func (m *LinkManagerMock) Parse(tempLink string) (*link_manager.Link, error) {
	m.Called(tempLink)

	return &link_manager.Link{
		ItemId:    42,
		UserLogin: "test_login",
	}, nil
}

type ItemRepositoryMock struct {
	mock.Mock
}

func (m *ItemRepositoryMock) Create(item *domain.Item) error {
	m.Called(item)

	return nil
}

func (m *ItemRepositoryMock) DeleteById(itemId int) error {
	m.Called(itemId)

	return nil
}

func (m *ItemRepositoryMock) GetAll() ([]*domain.Item, error) {
	m.Called()

	return []*domain.Item{
		{
			Id: 42,
		},
		{
			Id: 23,
		},
		{
			Id: 97,
		},
	}, nil
}

func (m *ItemRepositoryMock) Transfer(itemId int, userId int) error {
	m.Called(itemId, userId)

	return nil
}

func TestNewItemService(t *testing.T) {
	repository := new(ItemRepositoryMock)
	userRepository := new(UserRepositoryMock)
	linkManager := new(LinkManagerMock)

	serviceExpected := &ItemService{
		ItemRepository: repository,
		LinkManager:    linkManager,
		UserRepository: userRepository,
	}

	serviceEqual := NewItemService(repository, linkManager, userRepository)

	assert.Equal(t, serviceExpected, serviceEqual)
}

func TestCreate(t *testing.T) {
	repository := new(ItemRepositoryMock)
	linkManager := new(LinkManagerMock)
	userRepository := new(UserRepositoryMock)

	service := NewItemService(repository, linkManager, userRepository)

	item := &domain.Item{
		Name: "item_name",
	}

	repository.On("Create", item).Return(nil).Once()

	err := service.Create(item)

	repository.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	repository := new(ItemRepositoryMock)
	linkManager := new(LinkManagerMock)
	userRepository := new(UserRepositoryMock)

	service := NewItemService(repository, linkManager, userRepository)

	item := &domain.Item{
		Id: 42,
	}

	repository.On("DeleteById", item.Id).Return(nil).Once()

	err := service.Delete(item)

	repository.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestGetAll(t *testing.T) {
	repository := new(ItemRepositoryMock)
	linkManager := new(LinkManagerMock)
	userRepository := new(UserRepositoryMock)

	service := NewItemService(repository, linkManager, userRepository)

	items := []*domain.Item{
		{
			Id: 42,
		},
		{
			Id: 23,
		},
		{
			Id: 97,
		},
	}

	repository.On("GetAll").Return(items).Once()

	itemsReturned, err := service.GetAll()

	repository.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, items, itemsReturned)
}

func TestGetTempLink(t *testing.T) {
	repository := new(ItemRepositoryMock)
	linkManager := new(LinkManagerMock)
	userRepository := new(UserRepositoryMock)

	service := NewItemService(repository, linkManager, userRepository)

	link := &domain.Link{
		ItemId:    42,
		UserLogin: "test_login",
	}
	libLink := &link_manager.Link{
		ItemId:    42,
		UserLogin: "test_login",
	}

	tempLink := "temp_link"

	linkManager.On("NewLink", libLink, 24*time.Hour).Return(tempLink, nil).Once()

	linkReturned, err := service.GetTempLink(link)

	linkManager.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, tempLink, linkReturned)
}

func TestCanConfirm(t *testing.T) {
	repository := new(ItemRepositoryMock)
	linkManager := new(LinkManagerMock)
	userRepository := new(UserRepositoryMock)

	service := NewItemService(repository, linkManager, userRepository)

	tempLink := "temp_link"
	userId := 42
	link := &link_manager.Link{
		ItemId:    42,
		UserLogin: "test_login",
	}
	user := &domain.User{
		Id:       42,
		Login:    "test_login",
		Password: "password_hash",
	}

	linkManager.On("Parse", tempLink).Return(link, nil).Once()
	userRepository.On("GetById", userId).Return(user, nil).Once()

	canConfirm, err := service.CanConfirm(tempLink, userId)

	linkManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	assert.Equal(t, userId, user.Id)
	assert.Nil(t, err)
	assert.True(t, canConfirm)
}

func TestConfirm(t *testing.T) {
	repository := new(ItemRepositoryMock)
	linkManager := new(LinkManagerMock)
	userRepository := new(UserRepositoryMock)

	service := NewItemService(repository, linkManager, userRepository)

	userId := 42
	tempLink := "temp_link"
	link := &link_manager.Link{
		ItemId:    42,
		UserLogin: "test_login",
	}

	linkManager.On("Parse", tempLink).Return(link, nil).Once()
	repository.On("Transfer", link.ItemId, userId).Return(nil).Once()

	err := service.Confirm(tempLink, userId)

	linkManager.AssertExpectations(t)
	repository.AssertExpectations(t)
	assert.Nil(t, err)
}
