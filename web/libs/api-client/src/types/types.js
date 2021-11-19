"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
exports.__esModule = true;
exports.LocalRecord = exports.Record = exports.FolderList = exports.Folder = exports.FormDefinitionList = exports.FormDefinition = exports.FieldDefinition = exports.FieldTypeSubForm = exports.FieldTypeReference = exports.FieldTypeSingleSelect = exports.FieldTypeQuantity = exports.FieldTypeDate = exports.FieldTypeMultilineText = exports.FieldTypeText = exports.FieldType = exports.FieldKind = exports.DatabaseList = exports.Database = void 0;
var Database = /** @class */ (function () {
    function Database() {
        this.id = "";
        this.name = "";
    }
    return Database;
}());
exports.Database = Database;
var DatabaseList = /** @class */ (function () {
    function DatabaseList() {
        this.items = [];
    }
    return DatabaseList;
}());
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
var FieldType = /** @class */ (function () {
    function FieldType() {
    }
    return FieldType;
}());
exports.FieldType = FieldType;
var FieldTypeText = /** @class */ (function () {
    function FieldTypeText() {
    }
    return FieldTypeText;
}());
exports.FieldTypeText = FieldTypeText;
var FieldTypeMultilineText = /** @class */ (function () {
    function FieldTypeMultilineText() {
    }
    return FieldTypeMultilineText;
}());
exports.FieldTypeMultilineText = FieldTypeMultilineText;
var FieldTypeDate = /** @class */ (function () {
    function FieldTypeDate() {
    }
    return FieldTypeDate;
}());
exports.FieldTypeDate = FieldTypeDate;
var FieldTypeQuantity = /** @class */ (function () {
    function FieldTypeQuantity() {
    }
    return FieldTypeQuantity;
}());
exports.FieldTypeQuantity = FieldTypeQuantity;
var FieldTypeSingleSelect = /** @class */ (function () {
    function FieldTypeSingleSelect() {
    }
    return FieldTypeSingleSelect;
}());
exports.FieldTypeSingleSelect = FieldTypeSingleSelect;
var FieldTypeReference = /** @class */ (function () {
    function FieldTypeReference() {
        this.databaseId = "";
        this.formId = "";
    }
    return FieldTypeReference;
}());
exports.FieldTypeReference = FieldTypeReference;
var FieldTypeSubForm = /** @class */ (function () {
    function FieldTypeSubForm() {
        this.id = "";
        this.name = "";
        this.code = "";
        this.fields = [];
    }
    return FieldTypeSubForm;
}());
exports.FieldTypeSubForm = FieldTypeSubForm;
var FieldDefinition = /** @class */ (function () {
    function FieldDefinition() {
        this.id = "";
        this.code = "";
        this.name = "";
        this.description = "";
        this.required = false;
        this.options = [];
        this.key = false;
        this.fieldType = new FieldType();
    }
    return FieldDefinition;
}());
exports.FieldDefinition = FieldDefinition;
var FormDefinition = /** @class */ (function () {
    function FormDefinition() {
        this.id = "";
        this.code = "";
        this.databaseId = "";
        this.folderId = "";
        this.name = "";
        this.fields = [];
    }
    return FormDefinition;
}());
exports.FormDefinition = FormDefinition;
var FormDefinitionList = /** @class */ (function () {
    function FormDefinitionList() {
        this.items = [];
    }
    return FormDefinitionList;
}());
exports.FormDefinitionList = FormDefinitionList;
var Folder = /** @class */ (function () {
    function Folder() {
        this.id = "";
        this.databaseId = "";
        this.parentId = "";
        this.name = "";
    }
    return Folder;
}());
exports.Folder = Folder;
var FolderList = /** @class */ (function () {
    function FolderList() {
        this.items = [];
    }
    return FolderList;
}());
exports.FolderList = FolderList;
var Record = /** @class */ (function () {
    function Record() {
        this.id = "";
        this.databaseId = "";
        this.formId = "";
        this.parentId = undefined;
        this.values = {};
    }
    return Record;
}());
exports.Record = Record;
var LocalRecord = /** @class */ (function (_super) {
    __extends(LocalRecord, _super);
    function LocalRecord() {
        var _this = _super !== null && _super.apply(this, arguments) || this;
        _this.isNew = false;
        return _this;
    }
    return LocalRecord;
}(Record));
exports.LocalRecord = LocalRecord;
