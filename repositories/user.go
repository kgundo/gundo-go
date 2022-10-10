package repositories

import (
	"reflect"
	"time"

	"github.com/kgundo/gundo-go/loggers"
	"github.com/kgundo/gundo-go/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	CommonRepository
	logger *loggers.StandardLogger
}

type APIUserDetail struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy *uint      `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy *uint      `gorm:"column:updatedBy" json:"updatedBy"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt"`
	DeletedBy *uint      `gorm:"column:deletedBy" json:"deletedBy"`
}

func InitUserRepository() *UserRepository {
	commonRepo := InitCommonRepository()
	logger := loggers.Logger()
	return &UserRepository{CommonRepository: commonRepo, logger: logger}
}

func (r *UserRepository) CreateUser(db *gorm.DB, user interface{}) (interface{}, error) {
	err := db.Table((models.UserTableName)).Create(user).Error
	id := reflect.ValueOf(user).Elem().Field(0).Interface()
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (r *UserRepository) GetUserByID(db *gorm.DB, id int) (APIUserDetail, error) {
	var apiUserDetail APIUserDetail
	err := db.Model(&models.User{}).Select("users.*").Where("users.id = ?", id).Find(&apiUserDetail).Error
	if err != nil {
		return APIUserDetail{}, err
	}
	return apiUserDetail, nil
}

func (r *UserRepository) GetAllUsers(db *gorm.DB, limit int, offset int) (interface{}, error) {
	var list []APIUserDetail
	var totalItems int64
	var result APIList
	err := db.Model(&models.User{}).Select("*").Count(&totalItems).Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, err
	}
	tempList := []interface{}{}
	for _, item := range list {
		tempList = append(tempList, item)
	}
	result.Items = tempList
	result.TotalItem = totalItems
	return &result, nil
}

func (r *UserRepository) DeleteUser(db *gorm.DB, id int) error {
	now := ConvertTimeToUTC(time.Now())
	user := struct {
		ID        int        `gorm:"column:id"`
		DeletedAt *time.Time `gorm:"column:deletedAt"`
	}{
		ID:        id,
		DeletedAt: &now,
	}
	err := r.UpdateData(db, models.UserTableName, user, []string{})
	return err
}

func (r *UserRepository) UpdateUser(db *gorm.DB, user interface{}, id int) (*APIUserDetail, error) {
	err := db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return nil, err
	}
	result, err := r.GetUserByID(db, id)
	return &result, err
}
