package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kgundo/gundo-go/loggers"
	"github.com/kgundo/gundo-go/services"
	"github.com/kgundo/gundo-go/validators"
)

type UserController struct {
	service       *services.UserService
	commonService services.CommonService
	logger        *loggers.StandardLogger
	validator     *validators.UserValidator
}

func InitUserController() *UserController {
	userService := services.InitUserService()
	commonService := services.InitCommonService()
	logger := loggers.Logger()
	validator := &validators.UserValidator{}
	return &UserController{service: userService, commonService: commonService, validator: validator, logger: logger}
}

func (ctr *UserController) CreateUser(c *gin.Context) {
	bodyFormData, isValidated := ctr.validator.CreateUserValidator(c)
	if !isValidated {
		return
	}
	result, err := ctr.service.CreateUser(&bodyFormData)
	if err != nil {
		if e, ok := err.(validators.ErrorDescription); ok {
			c.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), validators.HandleError(http.StatusUnprocessableEntity, e))
			return
		}
		c.AbortWithStatusJSON(int(http.StatusInternalServerError), validators.HandleError(http.StatusInternalServerError, err))
		return
	}
	c.JSON(int(http.StatusOK), validators.HandleSuccess(result))
}

func (ctr *UserController) GetAllUsers(c *gin.Context) {
	query, isValidated := ctr.validator.GetAllUsersValidator(c)
	if !isValidated {
		return
	}
	limit, _ := strconv.Atoi(query.Limit)
	page, _ := strconv.Atoi(query.Page)
	list, err := ctr.service.GetAllUsers(limit, page)
	if err != nil {
		c.AbortWithStatusJSON(int(http.StatusInternalServerError), validators.HandleError(http.StatusInternalServerError, err))
		return
	}
	c.JSON(int(http.StatusOK), validators.HandleSuccess(list))
}

func (ctr *UserController) GetUserByID(c *gin.Context) {
	params, isValidated := ctr.validator.GetUserByIDValidator(c)
	if !isValidated {
		return
	}
	id, _ := strconv.Atoi(params.ID)
	result, err := ctr.service.GetUserByID(id)
	if err != nil {
		if err.Error() == "Not Found" {
			c.AbortWithStatusJSON(444, validators.HandleError(444, nil))
			return
		}
		c.AbortWithStatusJSON(int(http.StatusInternalServerError), validators.HandleError(http.StatusInternalServerError, err))
		return
	}
	c.JSON(int(http.StatusOK), validators.HandleSuccess(result))
}

func (ctr *UserController) DeleteUser(c *gin.Context) {
	params, isValidated := ctr.validator.DeleteUserValidator(c)
	if !isValidated {
		return
	}
	id, _ := strconv.Atoi(params.ID)
	err := ctr.service.DeleteUser(id)
	if err != nil {
		if err.Error() == "Not Found" {
			c.AbortWithStatusJSON(444, validators.HandleError(444, nil))
			return
		}
		if err.Error() == http.StatusText(http.StatusConflict) {
			c.AbortWithStatusJSON(444, validators.HandleError(http.StatusConflict, nil))
			return
		}
		c.AbortWithStatusJSON(int(http.StatusInternalServerError), validators.HandleError(http.StatusInternalServerError, err))
		return
	}
	c.JSON(int(http.StatusOK), validators.HandleSuccess(&struct {
		ID int `json:"id"`
	}{id}))
}

func (ctr *UserController) UpdateUser(c *gin.Context) {
	params, bodyFormData, isValidated := ctr.validator.UpdateUserValidator(c)
	if !isValidated {
		return
	}
	id, _ := strconv.Atoi(params.ID)
	result, err := ctr.service.UpdateUser(id, bodyFormData)
	if err != nil {
		if e, ok := err.(validators.ErrorDescription); ok {
			c.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), validators.HandleError(http.StatusUnprocessableEntity, e))
			return
		}
		if e, ok := err.(error); ok && e.Error() == "Not Found" {
			c.AbortWithStatusJSON(444, validators.HandleError(444, nil))
			return
		}
		c.AbortWithStatusJSON(int(http.StatusInternalServerError), validators.HandleError(http.StatusInternalServerError, err))
		return
	}
	c.JSON(int(http.StatusOK), validators.HandleSuccess(result))
}
