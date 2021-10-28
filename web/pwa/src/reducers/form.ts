import {createAsyncThunk, createEntityAdapter, createSlice} from "@reduxjs/toolkit";
import {FieldDefinition, FieldTypeSubForm, FormDefinition} from "../types/types";
import {RootState} from "../app/store";
import {defaultClient, FormListResponse} from "../data/client";

const adapter = createEntityAdapter<FormDefinition>({
    // Assume IDs are stored in a field other than `book.id`
    selectId: (formDefinition) => formDefinition.id,
    // Keep the "all IDs" array sorted based on book titles
    sortComparer: (a, b) => a.name.localeCompare(b.name),
})

export const fetchForms = createAsyncThunk<FormListResponse>(
    'forms/fetch',
    async (_, thunkAPI) => {
        try {
            const response = await defaultClient.listForms({})
            if (response.success) {
                return response
            } else {
                return thunkAPI.rejectWithValue(response)
            }
        } catch (err) {
            return thunkAPI.rejectWithValue(err)
        }
    }
)


export const formsSlice = createSlice({
    name: 'forms',
    initialState: {
        ...adapter.getInitialState(),
        fetchPending: false,
        fetchError: undefined as any,
        fetchSuccess: true
    },
    reducers: {
        addOne: adapter.addOne,
        addMany: adapter.addMany,
        removeAll: adapter.removeAll,
        removeMany: adapter.removeMany,
        removeOne: adapter.removeOne,
        updateMany: adapter.updateMany,
        updateOne: adapter.updateOne,
        upsertOne: adapter.upsertOne,
        upsertMany: adapter.upsertMany,
        setOne: adapter.setOne,
        setMany: adapter.setMany,
        setAll: adapter.setAll,
    },
    extraReducers: builder => {
        builder.addCase(fetchForms.pending, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = true
        })
        builder.addCase(fetchForms.rejected, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = false
            state.fetchError = action.payload
        })
        builder.addCase(fetchForms.fulfilled, (state, action) => {
            state.fetchSuccess = true
            state.fetchPending = false
            state.fetchError = undefined
            if (action.payload.response?.items) {
                adapter.setAll(state, action.payload.response.items)
            }
        })
    }
})

export const formActions = formsSlice.actions
export const formSelectors = adapter.getSelectors()
export const formGlobalSelectors = adapter.getSelectors<RootState>(state => state.forms)

export type FormInterface = {
    id: string
    code: string
    isSubForm: boolean
    name: string
    fields: FieldDefinition[]
}


export const mapSubForm = (subForm: FieldTypeSubForm): FormInterface => {
    return {
        id: subForm.id,
        code: subForm.code,
        fields: subForm.fields,
        isSubForm: true,
        name: subForm.name,
    }
}


export const findSubForm = (id: string, fields: FieldDefinition[]): FormInterface | undefined => {
    if (!fields) {
        return undefined
    }
    for (let field of fields) {
        if (field.fieldType.subForm) {
            const subForm = field.fieldType.subForm
            if (subForm && subForm.id === id) {
                return mapSubForm(subForm)
            }
            let childSf = findSubForm(id, subForm.fields)
            if (childSf) {
                return childSf
            }
        }
    }
    return undefined
}

export const hasSubFormWithId = (id: string, fields: FieldDefinition[]): boolean => {
    if (!fields) {
        return false
    }
    for (let field of fields) {
        if (!field.fieldType.subForm) {
            continue
        }
        const subForm = field.fieldType.subForm
        if (subForm.id === id) {
            return true
        }
        const childHasSubForm = hasSubFormWithId(id, subForm.fields)
        if (childHasSubForm) {
            return true
        }
    }
    return false
}

export const selectRootForm = (state: RootState, formOrSubFormId: string | undefined): FormDefinition | undefined => {
    if (!formOrSubFormId){
        return undefined
    }
    const allForms = formGlobalSelectors.selectAll(state)
    for (let form of allForms) {
        if (form.id === formOrSubFormId) {
            return form
        }
        if (hasSubFormWithId(formOrSubFormId, form.fields)) {
            return form
        }
    }
    return undefined
}

interface FlatForms {
    parentMap: { [key: string]: string }
    rootMap: { [key: string]: string }
    idMap: { [key: string]: FormInterface }
    ids: string[]
}

