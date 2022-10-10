package models

import "time"

const UserTableName string = "users"

type User struct {
	ID        uint       `gorm:"type:uint;size:32;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:varchar(255)" json:"name"`
	Email     string     `gorm:"type:varchar(255)" json:"email"`
	Password  string     `gorm:"type:varchar(63)" json:"password"`
	CreatedAt *time.Time `gorm:"type:datetime;column:createdAt" json:"createdAt"`
	CreatedBy *uint      `gorm:"type:uint;size:32;column:createdBy" json:"createdBy"`
	UpdatedAt *time.Time `gorm:"type:datetime;column:updatedAt" json:"updatedAt"`
	UpdatedBy *uint      `gorm:"type:uint;size:32;column:updatedBy" json:"updatedBy"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deletedAt" json:"deletedAt"`
	DeletedBy *uint      `gorm:"type:uint;size:32;column:deletedBy" json:"deletedBy"`
}
