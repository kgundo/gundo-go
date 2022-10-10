package repositories

import (
	"time"

	"github.com/kgundo/gundo-go/loggers"
	"gorm.io/gorm"
)

type CommonRepository struct {
	logger *loggers.StandardLogger
}

type APIList struct {
	Items     []interface{} `json:"items"`
	TotalItem int64         `json:"totalItem"`
}

func InitCommonRepository() CommonRepository {
	logger := loggers.Logger()
	return CommonRepository{logger: logger}
}

var defaultOmitColumns = []string{"createdAt"}

func (cr *CommonRepository) UpdateData(db *gorm.DB, model string, data interface{}, customOmitColumns []string) error {
	omitColumns := defaultOmitColumns
	if len(customOmitColumns) > 0 {
		omitColumns = append(omitColumns, customOmitColumns...)
	}
	err := db.Table(model).Omit(omitColumns...).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func ConvertTimeToUTC(inputTime time.Time) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return inputTime.In(loc)
}
