"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getFieldKind = void 0;
const types_1 = require("../types/types");
function getFieldKind(fieldType) {
    if (fieldType.text) {
        return types_1.FieldKind.Text;
    }
    if (fieldType.multilineText) {
        return types_1.FieldKind.MultilineText;
    }
    if (fieldType.date) {
        return types_1.FieldKind.Date;
    }
    if (fieldType.subForm) {
        return types_1.FieldKind.SubForm;
    }
    if (fieldType.reference) {
        return types_1.FieldKind.Reference;
    }
    if (fieldType.quantity) {
        return types_1.FieldKind.Quantity;
    }
    if (fieldType.singleSelect) {
        return types_1.FieldKind.SingleSelect;
    }
    throw new Error("unknown field kind");
}
exports.getFieldKind = getFieldKind;
