// Code generated by "stringer -type=FieldKind"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FieldKindUnknown-0]
	_ = x[FieldKindText-1]
	_ = x[FieldKindSubForm-2]
	_ = x[FieldKindReference-3]
	_ = x[FieldKindMultilineText-4]
	_ = x[FieldKindDate-5]
	_ = x[FieldKindQuantity-6]
	_ = x[FieldKindMonth-7]
	_ = x[FieldKindSingleSelect-8]
}

const _FieldKind_name = "FieldKindUnknownFieldKindTextFieldKindSubFormFieldKindReferenceFieldKindMultilineTextFieldKindDateFieldKindQuantityFieldKindMonthFieldKindSingleSelect"

var _FieldKind_index = [...]uint8{0, 16, 29, 45, 63, 85, 98, 115, 129, 150}

func (i FieldKind) String() string {
	if i < 0 || i >= FieldKind(len(_FieldKind_index)-1) {
		return "FieldKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldKind_name[_FieldKind_index[i]:_FieldKind_index[i+1]]
}
