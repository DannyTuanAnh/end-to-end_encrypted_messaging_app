package validation

import (
	"fmt"
	"strings"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return fmt.Errorf("failed to register validator")
	}

	RegisterCustomValidations(v)

	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationErr {
			root := strings.Split(e.Namespace(), ".")[0]

			raw := strings.TrimPrefix(e.Namespace(), root+".")

			paths := strings.Split(raw, ".")

			for i, path := range paths {
				if strings.Contains(path, "[") {
					idx := strings.Index(path, "[")
					base := utils.CamelToSnake(path[:idx])

					index := path[idx:]
					paths[i] = base + index
				} else {
					paths[i] = utils.CamelToSnake(path)
				}
			}

			fieldPath := strings.Join(paths, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s phải nhỏ hơn %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s phải nhỏ hơn hoặc bằng %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn hoặc bằng %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s phải là định dạng UUID hợp lệ", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s chỉ được chứa các chữ cái thường, số, dấu gạch ngang hoặc dấu chấm", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s phải nhiều hơn %s", fieldPath, e.Param()) + " ký tự"
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s phải ít hơn %s", fieldPath, e.Param()) + " ký tự"
			case "min_fl":
				errors[fieldPath] = fmt.Sprintf("%s phải lớn hơn hoặc bằng giá trị %s", fieldPath, e.Param())
			case "max_fl":
				errors[fieldPath] = fmt.Sprintf("%s phải nhỏ hơn hoặc bằng giá trị %s", fieldPath, e.Param())
			case "oneof":
				allowedValue := strings.ReplaceAll(e.Param(), " ", ", ")
				errors[fieldPath] = fmt.Sprintf("%s phải là một trong các giá trị sau: %s", fieldPath, allowedValue)
			case "file_ext":
				allowedValue := strings.ReplaceAll(e.Param(), " ", ", ")
				errors[fieldPath] = fmt.Sprintf("%s chỉ cho phép những file sau: %s", fieldPath, allowedValue)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s là bắt buộc", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s chỉ được chứa các chữ cái, số và khoảng trắng", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s phải đúng định dạng là email", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s phải đúng định dạng YYYY-MM-DD", fieldPath)
			case "email_advanced":
				errors[fieldPath] = fmt.Sprintf("%s này trong danh sách các tên miền bị cấm", fieldPath)
			case "strong_password":
				special := "!@#$%^&*()_+-=[]{},;.:|<>?/~'`" + `"`
				errors[fieldPath] = fmt.Sprintf("%s phải bao gồm chữ hoa, chữ thường, số và ký tự đặc biệt (%s)", fieldPath, special)
			default:
				errors[fieldPath] = fmt.Sprintf("%s không hợp lệ", fieldPath)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{
		"error":  "Yêu cầu không hợp lệ",
		"detail": err.Error(),
	}
}
