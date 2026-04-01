package validation

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
	var blockedDomains = map[string]bool{
		"blacklist.com": true,
		"edu.vn":        true,
		"abc.com":       true,
	}
	v.RegisterValidation("email_advanced", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		parts := strings.Split(email, "@")

		if len(parts) != 2 {
			return false
		}
		domain := utils.NormalizeString(parts[1])

		return !blockedDomains[domain]
	})

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_fl", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minVal, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return false
		}
		return fl.Field().Float() >= minVal
	})

	v.RegisterValidation("max_fl", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		maxVal, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return false
		}
		return fl.Field().Float() <= maxVal
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		filename := fl.Field().String()

		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}

		allowedExt := strings.Fields(allowedStr)

		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

		for _, allowed := range allowedExt {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}

		return false
	})

	var (
		lowercaseRegex = regexp.MustCompile(`[a-z]`)
		uppercaseRegex = regexp.MustCompile(`[A-Z]`)
		numberRegex    = regexp.MustCompile(`[0-9]`)
		specialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>?/~` + "`" + `]`)
	)
	v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		if !lowercaseRegex.MatchString(password) {
			return false
		}

		if !uppercaseRegex.MatchString(password) {
			return false
		}

		if !numberRegex.MatchString(password) {
			return false
		}

		if !specialRegex.MatchString(password) {
			return false
		}

		return true

	})
}
