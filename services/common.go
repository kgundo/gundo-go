package services

import (
	"github.com/kgundo/gundo-go/db"
	"github.com/kgundo/gundo-go/loggers"
	"gorm.io/gorm"
)

type CommonService struct {
	db     *gorm.DB
	logger *loggers.StandardLogger
}

func InitCommonService() CommonService {
	dbInstance := db.InitDB()
	logger := loggers.Logger()
	return CommonService{db: dbInstance, logger: logger}
}
