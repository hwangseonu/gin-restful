package gin_restful

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"sync"
)

var validate = new(defaultValidator)

func init() {
	validate.lazyinit()
	err := validate.validate.RegisterValidation("notblank", NotBlankValidate)
	if err != nil {
		panic(err)
	}
	binding.Validator = validate
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func mustString(i interface{}) (string, error) {
	if s, ok := i.(string); ok {
		return s, nil
	} else {
		return "", errors.New(" is must int")
	}
}

//구조체 필드에 `binding:"notblank"` 과 같이 사용합니다.
func NotBlankValidate(fl validator.FieldLevel) bool {
	if str, err := mustString(fl.Field().Interface()); err != nil {
		return false
	} else if regexp.MustCompile(`^\s*$`).MatchString(str) {
		return false
	} else {
		return true
	}
}
