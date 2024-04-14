package ponzu_driver_postgres

import (
	ponzuContent "github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/content/generator/types"
)

func (database *driver) Generate(contentType ponzuContent.Type, typeDefinition *types.TypeDefinition) error {
	return nil
}

func (database *driver) ValidateField(field *types.Field) error {
	return nil
}
