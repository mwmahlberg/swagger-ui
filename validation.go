package swaggerui

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"gopkg.in/yaml.v3"
)

var (
	// RegexValidFilename matches a valid filename for a swagger spec file.
	RegexValidFilename = regexp.MustCompile(`(?i)\.(y[a]?ml|json)$`)
)

func init() {
	govalidator.CustomTypeTagMap.Set("isYaml", isYaml)
	govalidator.CustomTypeTagMap.Set("correctContent", isCorrectContent)
	govalidator.CustomTypeTagMap.Set("acceptedFileName", isAcceptedFileName)

}

// isYaml checks if i is valid JSON data.
// It also explcitly checks that i is not JSON data, since JSON parses as YAML.
func isYaml(i interface{}, o interface{}) bool {
	var tmp = make(map[string]interface{})

	switch v := i.(type) {
	case []byte:
		return yaml.Unmarshal(v, &tmp) == nil && !json.Valid(v)
	case string:
		return yaml.Unmarshal([]byte(v), &tmp) == nil && !json.Valid([]byte(v))
	default:
		return false
	}

}

func isAcceptedFileName(i interface{}, o interface{}) bool {
	v, isString := i.(string)
	if !isString {
		return false
	}
	return RegexValidFilename.MatchString(v)
}

func isCorrectContent(i, o interface{}) bool {

	h, isHandler := o.(SwaggerUi)

	if !isHandler {
		return false
	}

	var foo string
	switch v := i.(type) {
	case string:
		foo = v
	case []byte:
		foo = string(v)
	default:
		return false
	}

	switch {
	case strings.HasSuffix(strings.ToLower(h.specFilename), ".yaml"):
		return isYaml(i, o)
	case strings.HasSuffix(strings.ToLower(h.specFilename), ".json"):
		return govalidator.IsJSON(foo)
	default:
		return false
	}

}
