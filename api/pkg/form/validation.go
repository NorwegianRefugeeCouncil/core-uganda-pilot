package form

import (
	"fmt"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"strconv"
)

func (f *FormElement) Validate(errList validation.ErrorList, path *validation.Path) {
	switch f.Type {
	case Checkbox:
		for i, option := range f.Attributes.CheckboxOptions {
			if option.Required && !utils.Contains(f.Attributes.Value, strconv.Itoa(i)) {
				err := validation.Required(path.Child(f.Attributes.Name).Index(i), fmt.Sprintf("%s is required", f.Attributes.Name))
				errList = append(errList, err)
			}
		}
		fallthrough
	default:
		if f.Validation.Required && utils.AllEmpty(f.Attributes.Value) {
			err := validation.Required(path.Child(f.Attributes.Name), fmt.Sprintf("%s is required", f.Attributes.Name))
			errList = append(errList, err)
		}
		break
	}
}
