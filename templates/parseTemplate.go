package templates

import "strings"

func ParseTemplate(template string, data map[string]string) string {
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		template = strings.ReplaceAll(template, placeholder, value)
	}

	return template
}
