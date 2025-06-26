package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ptBR_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// A função init é executada automaticamente quando o pacote é importado pela primeira vez.
// Ela inicializa uma única instância do validador e do tradutor.
func init() {
	// Cria a instância do validador.
	validate = validator.New()

	// Configura o tradutor para Português do Brasil.
	ptBR := pt_BR.New()
	uni := ut.New(ptBR, ptBR)
	trans, _ = uni.GetTranslator("pt_BR")

	// Registra as traduções padrão para pt_BR no validador.
	ptBR_translations.RegisterDefaultTranslations(validate, trans)

	// Registra uma função para obter o nome do campo a partir da tag `json`.
	// Isso fará com que as mensagens de erro usem "id_esporte" em vez de "EsporteID".
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Registra validadores customizados
	validate.RegisterValidation("user_type", validateUserType)
}

// ValidateStruct executa a validação em qualquer struct passada.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// TranslateError traduz os erros de validação para mensagens amigáveis em português.
func TranslateError(err error) map[string]string {
	if err == nil {
		return nil
	}

	errorMessages := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errorMessages[e.Field()] = e.Translate(trans)
	}
	return errorMessages
}

// validateUserType é uma função de validação customizada para o tipo de usuário.
func validateUserType(fl validator.FieldLevel) bool {
	userType := fl.Field().String()
	switch userType {
	case "organizador", "jogador":
		return true
	}
	return false
}
