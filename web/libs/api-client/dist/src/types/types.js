"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LocalRecord = exports.Record = exports.FolderList = exports.Folder = exports.FormDefinitionList = exports.FormDefinition = exports.FieldDefinition = exports.FieldTypeSubForm = exports.FieldTypeReference = exports.FieldTypeSingleSelect = exports.FieldTypeQuantity = exports.FieldTypeDate = exports.FieldTypeMultilineText = exports.FieldTypeText = exports.FieldType = exports.FieldKind = exports.DatabaseList = exports.Database = void 0;
class Database {
    constructor() {
        this.id = "";
        this.name = "";
    }
}
exports.Database = Database;
class DatabaseList {
    constructor() {
        this.items = [];
    }
}
exports.DatabaseList = DatabaseList;
var FieldKind;
(function (FieldKind) {
    FieldKind["Text"] = "text";
    FieldKind["MultilineText"] = "multilineText";
    FieldKind["Reference"] = "reference";
    FieldKind["SubForm"] = "subform";
    FieldKind["Date"] = "date";
    FieldKind["Quantity"] = "quantity";
    FieldKind["SingleSelect"] = "singleSelect";
})(FieldKind = exports.FieldKind || (exports.FieldKind = {}));
class FieldType {
}
exports.FieldType = FieldType;
class FieldTypeText {
}
exports.FieldTypeText = FieldTypeText;
class FieldTypeMultilineText {
}
exports.FieldTypeMultilineText = FieldTypeMultilineText;
class FieldTypeDate {
}
exports.FieldTypeDate = FieldTypeDate;
class FieldTypeQuantity {
}
exports.FieldTypeQuantity = FieldTypeQuantity;
class FieldTypeSingleSelect {
}
exports.FieldTypeSingleSelect = FieldTypeSingleSelect;
class FieldTypeReference {
    constructor() {
        this.databaseId = "";
        this.formId = "";
    }
}
exports.FieldTypeReference = FieldTypeReference;
class FieldTypeSubForm {
    constructor() {
        this.id = "";
        this.name = "";
        this.code = "";
        this.fields = [];
    }
}
exports.FieldTypeSubForm = FieldTypeSubForm;
class FieldDefinition {
    constructor() {
        this.id = "";
        this.code = "";
        this.name = "";
        this.description = "";
        this.required = false;
        this.options = [];
        this.key = false;
        this.fieldType = new FieldType();
    }
}
exports.FieldDefinition = FieldDefinition;
class FormDefinition {
    constructor() {
        this.id = "";
        this.code = "";
        this.databaseId = "";
        this.folderId = "";
        this.name = "";
        this.fields = [];
    }
}
exports.FormDefinition = FormDefinition;
class FormDefinitionList {
    constructor() {
        this.items = [];
    }
}
exports.FormDefinitionList = FormDefinitionList;
class Folder {
    constructor() {
        this.id = "";
        this.databaseId = "";
        this.parentId = "";
        this.name = "";
    }
}
exports.Folder = Folder;
class FolderList {
    constructor() {
        this.items = [];
    }
}
exports.FolderList = FolderList;
class Record {
    constructor() {
        this.id = "";
        this.databaseId = "";
        this.formId = "";
        this.parentId = undefined;
        this.values = {};
    }
}
exports.Record = Record;
class LocalRecord extends Record {
    constructor() {
        super(...arguments);
        this.isNew = false;
    }
}
exports.LocalRecord = LocalRecord;
