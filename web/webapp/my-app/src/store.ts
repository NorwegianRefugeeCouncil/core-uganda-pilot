import {
    client,
    Database,
    DatabaseCreator,
    DatabaseLister, FieldDefinition,
    FieldKind, Folder, FolderCreator, FolderLister, FolderListOptions,
    FormCreator,
    FormDefinition,
    FormLister, LocalRecord,
    Record,
    RecordCreator,
    RecordLister,
    RecordListOptions
} from "./client";
import {createReducer, on} from "./store/store";
import {Action, createAction, props} from "./store/actions";
import {
    BehaviorSubject,
    distinctUntilChanged,
    filter,
    Observable,
    of,
    ReplaySubject,
    shareReplay,
    Subject,
    Subscription
} from "rxjs";
import {ofType} from "./store/effects";
import {ActionReducer} from "./store/models";
import {catchError, map, switchMap} from "rxjs/operators";
import {v4 as uuidv4} from 'uuid';

export type AddFieldProps = {
    name: string,
    code: string,
    description: string,
    isKey: boolean,
    isRequired: boolean,
    type: undefined | FieldKind,
    subFormIdx: number,
}

export type FieldPropsState = AddFieldProps & {
    index: number
    formIdx: number
}

export type FormPropsState = {
    id: string,
    dirty: boolean,
    name: string,
    code: string,
    parentFormIdx: number,
    parentFieldIdx: number,
    index: number,
}

export type FormEditorState = {
    selectedFormIdx: number
    selectedFieldIdx: number
    fields: FieldPropsState[],
    forms: FormPropsState[]
    success: boolean
    error: any
    pending: boolean,
    createdFormId: string | null
}

export interface State {
    pending: boolean,
    error: any
    databaseIds: { [id: string]: number }
    databases: Database[]
    formIds: { [id: string]: number }
    forms: FormDefinition[]
    folderIds: { [id: string]: number }
    folders: Folder[]
    recordIds: { [id: string]: number }
    records: (LocalRecord | undefined)[]
    selectedDatabaseId: string
    selectedFormId: string
    selectedSubFormId: string
    selectedFolderId: string
    selectedRecordId: string
    createDatabaseName: string,
    createDatabaseId: string,
    createDatabaseError: any,
    createDatabasePending: boolean,
    createDatabaseSuccess: boolean,
    createRecordError: any,
    createRecordPending: boolean,
    createRecordSuccess: boolean,
    createFolderError: any,
    createFolderPending: boolean,
    createFolderSuccess: boolean,
    createFolderName: string,
    createFolderId: string,
    fetchFoldersError: any,
    fetchFoldersPending: boolean,
    fetchFoldersSuccess: boolean,
    form: FormEditorState
}

export const fieldPropsInitialState: FieldPropsState = {
    index: 0,
    name: "",
    code: "",
    isRequired: false,
    description: "",
    isKey: false,
    type: undefined,
    formIdx: 0,
    subFormIdx: -1,
}

export const formPropsInitialState: FormPropsState = {
    id: "",
    dirty: false,
    name: "",
    code: "",
    parentFormIdx: -1,
    parentFieldIdx: -1,
    index: 0,
}

export const formInitialState: FormEditorState = {
    selectedFormIdx: 0,
    selectedFieldIdx: -1,
    fields: [],
    forms: [formPropsInitialState],
    success: false,
    error: null,
    pending: false,
    createdFormId: "",
}

export const initialState: State = {
    pending: false,
    error: null,
    databaseIds: {},
    databases: [],
    formIds: {},
    forms: [],
    recordIds: {},
    records: [],
    folderIds: {},
    folders: [],
    selectedDatabaseId: "",
    selectedFormId: "",
    selectedFolderId: "",
    selectedRecordId: "",
    selectedSubFormId: "",
    createDatabaseName: "",
    createDatabaseError: null,
    createDatabasePending: false,
    createDatabaseSuccess: false,
    createDatabaseId: "",
    createRecordError: null,
    createRecordPending: false,
    createRecordSuccess: false,
    createFolderError: null,
    createFolderPending: false,
    createFolderSuccess: false,
    createFolderName: "",
    createFolderId: "",
    fetchFoldersError: null,
    fetchFoldersPending: false,
    fetchFoldersSuccess: false,
    form: formInitialState,
}

