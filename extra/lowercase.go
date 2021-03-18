package extra

import (
	"strings"
	"unicode"

	"github.com/let-s-go/jsoniter"
)

func init() {
	jsoniter.RegisterExtension(&lowercaseNamingExtension{})
}

type lowercaseNamingExtension struct {
	jsoniter.DummyExtension
}

func (extension *lowercaseNamingExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		tags := strings.Split(binding.Field.Tag.Get("json"), ",")
		if len(tags) == 0 || tags[0] == "" {
			binding.ToNames = []string{extension.toLowerCase(binding.Field.Name)}
			binding.FromNames = []string{extension.toLowerCase(binding.Field.Name)}
		}
	}
}

func (extension *lowercaseNamingExtension) toLowerCase(name string) string {
	newName := []rune(name)
	if len(newName) > 0 {
		newName[0] = unicode.ToLower(newName[0])
	}
	return string(newName)
}
