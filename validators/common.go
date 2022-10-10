package validators

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"net/http"
)

type CommonValidator struct{}

const allowEmpty = "allowEmpty"

type QueryGetAllSchema struct {
	Page  string `json:"page" form:"page" binding:"omitempty,numeric,integer,minVal=1,maxVal=10000"`
	Limit string `json:"limit" form:"limit" binding:"omitempty,numeric,integer,minVal=1,maxVal=1000"`
}

type ParamIDSchema struct {
	ID string `json:"id" uri:"id" binding:"required,numeric,integer,minVal=1,maxVal=4294967295"`
}

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone", RegexPhone)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("integer", Integer)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("minVal", MinVal)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("maxVal", MaxVal)
	}
}

var RegexPhone validator.Func = func(fl validator.FieldLevel) bool {
	phoneNumber, ok := fl.Field().Interface().(string)
	if ok {
		if checkAllowEmptyString(fl.Param()) && phoneNumber == "" {
			return true
		}
		matched, _ := regexp.Match(`^[0-9]{10,11}$`, []byte(phoneNumber))
		return matched || phoneNumber == ""
	}
	return true
}

var Integer validator.Func = func(fl validator.FieldLevel) bool {
	qP, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, err := strconv.Atoi(qP)
	return err == nil
}

var MinVal validator.Func = func(fl validator.FieldLevel) bool {
	qP, _ := fl.Field().Interface().(string)
	value, err := strconv.ParseFloat(qP, 64)
	if err != nil {
		return false
	}
	md, _ := strconv.ParseFloat(fl.Param(), 64)
	return value >= md
}

var MaxVal validator.Func = func(fl validator.FieldLevel) bool {
	qP, _ := fl.Field().Interface().(string)
	value, err := strconv.ParseFloat(qP, 64)
	if err != nil {
		return false
	}
	md, _ := strconv.ParseFloat(fl.Param(), 64)
	return value <= md
}

func commonValidationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", lowerFirstCharacter(e.Field()))
	case "integer":
		return fmt.Sprintf("%s must be an integer", lowerFirstCharacter(e.Field()))
	}
	return fmt.Sprintf("%s is not valid", lowerFirstCharacter(e.Field()))
}

func checkAllowEmptyString(rule string) bool {
	return rule == allowEmpty
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorDescription struct {
	Key       string `json:"key"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

func HandleSuccess(data interface{}) Response {
	return Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}
}

func HandleError(code int, err interface{}) Response {
	var message = ""
	if code == 444 {
		message = "Not Found"
	} else {
		message = http.StatusText(code)
	}
	if err == nil {
		return Response{
			Code:    code,
			Message: message,
			Errors:  []ErrorDescription{},
		}
	}
	v := reflect.ValueOf(err)
	isSlice := v.Kind() == reflect.Slice
	if !isSlice {
		ret := make([]interface{}, 1)
		ret[0] = err
		return Response{
			Code:    code,
			Message: message,
			Errors:  ret,
		}
	}
	return Response{
		Code:    code,
		Message: message,
		Errors:  err,
	}
}

func lowerFirstCharacter(str string) string {
	if strings.ToLower(str) == "id" {
		return "id"
	}
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func upperFirstCharacter(str string) string {
	if strings.ToLower(str) == "id" {
		return "ID"
	}
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func formatBindJsonError(err *json.UnmarshalTypeError) []ErrorDescription {
	errs := []ErrorDescription{}
	errMessage := fmt.Sprintf("%s must be %s", err.Field, err.Type)
	errs = append(errs, ErrorDescription{Key: err.Field, Message: errMessage, ErrorCode: http.StatusUnprocessableEntity})
	return errs
}

func formatValidatorError(verr validator.ValidationErrors) []ErrorDescription {
	errs := []ErrorDescription{}
	for _, f := range verr {
		err := commonValidationErrorToText(f)
		errs = append(errs, ErrorDescription{Key: lowerFirstCharacter(f.Field()), Message: err, ErrorCode: http.StatusUnprocessableEntity})
	}
	return errs
}

func (cv *CommonValidator) HandleValidationError(c *gin.Context, err error) bool {
	if err != nil {
		if uErr, ok := err.(*json.UnmarshalTypeError); ok {
			c.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), HandleError(http.StatusUnprocessableEntity, formatBindJsonError(uErr)))
			return false
		} else if vErr, ok := err.(validator.ValidationErrors); ok {
			c.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), HandleError(http.StatusUnprocessableEntity, formatValidatorError(vErr)))
			return false
		} else {
			c.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), HandleError(http.StatusUnprocessableEntity, nil))
		}
	}
	return true
}

func (cv *CommonValidator) HandleRedundantFieldError(c *gin.Context, schema interface{}) bool {
	m := make(map[string]json.RawMessage)
	e := c.ShouldBindBodyWith(&m, binding.JSON)
	if e != nil {
		c.AbortWithStatusJSON(int(http.StatusBadRequest), HandleError(http.StatusBadRequest, e))
		return false
	}
	t := reflect.TypeOf(schema)
	errs := []ErrorDescription{}
	for s := range m {
		f, valid := t.FieldByName(upperFirstCharacter(s))
		if valid {
			_, v := f.Tag.Lookup("input")
			valid = v
		}
		if !valid {
			errs = append(errs, ErrorDescription{Key: s, Message: fmt.Sprintf("\"%s\" is not allowed", s), ErrorCode: http.StatusUnprocessableEntity})
			c.AbortWithStatusJSON(int(http.StatusUnprocessableEntity), HandleError(http.StatusUnprocessableEntity, errs))
			return false
		}
	}
	return true
}