export const fetchDatabases = createAction("[databases] fetch")
export const fetchDatabasesSuccess = createAction("[databases] fetch success", props<{ databases: Database[] }>())
export const fetchDatabasesError = createAction("[databases] fetch error", props<{ error: any }>())

export const fetchForms = createAction("[forms] fetch")
export const fetchFormsSuccess = createAction("[forms] fetch success", props<{ forms: FormDefinition[] }>())
export const fetchFormsError = createAction("[forms] fetch error", props<{ error: any }>())

export const fetchRecords = createAction("[records] fetch", props<{ options: RecordListOptions }>())
export const fetchRecordsSuccess = createAction("[records] fetch success", props<{ request: { options: RecordListOptions }, result: Record[] }>())
export const fetchRecordsError = createAction("[records] fetch error", props<{ error: any }>())
export const createRecord = createAction("[records] create", props<{ record: LocalRecord }>())
export const createRecordSuccess = createAction("[records] create success", props<{ request: { record: LocalRecord }, response: { record: Record } }>())
export const createRecordError = createAction("[records] create error", props<{ error: any }>())
export const createRecordCancel = createAction("[records] reset")
export const setSelectedRecord = createAction("[records] set selected", props<{ id: string }>())
export const setRecordValue = createAction("[records] set value", props<{ id: string, key: string, value: any }>())
export const setRecordValues = createAction("[records] set value", props<{ id: string, values: { [key: string]: any } }>())
export const recordOpenSubForm = createAction("[records] open sub form", props<{ id: string }>())
export const recordCloseSubForm = createAction("[records] close sub form")
export const initNewRecord = createAction("[records] init new record", props<{ databaseId: string, formId: string, parentId: string | undefined }>())

export const createDatabaseReset = createAction("[databases] reset")
export const createDatabaseSetName = createAction("[databases] set name", props<{ name: string }>())
export const createDatabase = createAction("[databases] create", props<{ database: Partial<Database> }>())
export const createDatabaseSuccess = createAction("[databases] create success", props<{ database: Database }>())
export const createDatabaseError = createAction("[databases] create error", props<{ error: any }>())

export const createFormReset = createAction("[forms] reset")
export const createFormSetName = createAction("[forms] set name", props<{ name: string }>())
export const createForm = createAction("[forms] create", props<{ form: FormDefinition }>())
export const createFormSuccess = createAction("[forms] create success", props<{ form: FormDefinition }>())
export const createFormError = createAction("[forms] create error", props<{ error: any }>())
export const createFormSetSelectedForm = createAction("[forms] set selected form", props<{ index: number }>())
export const createFormSetSelectedField = createAction("[forms] set selected field", props<{ index: number }>())
export const createFormSetSelectedFieldCode = createAction("[forms] set field code", props<{ code: string }>())
export const createFormSetSelectedFieldName = createAction("[forms] set field name", props<{ name: string }>())
export const createFormSetSelectedFieldRequired = createAction("[forms] set field required", props<{ required: boolean }>())
export const createFormSetSelectedFieldIsKey = createAction("[forms] set field key", props<{ isKey: boolean }>())
export const createFormSetSelectedFieldDescription = createAction("[forms] set field description", props<{ description: string }>())
export const createFormSetSelectedFieldType = createAction("[forms] set field kind", props<{ kind: FieldKind }>())
export const createFormSetOpenSubForm = createAction("[forms] open subform")
export const createFormSetSaveSubForm = createAction("[forms] save subform")
export const createFormAddField = createAction("[forms] add field", props<{ field: AddFieldProps, select?: boolean }>())
export const createFormSaveField = createAction("[forms] save field", props<{ index: number }>())
export const createFormRemoveField = createAction("[forms] add field", props<{ index: number }>())

export const fetchFolders = createAction("[folders] fetch", props<{ options: FolderListOptions }>())
export const fetchFoldersError = createAction("[folders] fetch error", props<{ error: any }>())
export const fetchFoldersSuccess = createAction("[folders] fetch success", props<{ request: { options: FolderListOptions }, result: Folder[] }>())
export const createFolderReset = createAction("[folders] reset")
export const createFolderSetName = createAction("[folders] set name", props<{ name: string }>())
export const createFolder = createAction("[folders] create", props<{ folder: Partial<Folder> }>())
export const createFolderSuccess = createAction("[folders] create success", props<{ folder: Folder }>())
export const createFolderError = createAction("[folders] create error", props<{ error: any }>())
export const setSelectedFolder = createAction("[folders] set selected", props<{ id: string | undefined }>())

