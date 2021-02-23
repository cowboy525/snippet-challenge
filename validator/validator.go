package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/utils"
)

func ValidateNull(value interface{}, option string) *model.Error {
	defer func() {
		if r := recover(); r != nil {
			mlog.Error(fmt.Sprintln("an error occured inside of nil vaidator: ", r))
		}
	}()
	b, err := strconv.ParseBool(option)
	if err != nil {
		mlog.Error("invalid option detected for null validator")
		return nil
	}
	if b {
		return nil
	}
	if reflect.ValueOf(value).IsNil() {
		return model.NewError("model.validation.notnull", nil)
	}
	return nil
}

// TODO: only string at the moment
func ValidateBlank(value interface{}, option string) *model.Error {
	defer func() {
		if r := recover(); r != nil {
			mlog.Error(fmt.Sprintln("an error occured inside of blank validator: ", r))
		}
	}()
	b, err := strconv.ParseBool(option)
	if err != nil {
		mlog.Error("invalid option detected for blank validator")
		return nil
	}
	if b {
		return nil
	}
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if str, ok := value.(*string); ok && str != nil && len(*str) == 0 {
			return model.NewError("model.validation.blank", nil)
		}
	} else if str, ok := value.(string); ok && len(str) == 0 {
		return model.NewError("model.validation.blank", nil)
	}
	return nil
}

func ValidateEmail(value interface{}, option string) *model.Error {
	defer func() {
		if r := recover(); r != nil {
			mlog.Error(fmt.Sprintln("an error occured inside of email vaidator: ", r))
		}
	}()
	b, err := strconv.ParseBool(option)
	if err != nil {
		mlog.Error("invalid option detected for blank validator")
		return nil
	}
	if !b {
		return nil
	}
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if str, ok := value.(*string); ok && str != nil && !utils.IsValidEmail(*str) {
			return model.NewError("model.validation.invalid_email", nil)
		}
	} else if str, ok := value.(string); ok && !utils.IsValidEmail(str) {
		return model.NewError("model.validation.invalid_email", nil)
	}
	return nil
}

func Validate(item interface{}, data interface{}) *model.AppError {
	if utils.IsArrayInterface(item) {
		errors := ValidateArrayItem(item, data)
		if len(errors) > 0 {
			return model.ValidationErrorWithManyDetails("Validate", errors)
		}
	} else {
		errors := ValidateItem(item, data.(map[string]interface{}))
		if len(errors) > 0 {
			return model.ValidationErrorWithDetails("Validate", []map[string]interface{}{errors})
		}
	}
	return nil
}

func ValidateArrayItem(item interface{}, data interface{}) []map[string]interface{} {
	errors := []map[string]interface{}{}

	val := reflect.ValueOf(item)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	dataVal := reflect.ValueOf(data)
	for dataVal.Kind() == reflect.Ptr {
		dataVal = dataVal.Elem()
	}

	errorExists := false
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			valI := val.Index(i)
			for valI.Kind() == reflect.Ptr {
				valI = valI.Elem()
			}
			if valI.Kind() != reflect.Struct {
				return nil
			}
			err := ValidateItem(valI.Interface(), dataVal.Index(i).Interface().(map[string]interface{}))
			if len(err) > 0 {
				errorExists = true
			}
			errors = append(errors, err)
		}
	}
	if !errorExists {
		return nil
	}
	return errors
}

func ValidateItem(item interface{}, data map[string]interface{}) map[string]interface{} {
	errors := map[string]interface{}{}

	fields := structs.Fields(item)
	for _, field := range fields {
		fieldValue := field.Value()
		fieldKey := field.Tag("json")

		var isRequiredField bool = false
		tags := strings.Split(field.Tag("validate"), ";")
		for _, tag := range tags {
			v := strings.Split(tag, ":")
			if strings.TrimSpace(strings.ToLower(v[0])) == "required" {
				isRequiredField = true
				if len(v) > 1 {
					if b, err := strconv.ParseBool(v[1]); err == nil {
						isRequiredField = b
					}
				}
				break
			}
		}

		fieldData, isExisting := data[fieldKey]
		if isRequiredField && !isExisting {
			errors[fieldKey] = []*model.Error{model.NewError("model.validation.field_required", nil)}
		} else if isExisting {
			// if the field is an array
			if field.Kind() == reflect.Array || field.Kind() == reflect.Slice {
				errorsArray := ValidateArrayItem(fieldValue, fieldData)
				if len(errorsArray) > 0 {
					errors[fieldKey] = errorsArray
				}
			} else if field.Kind() == reflect.Struct {
				structErrors := ValidateItem(fieldValue, fieldData.(map[string]interface{}))
				if len(structErrors) != 0 {
					errors[fieldKey] = structErrors
				}
			} else {
				var fieldErrors []*model.Error
				for _, tag := range tags {
					v := strings.Split(tag, ":")
					if len(v) < 2 {
						continue
					}
					tagKey := strings.TrimSpace(strings.ToLower(v[0]))
					tagValue := v[1]
					if tagKey == "null" {
						if err := ValidateNull(fieldValue, tagValue); err != nil {
							fieldErrors = []*model.Error{err}
							break
						}
					} else if tagKey == "blank" {
						if err := ValidateBlank(fieldValue, tagValue); err != nil {
							fieldErrors = []*model.Error{err}
							break
						}
					} else if tagKey == "email" {
						if err := ValidateEmail(fieldValue, tagValue); err != nil {
							fieldErrors = []*model.Error{err}
							break
						}
					}
				}
				if len(fieldErrors) > 0 {
					errors[fieldKey] = fieldErrors
				}
			}
		}
	}
	return errors
}
