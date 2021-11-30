export type Session = {
    active: boolean
    expiry: string
    expiredInSeconds: number
    subject: string
    username: string
}

export class Database {
    public id: string = ""
    public name: string = ""
}

export class DatabaseList {
    public items: Database[] = []
}

export enum FieldKind {
    Text = "text",
    MultilineText = "multilineText",
    Reference = "reference",
    SubForm = "subform",
    Date = "date",
    Quantity = "quantity",
    SingleSelect = "singleSelect",
}

export class FieldType {
    public text?: FieldTypeText
    public reference?: FieldTypeReference
    public subForm?: FieldTypeSubForm
    public multilineText?: FieldTypeMultilineText
    public date?: FieldTypeDate
    public quantity?: FieldTypeQuantity
    public singleSelect?: FieldTypeSingleSelect
}

export class FieldTypeText {
}

export class FieldTypeMultilineText {
}

export class FieldTypeDate {
}

export class FieldTypeQuantity {
}

export class FieldTypeSingleSelect {
}

export class FieldTypeReference {
    public databaseId: string = ""
    public formId: string = ""
}

export class FieldTypeSubForm {
    public id: string = ""
    public name: string = ""
    public code: string = ""
    public fields: FieldDefinition[] = []
}

export class FieldDefinition {
    public id: string = ""
    public code: string = ""
    public name: string = ""
    public description: string = ""
    public required: boolean = false
    public options: string[] = []
    public key: boolean = false
    public fieldType: FieldType = new FieldType()
}

export class FormDefinition {
    public id: string = ""
    public code: string = ""
    public databaseId: string = ""
    public folderId: string = ""
    public name: string = ""
    public fields: FieldDefinition[] = []

    constructor(id = "", code = "", databaseId = "", folderId = "", name = "", fields = []) {
        this.id = id;
        this.code = code;
        this.databaseId = databaseId;
        this.folderId = folderId;
        this.name = name;
        this.fields = fields;
    }
}

export class FormDefinitionList {
    public items: FormDefinition[] = []
}

export class Folder {
    public id: string = ""
    public databaseId: string = ""
    public parentId: string = ""
    public name: string = ""
}

export class FolderList {
    public items: Folder[] = []
}

export class RecordInstance {
    public id: string = ""
    public databaseId: string = ""
    public formId: string = ""
    public parentId: string | undefined = undefined
    public values: { [key: string]: any } = {}
}

export class LocalRecord extends RecordInstance {
    public isNew: boolean = false
}

export type RecordList = { items: RecordInstance[] }