export const setSelectedDatabase = createAction("[databases] set selected", props<{ id: string }>())
export const setSelectedForm = createAction("[form] set selected", props<{ id: string }>())


export const reducer = createReducer<State>(initialState,
    on(fetchDatabases, state => ({...state, pending: true})),
    on(fetchDatabasesSuccess, (state, {databases}) => {
        const newState = {...state}
        newState.databases = databases
        newState.databaseIds = {}
        for (let i = 0; i < databases.length; i++) {
            const db = databases[i]
            newState.databaseIds[db.id] = i
        }
        return newState;
    }),
    on(fetchDatabasesError, (state, {error}) => ({...state, pending: false, error})),
    on(fetchForms, state => ({...state, pending: true})),
    on(fetchFormsSuccess, (state, {forms}) => {
        const newState = {...state}
        newState.forms = forms
        newState.formIds = {}
        for (let i = 0; i < forms.length; i++) {
            const form = forms[i]
            newState.formIds[form.id] = i
        }
        return newState;
    }),
    on(fetchFormsError, (state, {error}) => ({...state, pending: false, error})),
    on(fetchRecords, state => ({...state, pending: true})),
    on(fetchRecordsSuccess, (state, {request, result}) => {

        let records = result ? result : []

        const newState = {
            ...state,
            records: [...state.records],
            recordIds: {...state.recordIds}
        }

        for (let record of records) {
            const localRecord = record as LocalRecord
            localRecord.isNew = false
            if (newState.recordIds.hasOwnProperty(record.id)) {
                const currentRecordIdx = newState.recordIds[record.id]
                newState.records[currentRecordIdx] = localRecord
            } else {
                newState.records.push(localRecord)
                newState.recordIds[record.id] = newState.records.length - 1
            }
        }

        return newState

    }),
    on(fetchRecordsError, (state, {error}) => ({...state, pending: false, error})),
    on(createDatabaseSetName, (state, {name}) => ({...state, createDatabaseName: name})),
    on(createDatabase, (state) => ({...state, createDatabasePending: true})),
    on(setSelectedDatabase, (state, {id}) => state.selectedDatabaseId === id ? state : ({
        ...state,
        selectedDatabaseId: id
    })),
    on(setSelectedForm, (state, {id}) => state.selectedFormId === id ? state : ({
        ...state,
        selectedFormId: id
    })),
    on(setSelectedFolder, (state, {id}) => state.selectedFolderId === id ? state : ({
        ...state,
        selectedFolderId: id ? id : "",
    })),
    on(createDatabaseReset, (state) => ({
        ...state,
        createDatabasePending: false,
        createDatabaseError: null,
        createDatabaseSuccess: false,
        createDatabaseName: "",
    })),
    on(createDatabaseSuccess, (state, {database}) => {

        const newState = {
            ...state,
            createDatabasePending: false,
            createDatabaseError: null,
            createDatabaseSuccess: true,
            createDatabaseId: database.id,
        }

        newState.databases.push(database)
        newState.databaseIds[database.id] = newState.databases.length - 1

        return newState
    }),
    on(createDatabaseError, (state, {error}) => ({
        ...state,
        createDatabasePending: false,
        createDatabaseError: error,
        createDatabaseSuccess: false
    })),
    on(createFormReset, state => ({
        ...state,
        form: formInitialState
    })),
    on(createForm, (state) => {
        return {
            ...state,
            form: {
                ...state.form,
                pending: true
            }
        }
    }),
    on(createFormSuccess, (state, {form}) => {
        const newState = {
            ...state,
            form: {
                ...state.form,
                ...formInitialState,
                success: true,
                createdFormId: form.id,
            }
        }
        newState.forms.push(form)
        newState.formIds[form.id] = newState.forms.length - 1
        return newState

    }),
    on(createFormError, (state, {error}) => ({
        ...state,
        form: {
            ...state.form,
            pending: false,
            error
        }
    })),
    on(createFormSetName, (state, {name}) => {
        const newForms = [...state.form.forms]
        const formIdx = state.form.selectedFormIdx
        newForms[formIdx] = {
            ...state.form.forms[formIdx],
            name
        }
        return {
            ...state,
            form: {
                ...state.form,
                forms: newForms
            }
        }
    }),
    on(createFormSetSelectedForm, (state, {index}) => {
        return {
            ...state,
            form: {
                ...state.form,
                selectedFormIdx: index,
                selectedFieldIdx: -1,
            }
        }
    }),
    on(createFormSetSelectedField, (state, {index}) => {
        return {
            ...state,
            form: {
                ...state.form,
                selectedFieldIdx: index,
            }
        }
    }),
    on(createFormAddField, (state, {field, select}) => {
        const formIdx = state.form.selectedFormIdx
        const newField: FieldPropsState = {
            ...field,
            formIdx,
            index: state.form.fields.length
        }
        const newFields = [...state.form.fields, newField]
        return {
            ...state,
            form: {
                ...state.form,
                fields: newFields,
                selectedFieldIdx: select ? newFields.length - 1 : state.form.selectedFieldIdx,
            }
        }
    }),
    on(createFormRemoveField, (state, {index}) => {
        const newFields = state.form.fields.filter((_, i) => i !== index).map((f, i) => ({...f, index: i}))
        return {
            ...state,
            form: {
                ...state.form,
                fields: newFields
            }
        }
    }),
    on(createFormSetSelectedFieldCode, (state, {code}) => {
        const newFields = [...state.form.fields]
        const field = newFields[state.form.selectedFieldIdx]
        field.code = code
        return {
            ...state,
            form: {
                ...state.form,
                fields: newFields
            }
        };
    }),
    on(createFormSetSelectedFieldName, (state, {name}) => {
        const newFields = [...state.form.fields]
        const field = newFields[state.form.selectedFieldIdx]
        field.name = name
        return {
            ...state,
            form: {
                ...state.form,
                fields: newFields
            }
        };
    }),
    on(createFormSetSelectedFieldRequired, (state, {required}) => {
        const newFields = [...state.form.fields]
        const field = newFields[state.form.selectedFieldIdx]
        field.isRequired = required
        return {
            ...state,
            form: {
                ...state.form,
                fields: newFields
            }
        };
    }),
    on(createFormSetSelectedFieldIsKey, (state, {isKey}) => {
        const newFields = [...state.form.fields]
        const field = newFields[state.form.selectedFieldIdx]
        field.isKey = isKey
        return {
            ...state, form: {...state.form, fields: newFields}
        };
    }),
    on(createFormSetSelectedFieldDescription, (state, {description}) => {
        const newFields = [...state.form.fields]
        const field = newFields[state.form.selectedFieldIdx]
        field.description = description
        return {
            ...state, form: {...state.form, fields: newFields}
        };
    }),
    on(createFormSetSelectedFieldType, (state, {kind}) => {
        const newFields = [...state.form.fields]
        const field = newFields[state.form.selectedFieldIdx]
        field.type = kind
        return {
            ...state, form: {...state.form, fields: newFields}
        };
    }),
    on(createFormSetOpenSubForm, (state) => {
        const newState = {
            ...state,
            form: {
                ...state.form,
                forms: [...state.form.forms]
            }
        }
        const field = newState.form.fields[newState.form.selectedFieldIdx]
        let subFormIdx = field.subFormIdx
        if (subFormIdx < 0) {
            const newSubForm = {...formPropsInitialState}
            newSubForm.parentFormIdx = newState.form.selectedFormIdx
            newSubForm.parentFieldIdx = newState.form.selectedFieldIdx
            newSubForm.index = newState.form.forms.length
            newState.form.forms.push(newSubForm)
            subFormIdx = newState.form.forms.length - 1
            field.subFormIdx = subFormIdx
        }
        newState.form.selectedFormIdx = subFormIdx
        newState.form.selectedFieldIdx = -1
        return newState
    }), on(createFormSetSaveSubForm, (state) => {
        const form = state.form.forms[state.form.selectedFormIdx]
        const parent = state.form.forms[form.parentFormIdx]
        return {
            ...state,
            form: {
                ...state.form,
                selectedFormIdx: parent.index,
                selectedFieldIdx: -1
            }
        }
    }),
    on(createFormSaveField, (state) => ({
        ...state
    })),
    on(setSelectedRecord, (state, {id}) => ({
        ...state,
        selectedRecordId: id,
    })),
    on(recordOpenSubForm, (state, {id}) => {
        const newState = {...state}
        newState.selectedSubFormId = id
        return newState
    }),
    on(recordCloseSubForm, (state) => {
        const newState = {...state}

        const findParentFormId = (parentId: string, fields: FieldDefinition[]): string => {
            for (let i = 0; i < fields.length; i++) {
                const field = fields[i]
                if (field.fieldType.subForm) {
                    if (field.fieldType.subForm.id === newState.selectedSubFormId) {
                        return field.fieldType.subForm.id
                    } else {
                        const parentFormId = findParentFormId(field.fieldType.subForm.id, field.fieldType.subForm.fields)
                        if (parentFormId) {
                            return parentFormId
                        }
                    }
                }
            }
            return ""
        }

        let parentFormId: string = ""
        for (let i = 0; i < state.forms.length; i++) {
            const form = state.forms[i]
            parentFormId = findParentFormId(form.id, form.fields)
            if (parentFormId) {
                break
            }
        }

        const currentForm = state.forms[state.formIds[newState.selectedFormId]]
        if (currentForm.id !== parentFormId) {
            newState.selectedSubFormId = parentFormId
        } else {
            newState.selectedSubFormId = ""
        }
        return newState
    }),
    on(initNewRecord, (state, {databaseId, formId, parentId}) => {
        const newRecord = new LocalRecord()
        newRecord.databaseId = databaseId
        newRecord.formId = formId
        newRecord.isNew = true
        newRecord.id = uuidv4()
        if (parentId) {
            newRecord.parentId = parentId
        }
        const newState = {...state}
        newState.recordIds[newRecord.id] = newState.records.length
        newState.records.push(newRecord)
        newState.selectedRecordId = newRecord.id
        return newState
    }),
    on(setRecordValue, (state, {id, key, value}) => {
        const recordIdx = state.recordIds[id]
        if (recordIdx >= 0) {
            const record = state.records[recordIdx]
            if (record) {
                const newRecord = {
                    ...record,
                    values: {
                        ...record.values,
                        ...{[key]: value},
                    }
                }
                const newState = {...state}
                newState.records[recordIdx] = newRecord
                return newState
            }
        }
        return state
    }),
    on(setRecordValues, (state, {id, values}) => {
        const recordIdx = state.recordIds[id]
        if (recordIdx >= 0) {
            const record = state.records[recordIdx]
            if (record) {
                const newRecord = {
                    ...record,
                    values: {
                        ...record.values,
                        ...values
                    }
                }
                const newState = {...state}
                newState.records[recordIdx] = newRecord
                return newState
            }
        }
        return state
    }),
    on(createRecord, (state) => ({
        ...state,
        createRecordPending: true,
        createRecordSuccess: false,
    })),
    on(createRecordSuccess, (state, {request, response}) => {
        const newState = {
            ...state,
            createRecordPending: false,
            createRecordError: null,
            createRecordSuccess: true,
        }

        const newRecordId = response.record.id;
        const newRecord = response.record as LocalRecord
        newRecord.isNew = false

        // Updating the records "parentId" to match the created record ID
        for (let i = 0; i < newState.records.length; i++) {
            const currentRecord = newState.records[i]
            if (currentRecord && currentRecord.parentId === newRecordId) {
                newState.records[i] = {...currentRecord, parentId: newRecordId}
            }
        }

        const oldRecordIdx = newState.recordIds[newRecordId]
        let oldRecord: Record | undefined
        if (oldRecordIdx >= 0) {
            oldRecord = newState.records[oldRecordIdx]
        }

        // Replacing the selected record ID
        if (oldRecord && oldRecord.id === newState.selectedRecordId) {
            newState.selectedRecordId = response.record.id
        }

        // Replacing/Pushing the record in the store
        if (!oldRecord) {
            newState.records[oldRecordIdx] = newRecord
        } else {
            newState.records.push(newRecord)
            newState.recordIds[response.record.id] = newState.records.length - 1
        }

        return newState
    }),
    on(createRecordError, (state, {error}) => ({
        ...state,
        createRecordPending: true,
        createRecordError: error,
        createRecordSuccess: false,
    })),
    on(createRecordCancel, (state) => {

        const newState = {
            ...state,
            createRecordPending: false,
            createRecordError: null,
            createRecordSuccess: false,
        }

        const selectedRecordIdx = state.selectedRecordId
        if (state.recordIds.hasOwnProperty(state.selectedRecordId)) {
            const selectedRecordIdx = state.recordIds[state.selectedRecordId]
            const selectedRecord = state.records[selectedRecordIdx]
            if (selectedRecord?.isNew) {
                delete newState.recordIds[selectedRecordIdx]
                newState.records[selectedRecordIdx] = undefined
            }
            if (selectedRecord?.parentId) {
                const parentIdx = state.recordIds[selectedRecord.parentId]
                const parent = state.records[parentIdx]
                if (parent) {
                    newState.selectedRecordId = parent.id
                }
            }

        } else {

        }

        return newState
    }),
    on(createFolderSetName, (state, {name}) => ({...state, createFolderName: name})),
    on(fetchFolders, (state) => ({
        ...state,
        fetchFoldersPending: true,
    })),
    on(fetchFoldersError, (state, {error}) => ({
        ...state,
        fetchFoldersPending: false,
        fetchFoldersSuccess: false,
        fetchFoldersError: error,
    })),
    on(fetchFoldersSuccess, (state, {result}) => {
        const folders = result
        const newState = {...state}
        newState.folders = folders
        newState.folderIds = {}
        for (let i = 0; i < folders.length; i++) {
            const folder = folders[i]
            newState.folderIds[folder.id] = i
        }
        return newState
    }),
    on(createFolder, (state) => ({
        ...state,
        createFolderPending: true,
    })),
    on(createFolderError, (state, {error}) => ({
        ...state,
        createFolderPending: false,
        createFolderSuccess: false,
        createFolderError: error,
    })),
    on(createFolderSuccess, (state, {folder}) => {
        const newState = {
            ...state,
            createFolderPending: false,
            createFolderSuccess: true,
            createFolderError: undefined,
        }
        newState.folders.push(folder)
        newState.folderIds[folder.id] = newState.folders.length - 1
        newState.createFolderId = folder.id
        return newState
    }),
    on(createFolderReset, (state) => ({
        ...state,
        createFolderPending: false,
        createFolderSuccess: false,
        createFolderError: undefined,
        createFolderName: "",
        createFolderId: "",
    }))
)