function flattenSubForms(result: FlatForms, rootId: string, parentId: string, fields: FieldDefinition[]): void {
    if (!fields) {
        return
    }
    for (let field of fields) {
        const subForm = field.fieldType.subForm
        if (!subForm) {
            continue
        }
        result.idMap[subForm.id] = {
            ...subForm,
            isSubForm: true
        }
        result.parentMap[subForm.id] = parentId
        result.rootMap[subForm.id] = rootId
        result.ids.push(subForm.id)
        flattenSubForms(result, rootId, subForm.id, subForm.fields)
    }
}

const selectFlattenedForms = (state: RootState, rootFormId?: string): FlatForms => {

    const result: FlatForms = {
        idMap: {},
        rootMap: {},
        parentMap: {},
        ids: []
    }

    let allForms: FormDefinition[] = []
    if (rootFormId) {
        const f = formGlobalSelectors.selectById(state, rootFormId)
        if (!f) {
            return result
        }
        allForms = [f]
    } else {
        allForms = formGlobalSelectors.selectAll(state)
    }

    for (let form of allForms) {
        result.idMap[form.id] = {
            ...form,
            isSubForm: false
        }
        result.rootMap[form.id] = form.id
        result.ids.push(form.id)
        flattenSubForms(result, form.id, form.id, form.fields)
    }

    return result
}

export const selectSubFormParents = (state: RootState, subFormId?: string, includeSelf: boolean = false): FormInterface[] => {

    if (!subFormId) {
        return []
    }

    const flat = selectFlattenedForms(state)
    let result: FormInterface[] = []
    let walk = subFormId
    while (walk) {
        const parentId = flat.parentMap[walk]
        if (!parentId) {
            break
        }
        const parent = flat.idMap[parentId]
        if (!parent) {
            break
        }
        result.push(parent)
        walk = parent.id
    }

    if (includeSelf) {
        const self = flat.idMap[subFormId]
        result = [...result, self]
    }

    return result
}

export const selectFormOrSubFormById = (state: RootState, formOrSubFormId: string): FormInterface | undefined => {
    for (let f of formGlobalSelectors.selectAll(state)) {
        if (f.id === formOrSubFormId) {
            return {
                ...f,
                isSubForm: false
            }
        }
        const childSf = findSubForm(formOrSubFormId, f.fields)
        if (childSf) {
            return childSf
        }
    }
    return undefined
}

export const selectByFolderOrDBId = (state: RootState, folderOrDbId?: string): FormDefinition[] => {
    return formGlobalSelectors.selectAll(state).filter(f => {
        return (f.folderId === folderOrDbId || f.databaseId === folderOrDbId) && !f.folderId
    })
}

function findField(fieldId: string, fields: FieldDefinition[]): FieldDefinition | undefined {
    for (let field of fields) {
        if (field.id === fieldId) {
            return field
        }
        if (field.fieldType.subForm) {
            const sfField = findField(fieldId, field.fieldType.subForm.fields)
            if (sfField) {
                return sfField
            }
        }
    }
    return undefined
}


export const selectField = (fieldId: string): ((rootState: RootState) => FieldDefinition | undefined) => {
    return rootState => {
        const allForms = formGlobalSelectors.selectAll(rootState)
        let field: FieldDefinition | undefined
        for (let form of allForms) {
            field = findField(fieldId, form.fields)
            if (field) {
                break
            }
        }
        if (!field) {
            return undefined
        }
        return field
    }
}

export const selectSubFormForField = (fieldId: string): ((rootState: RootState) => FieldTypeSubForm | undefined) => {
    return rootState => {
        const field = selectField(fieldId)(rootState)
        if (!field) {
            return undefined
        }
        if (!field.fieldType.subForm) {
            return undefined
        }
        return field.fieldType.subForm
    }
}

export const findFieldForSubForm: (fields: FieldDefinition[], subFormId: string) => FieldDefinition | undefined = (fields, subFormId) => {
    for (const field of fields) {
        if (field.fieldType.subForm) {
            if (field.fieldType.subForm?.id === subFormId) {
                return field
            }
            let childField = findFieldForSubForm(field.fieldType.subForm?.fields, subFormId)
            if (childField) {
                return childField
            }
        }
    }
    return undefined
}

export const selectFieldForSubForm: (form: FormDefinition, subFormId: string) => FieldDefinition | undefined = (form, subFormId) => {
    return findFieldForSubForm(form.fields, subFormId)
}

export const selectSubFormFields: (rootState: RootState, formId: string) => FieldDefinition[] = (rootState, formId) => {

    const form = formGlobalSelectors.selectById(rootState, formId)
    if (!form) {
        return []
    }

    const result: FieldDefinition[] = []
    for (let field of form.fields) {
        if (field.fieldType.subForm) {
            result.push(field)
        }
    }

    return result
}

export default formsSlice.reducer

