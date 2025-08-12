package processors

import "encoding/json"

// ValidateJSON vérifie si une chaîne est un JSON valide
func ValidateJSON(input string) error {
	var js interface{}
	return json.Unmarshal([]byte(input), &js)
}

// PrettyValidationError formate les erreurs de validation pour l'affichage
func PrettyValidationError(err error) string {
	if err == nil {
		return ""
	}

	switch e := err.(type) {
	case *json.SyntaxError:
		return "Erreur de syntaxe à la position " + string(e.Offset) + ": " + e.Error()
	case *json.UnmarshalTypeError:
		return "Type incorrect à la position " + string(e.Offset) +
			". Attendu: " + e.Type.String() + ", Reçu: " + e.Value
	default:
		return "Erreur de validation: " + err.Error()
	}
}