export interface FormIntf {
    fields: FieldDefinition[]
    id: string
    name: string
    code: string
}

interface Store<T> {
    dispatch(action: Action): void

    actions$: Observable<Action>
    state$: Observable<T>
    selectedDatabase$: Observable<Database | undefined>
    selectedForm$: Observable<FormDefinition | undefined>
    selectedDatabaseForms$: Observable<FormDefinition[] | undefined>
    createFormName$: Observable<string>,
    createFormFields$: Observable<FieldPropsState[]>,
    createFormAllFields$ : Observable<FieldPropsState[]>,
    createFormForms$: Observable<FormPropsState[]>,
    createFormSelectedFieldIdx$: Observable<number>,
    createFormSelectedField$: Observable<FieldPropsState | undefined>
    createFormSelectedFormIdx$: Observable<number>,
    createFormSelectedForm$: Observable<FormPropsState | undefined>
    createFormSelectedFormIsSubForm$: Observable<boolean>,
    createFormParentFormIdx$: Observable<number>,
    createFormParentForm$: Observable<FormPropsState | undefined>,
    postRecordSuccess$: Observable<boolean>
    selectedRecord: Observable<LocalRecord | undefined>
    selectedSubForm$: Observable<FormIntf | undefined>

    records$(databaseName: string, formName: string): Observable<Record[]>

