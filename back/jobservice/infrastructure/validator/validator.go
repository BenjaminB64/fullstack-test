package validator

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"net/url"
	"regexp"
)

var slackRegex = regexp.MustCompile(`^/services/T[0-9A-Z]{5,15}/B[0-9A-Z]{5,15}/[0-9A-Za-z]{20,30}$`)

type Validator struct {
	universalTranslator *ut.UniversalTranslator
}

func NewValidator() (*Validator, error) {
	enLocale := en.New()
	universalTranslator := ut.New(enLocale, enLocale)

	return &Validator{
		universalTranslator: universalTranslator,
	}, nil
}

func (v Validator) RegisterOn(validate *validator.Validate) error {
	err := validate.RegisterValidation("slack_webhook_url", slackWebhookValidation)
	if err != nil {
		return errors.Join(errors.New("failed to register slack_webhook_url validation"), err)
	}

	err = validate.RegisterValidation("enum", enumValidation)
	if err != nil {
		return errors.Join(errors.New("failed to register enum validation"), err)
	}

	err = en_translations.RegisterDefaultTranslations(validate, v.universalTranslator.GetFallback())
	if err != nil {
		return errors.Join(errors.New("failed to register default translations"), err)
	}

	return nil
}

func (v Validator) GetValidationErrorsMap(validationErrors validator.ValidationErrors) (map[string]FieldError, error) {
	fieldErrors := make(map[string]FieldError, len(validationErrors))

	translator := v.universalTranslator.GetFallback()
	for _, err := range validationErrors {
		fieldErrors[err.Field()] = FieldError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Error: err.Translate(translator),
		}
	}

	return fieldErrors, nil
}

func enumValidation(fl validator.FieldLevel) bool {
	if validatable, ok := fl.Field().Interface().(EnumValidator); ok {
		return validatable.IsValid()
	}
	return false
}

func slackWebhookValidation(fl validator.FieldLevel) bool {
	rawURL := fl.Field().String()

	if rawURL == "" {
		return true
	}

	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if parsedUrl.Scheme != "https" {
		return false
	}

	if parsedUrl.Host != "hooks.slack.com" {
		return false
	}

	if !slackRegex.MatchString(parsedUrl.Path) {
		return false
	}

	return true
}

type ValidationErrorsTranslations map[string]FieldError
type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Error string `json:"error"`
}

type EnumValidator interface {
	IsValid() bool
}
