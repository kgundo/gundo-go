package tables

import (
	"github.com/kgundo/gundo-go/models"
	"gorm.io/gorm"
)

func MigrateUser(dbInstance *gorm.DB) {
	_ = dbInstance.AutoMigrate(&models.User{})
}