    forms$(databaseId: string, folderId: string): Observable<FormDefinition[]>

    form$(databaseId: string, id: string): Observable<FormDefinition | undefined>
}

export function newStore(reducer: ActionReducer<State>): Store<State> {

    const actions = new ReplaySubject<Action>()
    const nextActions = new Subject<Action>()
    const state$ = new BehaviorSubject<State>(initialState)

    actions.subscribe({
        next: (a) => {
            const currentState = state$.value;
            const nextState = reducer(currentState, a)
            if (nextState !== currentState) {
                state$.next(nextState)
                nextActions.next(a)
            }
            console.log(a)
            console.log(nextState)
        }
    })

    const s: Store<State> = {
        actions$: nextActions,
        state$: state$,
        dispatch(action: Action) {
            actions.next(action)
        },
        selectedDatabase$: state$.pipe(
            map(s => {
                for (let database of s.databases) {
                    if (database.id === s.selectedDatabaseId) {
                        return database
                    }
                }
                return undefined
            }),
            shareReplay()
        ),
        selectedForm$: state$.pipe(
            map(s => {
                if (s.formIds.hasOwnProperty(s.selectedFormId)) {
                    return s.forms[s.formIds[s.selectedFormId]]
                }
                return undefined
            })
        ),
        selectedDatabaseForms$: state$.pipe(
            map(s => {
                const result: FormDefinition[] = []
                for (let form of s.forms) {
                    if (form.databaseId === s.selectedDatabaseId) {
                        result.push(form)
                    }
                }
                return result
            }),
        ),
        createFormName$: state$.pipe(
            map(s => s.form.forms[s.form.selectedFormIdx]),
            map(s => s.name),
            distinctUntilChanged(),
            shareReplay()),
        createFormFields$: state$.pipe(
            map(s => s.form.fields.filter(f => f.formIdx === s.form.selectedFormIdx)),
            distinctUntilChanged(),
            shareReplay()),
        createFormAllFields$: state$.pipe(
            map(s => s.form.fields),
            distinctUntilChanged(),
            shareReplay()),
        createFormForms$: state$.pipe(
            map(s => s.form.forms),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormSelectedFormIdx$: state$.pipe(
            map(s => s.form.selectedFormIdx),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormSelectedForm$: state$.pipe(
            map(s => s.form.forms[s.form.selectedFormIdx]),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormSelectedFormIsSubForm$: state$.pipe(
            map(s => s.form.selectedFormIdx > 0),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormParentFormIdx$: state$.pipe(
            map(s => {
                if (s.form.selectedFormIdx < 0) {
                    return -1
                }
                return s.form.forms[s.form.selectedFormIdx].parentFormIdx
            }),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormParentForm$: state$.pipe(
            map(s => {
                if (s.form.selectedFormIdx < 0) {
                    return undefined
                }
                const form = s.form.forms[s.form.selectedFormIdx]
                if (form.parentFormIdx < 0) {
                    return undefined
                }
                return s.form.forms[form.parentFormIdx]
            }),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormSelectedFieldIdx$: state$.pipe(
            map(s => s.form.selectedFieldIdx),
            distinctUntilChanged(),
            shareReplay(),
        ),
        createFormSelectedField$: state$.pipe(
            map(s => {
                if (s.form.selectedFieldIdx === -1) {
                    return undefined
                }
                return s.form.fields[s.form.selectedFieldIdx]
            }),
            distinctUntilChanged(),
            shareReplay(),
        ),
        postRecordSuccess$: state$.pipe(
            map(s => s.createRecordSuccess),
            distinctUntilChanged(),
            shareReplay()),
        selectedRecord: state$.pipe(
            map(s => {
                if (s.selectedRecordId) {
                    const recordIdx = s.recordIds[s.selectedRecordId]
                    return s.records[recordIdx]
                } else {
                    return undefined
                }
            })
        ),
        selectedSubForm$: state$.pipe(
            map(s => {
                if (!s.selectedSubFormId && !s.selectedFormId) {
                    return
                }
                if (s.selectedFormId && !s.selectedSubFormId) {
                    return s.forms[s.formIds[s.selectedFormId]]
                }

                const findSubForm = (fields: FieldDefinition[], id: string): FormIntf | undefined => {
                    for (let i = 0; i < fields.length; i++) {
                        const field = fields[i]
                        if (field.fieldType.subForm) {
                            if (field.fieldType.subForm.id === id) {
                                return field.fieldType.subForm
                            }
                            const childSf = findSubForm(field.fieldType.subForm.fields, id)
                            if (childSf) {
                                return childSf
                            }
                        }
                    }
                }

                return findSubForm(s.forms[s.formIds[s.selectedFormId]].fields, s.selectedSubFormId)

            })
        ),
        forms$(databaseId: string, folderId: string): Observable<FormDefinition[]> {
            return state$.pipe(
                map(s => s.forms),
                map(forms => {
                    return forms.filter(r => {
                        if (folderId) {
                            return r.databaseId === databaseId && r.folderId === folderId
                        } else {
                            return r.databaseId === databaseId && !r.folderId
                        }
                    });
                }),
                shareReplay())
        },
        form$(databaseId: string, id: string): Observable<FormDefinition | undefined> {
            return state$.pipe(
                map(s => s.forms.find(f => f.id === id && f.databaseId === databaseId)),
            )
        },
        records$(databaseId: string, formId: string): Observable<LocalRecord[]> {
            return state$.pipe(
                map(s => s.records),
                filter(r => !!r),
                map(records => {
                    return (records as LocalRecord[]).filter(r => {
                        return r.formId === formId && r.databaseId === databaseId;
                    });
                }),
                distinctUntilChanged(),
                shareReplay())
        }
    }

    return s
}

export function setupFetchDBEffects(store: Store<State>, client: DatabaseLister): () => void {

    const subs: Subscription[] = []

    subs.push(store.actions$.pipe(
        ofType(fetchDatabases.type),
        switchMap(() => {
            return of(null).pipe(
                client.getDatabases(),
                switchMap(value => of(fetchDatabasesSuccess({databases: value.items}))),
                catchError(err => of(fetchDatabasesError({error: err}))))
        })
    ).subscribe(action => {
        store.dispatch(action)
    }))

    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupFetchFormsEffects(store: Store<State>, client: FormLister): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(fetchForms.type),
        switchMap(() => {
            return of(null).pipe(
                client.getForms(),
                switchMap(value => of(fetchFormsSuccess({forms: value.items}))),
                catchError(err => of(fetchFormsError({error: err}))))
        })
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupCreateDBEffects(store: Store<State>, client: DatabaseCreator): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(createDatabase),
        switchMap((action) => {
            return of(action).pipe(
                map(a => a.database),
                client.createDatabase(),
                switchMap(value => of(createDatabaseSuccess({database: value}))),
                catchError(err => of(createDatabaseError({error: err}))))
        })
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupCreateFormEffects(store: Store<State>, client: FormCreator): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(createForm),
        switchMap((action) => {
            return of(action).pipe(
                map(a => a.form),
                client.createForm(),
                switchMap(value => of(createFormSuccess({form: value}))),
                catchError(err => of(createFormError({error: err}))))
        })
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupCreateRecordEffect(store: Store<State>, client: RecordCreator): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(createRecord),
        switchMap((action) => {
            return of(action).pipe(
                map(a => a.record),
                client.postRecord(),
                switchMap(value => of(createRecordSuccess(value))),
                catchError(err => of(createRecordError({error: err}))))
        })
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupFetchRecordsEffect(store: Store<State>, client: RecordLister): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(fetchRecords),
        switchMap((s) =>
            of(s).pipe(
                map(s => s.options),
                client.fetchRecords(),
                switchMap(value => of(fetchRecordsSuccess({request: value.request, result: value.result.items}))),
                catchError(err => of(fetchRecordsError({error: err})))))
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupFetchFoldersEffect(store: Store<State>, client: FolderLister): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(fetchFolders),
        switchMap((s) =>
            of(s).pipe(
                map(s => s.options),
                client.listFolders(),
                switchMap(value => of(fetchFoldersSuccess({request: value.request, result: value.result.items}))),
                catchError(err => of(fetchFoldersError({error: err})))))
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}


export function setupCreateFolderEffects(store: Store<State>, client: FolderCreator): () => void {
    const subs: Subscription[] = []
    subs.push(store.actions$.pipe(
        ofType(createFolder),
        switchMap((action) => {
            return of(action).pipe(
                map(a => a.folder),
                client.createFolder(),
                switchMap(value => of(createFolderSuccess({folder: value}))),
                catchError(err => of(createFolderError({error: err}))))
        })
    ).subscribe(action => {
        store.dispatch(action)
    }))
    return () => {
        for (let sub of subs) {
            sub.unsubscribe()
        }
    }
}

export const store = newStore(reducer)
export const apiClient = new client()
setupFetchDBEffects(store, apiClient)
setupCreateDBEffects(store, apiClient)
setupFetchFormsEffects(store, apiClient)
setupCreateFormEffects(store, apiClient)
setupCreateRecordEffect(store, apiClient)
setupFetchRecordsEffect(store, apiClient)
setupFetchFoldersEffect(store, apiClient)
setupCreateFolderEffects(store, apiClient)



