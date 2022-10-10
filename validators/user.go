package validators

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserValidator struct {
	CommonValidator
}

type BodyCreateUserSchema struct {
	ID        int        `gorm:"primaryKey;autoIncrement" form:"id" json:"id" binding:"omitempty"`
	Name      string     `json:"name" binding:"required,max=255" input:"1"`
	Email     string     `json:"email" binding:"required,max=255" input:"1"`
	Password  string     `json:"password" binding:"required,max=63" input:"1"`
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt,omitempty" binding:"omitempty"`
	CreatedBy *int       `gorm:"column:createdBy;default:null" json:"createdBy,omitempty" binding:"omitempty,min=1,max=4294967295"`
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt,omitempty" binding:"omitempty"`
	UpdatedBy *int       `gorm:"column:updatedBy;default:null" json:"updatedBy,omitempty" binding:"omitempty,min=1,max=4294967295"`
	DeletedAt *time.Time `gorm:"column:deletedAt" json:"deletedAt,omitempty" binding:"omitempty"`
	DeletedBy *int       `gorm:"column:deletedBy;default:null" json:"deletedBy,omitempty" binding:"omitempty,min=1,max=4294967295"`
}

type BodyUpdateUserSchema struct {
	ID        int        `json:"id" binding:"omitempty"`
	Name      string     `json:"name" binding:"required,max=255" input:"1"`
	Email     string     `json:"email" binding:"required,max=255" input:"1"`
	Password  string     `json:"password" binding:"required,max=63" input:"1"`
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt,omitempty" binding:"omitempty"`
	UpdatedBy *int       `gorm:"column:updatedBy;default:null" json:"updatedBy,omitempty" binding:"omitempty,min=1,max=4294967295"`
}

func (v *UserValidator) CreateUserValidator(c *gin.Context) (BodyCreateUserSchema, bool) {
	var body BodyCreateUserSchema
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	return body, v.HandleRedundantFieldError(c, body) && v.HandleValidationError(c, err)
}

func (v *UserValidator) GetAllUsersValidator(c *gin.Context) (QueryGetAllSchema, bool) {
	var query QueryGetAllSchema
	err := c.ShouldBindQuery(&query)
	return query, v.HandleValidationError(c, err)
}

func (v *UserValidator) GetUserByIDValidator(c *gin.Context) (ParamIDSchema, bool) {
	var param ParamIDSchema
	err := c.ShouldBindUri(&param)
	return param, v.HandleValidationError(c, err)
}

func (v *UserValidator) DeleteUserValidator(c *gin.Context) (ParamIDSchema, bool) {
	var param ParamIDSchema
	err := c.ShouldBindUri(&param)
	return param, v.HandleValidationError(c, err)
}

func (v *UserValidator) UpdateUserValidator(c *gin.Context) (ParamIDSchema, map[string]interface{}, bool) {
	var param ParamIDSchema
	errParam := c.ShouldBindUri(&param)
	body := make(map[string]interface{})
	var validator BodyUpdateUserSchema
	errBody := c.ShouldBindBodyWith(&validator, binding.JSON)
	_ = c.ShouldBindBodyWith(&body, binding.JSON)
	return param, body, v.HandleValidationError(c, errParam) && v.HandleValidationError(c, errBody)
}
