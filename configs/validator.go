package configs

import (
	"fmt"
	"github.com/go-playground/locales/en"
	universalTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	engTranslation "github.com/go-playground/validator/v10/translations/en"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func isValidIndonesianPhoneNumber(phone string) bool {
	phoneNumberRegex := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,9}$`)
	return phoneNumberRegex.MatchString(phone)
}

func InitializeValidator() (*validator.Validate, universalTranslator.Translator) {
	var validatorInstance = validator.New()
	// Tambahkan validasi kustom untuk 'datetime'
	validatorInstance.RegisterValidation("date", func(fieldLevel validator.FieldLevel) bool {
		if fieldLevel.Field().String() == "" {
			return true
		}
		format := "2006-01-02"
		_, err := time.Parse(format, fieldLevel.Field().String())
		return err == nil
	})
	validatorInstance.RegisterValidation("datetime", func(fieldLevel validator.FieldLevel) bool {
		if fieldLevel.Field().String() == "" {
			return true
		}
		format := "2006-01-02 15:04"
		_, err := time.Parse(format, fieldLevel.Field().String())
		return err == nil
	})
	validatorInstance.RegisterValidation("maxSize", maxFileSizeValidation)
	validatorInstance.RegisterValidation("extensionFile", validateFileExtensionValidation)
	validatorInstance.RegisterValidation("obligatoryFile", requiredFileValidationValidation)
	validatorInstance.RegisterValidation("phoneNumber", phoneNumberValidation)
	validatorInstance.RegisterValidation("conditionalRequired", conditionalRequired)
	validatorInstance.RegisterValidation("weakPassword", weakPassword)
	validatorInstance.RegisterValidation("decimalValue", decimalValue)

	englishLang := en.New()
	universalTranslatorInstance := universalTranslator.New(englishLang, englishLang)
	engTranslator, _ := universalTranslatorInstance.GetTranslator("en")
	_ = engTranslation.RegisterDefaultTranslations(validatorInstance, engTranslator)

	// Register custom error message
	validatorInstance.RegisterTranslation("maxSize", engTranslator, func(ut universalTranslator.Translator) error {
		return ut.Add("maxSize", "File size must be at most {0} MB", true)
	}, func(ut universalTranslator.Translator, fe validator.FieldError) string {
		return fmt.Sprintf("File size must be at most %s MB", fe.Param())
	})
	// Register translations for "obligatoryFile"
	validatorInstance.RegisterTranslation("obligatoryFile", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("obligatoryFile", "File is required", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return "File is required"
		})

	// Register translations for "fileExtension"
	validatorInstance.RegisterTranslation("extensionFile", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("extensionFile", "File extension must be one of {0}", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return fmt.Sprintf("File extension must be one of %s", fe.Param())
		})

	validatorInstance.RegisterTranslation("phoneNumber", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("phoneNumber", "Format phone number not valid", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return "Format phone number not valid"
		})

	validatorInstance.RegisterTranslation("conditionalRequired", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("conditionalRequired", "One of the field must be filled", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return fmt.Sprintf("You must fill %s if %s blank", fe.Field(), fe.Param())
		})

	validatorInstance.RegisterTranslation("conditionalRequired", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("conditionalRequired", "One of the field must be filled", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return fmt.Sprintf("You must fill %s if %s blank", fe.Field(), fe.Param())
		})

	// Register Translation
	validatorInstance.RegisterTranslation("weakPassword", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("weak_password", "Password is weak", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return "Password is weak, please add upper case, lower case, and digit"
		},
	)

	validatorInstance.RegisterTranslation("decimalValue", engTranslator,
		func(ut universalTranslator.Translator) error {
			return ut.Add("decimal_value", "Value must be decimal", true)
		},
		func(ut universalTranslator.Translator, fe validator.FieldError) string {
			return fmt.Sprintf("%s must be decimal", fe.Field())
		},
	)
	return validatorInstance, engTranslator
}

func phoneNumberValidation(fieldLevel validator.FieldLevel) bool {
	stringValue, isValid := fieldLevel.Field().Interface().(string)
	if !isValid {
		return false
	}
	return isValidIndonesianPhoneNumber(stringValue)
}

func requiredFileValidationValidation(fieldLevel validator.FieldLevel) bool {
	_, ok := fieldLevel.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}
	return true
}

func maxFileSizeValidation(fieldLevel validator.FieldLevel) bool {
	file, ok := fieldLevel.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	// Dapatkan parameter maxSize
	maxSize, err := strconv.ParseInt(fieldLevel.Param(), 10, 64)
	if err != nil {
		return false
	}

	// Validasi ukuran file
	return file.Size <= maxSize*1024*1024
}

// Validasi ekstensi file
func validateFileExtensionValidation(fieldLevel validator.FieldLevel) bool {
	// Ambil file dari field
	file, ok := fieldLevel.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	// Ambil parameter ekstensi yang diperbolehkan
	allowedExtension := fieldLevel.Param() // Misalnya "png,jpg,jpeg"

	// Pisahkan parameter ekstensi menjadi slice
	extSlice := strings.Split(allowedExtension, " ")

	// Ambil ekstensi file
	filename := strings.ToLower(file.Filename)
	for _, ext := range extSlice {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func conditionalRequired(fieldLevel validator.FieldLevel) bool {
	fieldName := fieldLevel.Param()                    // Nama field yang terkait, misalnya "PhoneNumber"
	structValue := fieldLevel.Parent()                 // Struct yang sedang divalidasi
	relatedField := structValue.FieldByName(fieldName) // Ambil field berdasarkan nama

	// Pastikan field terkait ada dalam struct
	if !relatedField.IsValid() {
		return false
	}

	// Ambil nilai dari field yang sedang divalidasi dan field terkait
	currentValue := fieldLevel.Field().String()
	relatedValue := relatedField.String()

	// Jika field terkait kosong, maka field ini wajib diisi
	return relatedValue != "" || currentValue != ""
}

func decimalValue(fieldLevel validator.FieldLevel) bool {
	switch value := fieldLevel.Field().Interface().(type) {
	case float64:
		return true
	case float32:
		return true
	case int, int32, int64:
		return true
	case string:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	default:
		return false
	}
}

// WeakPasswordValidator untuk mengecek kelemahan password
func weakPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	errors := []string{}

	// Minimal panjang 8 karakter
	if len(password) < 8 {
		errors = append(errors, "must be at least 8 characters long")
	}

	// Harus ada huruf besar
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		errors = append(errors, "must contain at least one uppercase letter")
	}

	// Harus ada huruf kecil
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		errors = append(errors, "must contain at least one lowercase letter")
	}

	// Harus ada angka
	if !regexp.MustCompile(`\d`).MatchString(password) {
		errors = append(errors, "must contain at least one digit")
	}

	// Harus ada karakter spesial
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		errors = append(errors, "must contain at least one special character")
	}

	// Jika ada error, simpan dalam `fl.Param()` untuk translasi
	if len(errors) > 0 {
		return false
	}

	return true
}
