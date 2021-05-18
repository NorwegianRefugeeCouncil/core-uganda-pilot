package core

func CanHaveChildren(elementType FormElementType) bool {
	switch elementType {
	case SectionType:
		return true
	default:
		return false
	}
}

func MustHaveChildren(elementType FormElementType) bool {
	switch elementType {
	case SectionType:
		return true
	default:
		return false
	}
}
