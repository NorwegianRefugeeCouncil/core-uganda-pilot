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

export class Record {
    public id: string = ""
    public databaseId: string = ""
    public formId: string = ""
    public parentId: string | undefined = undefined
    public values: { [key: string]: any } = {}
}

export class LocalRecord extends Record {
    public isNew: boolean = false
}

export type RecordList = { items: Record[] }

export class Organization {
    public id: string = ""
    public key: string = ""
    public name: string = ""
}

export type OrganizationList = { items: Organization[] }

export class IdentityProvider {
    public id: string = ""
    public organizationId: string = ""
    public kind: string = ""
    public domain: string = ""
    public clientId: string = ""
    public clientSecret: string = ""
}

export type IdentityProviderList = { items: IdentityProvider[] }
