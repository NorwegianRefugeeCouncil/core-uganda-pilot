import {of, OperatorFunction, zip} from "rxjs"
import {map, switchMap, tap} from "rxjs/operators"
import {ajax} from "rxjs/ajax"

export class Database {
    public id: string = ""
    public name: string = ""
}

export class DatabaseList {
    public items: Database[] = []
}

export enum FieldKind {
    Text = "text",
    Reference = "reference",
    SubForm = "subform",
}

export class FieldType {
    public text?: FieldTypeText
    public reference?: FieldTypeReference
    public subForm?: FieldTypeSubForm
}

export class FieldTypeText {
}

export class FieldTypeMultiLineText {
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


export interface DatabaseCreator {
    createDatabase(): OperatorFunction<Partial<Database>, Database>
}

export interface DatabaseLister {
    getDatabases(): OperatorFunction<null, DatabaseList>
}

export interface FormLister {
    getForms(): OperatorFunction<null, FormDefinitionList>
}

export interface FormCreator {
    createForm(): OperatorFunction<FormDefinition, FormDefinition>
}

export interface RecordCreator {
    postRecord(): OperatorFunction<LocalRecord, { request: { record: LocalRecord }, response: { record: Record } }>
}

export interface FolderListOptions {

}

export interface FolderLister {
    listFolders(): OperatorFunction<FolderListOptions, { request: { options: FolderListOptions }, result: FolderList }>
}

export interface FolderCreator {
    createFolder(): OperatorFunction<Partial<Folder>, Folder>
}

export type RecordListOptions = { databaseId: string, formId: string }

export interface RecordLister {
    fetchRecords(): OperatorFunction<RecordListOptions, { request: { options: RecordListOptions }, result: RecordList }>
}

export interface Client
    extends DatabaseCreator,
        DatabaseLister,
        FormLister,
        FormCreator,
        RecordCreator,
        FolderLister,
        FolderCreator {
}

export class client implements Client {
    public address = "http://localhost:9000"

    createDatabase(): OperatorFunction<Partial<Database>, Database> {
        return (s) => s.pipe(
            tap(console.log),
            switchMap((s) => ajax({url: `${this.address}/databases`, method: "POST", responseType: "json", body: s})),
            map(r => {
                return r.response as Database
            })
        )
    }

    getDatabases(): OperatorFunction<null, DatabaseList> {
        return () => of(null).pipe(
            switchMap(() => ajax({url: `${this.address}/databases`, method: "GET", responseType: "json",})),
            map(r => {
                return r.response as DatabaseList
            })
        )
    }

    getForms(): OperatorFunction<null, FormDefinitionList> {
        return () => of(null).pipe(
            switchMap(() => ajax({url: `${this.address}/forms`, method: "GET", responseType: "json",})),
            map(r => {
                return r.response as FormDefinitionList
            })
        )
    }

    createForm(): OperatorFunction<FormDefinition, FormDefinition> {
        return (s) => s.pipe(
            tap(console.log),
            switchMap((s) => ajax({url: `${this.address}/forms`, method: "POST", responseType: "json", body: s})),
            map(r => {
                return r.response as FormDefinition
            })
        )
    }

    postRecord(): OperatorFunction<LocalRecord, { request: { record: LocalRecord }, response: { record: Record } }> {
        return (s) => s.pipe(
            switchMap((s) => {
                const {databaseId, formId, values} = s
                return zip(
                    of(s),
                    ajax({
                        url: `${this.address}/databases/${databaseId}/forms/${formId}/records`,
                        method: "POST",
                        responseType: "json",
                        body: {
                            formId,
                            databaseId,
                            values
                        }
                    }))
            }),
            map(([localRecord, response]) => {
                return {request: {record: localRecord}, response: {record: response.response as Record}}
            })
        )
    }

    fetchRecords(): OperatorFunction<RecordListOptions, { request: { options: RecordListOptions }, result: RecordList }> {
        return (s) => s.pipe(
            switchMap((o) => {
                return zip(
                    of(o),
                    ajax({
                        url: `${this.address}/databases/${o.databaseId}/forms/${o.formId}/records`,
                        method: "GET",
                        responseType: "json",
                    })
                )
            }),
            map(([options, r]) => {
                return {request: {options}, result: r.response as RecordList}
            })
        )
    }

    createFolder(): OperatorFunction<Partial<Folder>, Folder> {
        return (s) => s.pipe(
            switchMap((s) => ajax({url: `${this.address}/folders`, method: "POST", responseType: "json", body: s})),
            map(r => {
                return r.response as Folder
            })
        )
    }

    listFolders(): OperatorFunction<FolderListOptions, { request: { options: FolderListOptions }, result: FolderList }> {
        return (s) => s.pipe(
            switchMap((o) => {
                return zip(
                    of(o),
                    ajax({
                        url: `${this.address}/folders`,
                        method: "GET",
                        responseType: "json",
                    })
                )
            }),
            map(([options, r]) => {
                return {request: {options}, result: r.response as FolderList}
            })
        )
    }

}


