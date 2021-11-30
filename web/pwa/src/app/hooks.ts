import {TypedUseSelectorHook, useDispatch, useSelector} from 'react-redux';
import type {AppDispatch, RootState} from './store';
import {useLocation, useParams} from "react-router-dom";
import {useEffect, useMemo, useState} from "react";
import {folderGlobalSelectors} from "../reducers/folder";
import {databaseGlobalSelectors} from "../reducers/database";
import {recordGlobalSelectors, selectSubRecords, SubRecordResult} from "../reducers/records";
import {formGlobalSelectors, selectFormOrSubFormById} from "../reducers/form";
import {Database, FormDefinition, Record} from "../types/types";
import client, {ClientDefinition} from "core-js-api-client";

// Use throughout your app instead of plain `useDispatch` and `useSelector`
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

export const useQueryParam = (queryParamName: string) => {
    const location = useLocation()
    const [value, setValue] = useState<string | undefined>(undefined)
    useEffect(() => {
        const search = new URLSearchParams(location.search)
        const queryParamValue = search.get(queryParamName)
        if (value !== queryParamValue) {
            setValue(queryParamValue ? queryParamValue : undefined)
        }
    }, [queryParamName, value, location])
    return value
}


export const useFolderFromQueryParam = (queryParamName: string) => {
    const folderId = useQueryParam(queryParamName)
    return useAppSelector(state => {
        if (!folderId) {
            return undefined
        }
        return folderGlobalSelectors.selectById(state, folderId)
    })
}


export const useDatabaseFromQueryParam = (queryParamName: string) => {
    const databaseId = useQueryParam(queryParamName)
    return useAppSelector(state => {
        if (!databaseId) {
            return undefined
        }
        return databaseGlobalSelectors.selectById(state, databaseId)
    })
}

export const useRecordFromPath = (queryParamName: string) => {
    const params = useParams() as any
    const recordId = params[queryParamName] as string
    return useAppSelector(state => {
        if (!recordId) {
            return undefined
        }
        return recordGlobalSelectors.selectById(state, recordId)
    })
}

export const useFormOrSubForm = (formId: string | undefined) => {
    return useAppSelector(state => {
        if (!formId) {
            return undefined
        }
        return selectFormOrSubFormById(state, formId)
    })
}

export const useSubRecords: (recordId: string | undefined) => SubRecordResult | undefined = (recordId) => {
    return useAppSelector(state => {
        if (!recordId) {
            return undefined
        }
        return selectSubRecords(state, recordId)
    })
}

export const useParentRecord: (childRecordId: string | undefined) => Record | undefined = (recordId) => {
    return useAppSelector(state => {
        if (!recordId) {
            return undefined
        }
        const childRecord = recordGlobalSelectors.selectById(state, recordId)
        if (!childRecord) {
            return undefined
        }
        if (!childRecord.parentId) {
            return undefined
        }
        return recordGlobalSelectors.selectById(state, childRecord.parentId)
    })
}

export const useDatabases: () => Database[] = () => {
    return useAppSelector(databaseGlobalSelectors.selectAll)
}


export const useDatabase: (databaseId: string | undefined) => Database | undefined = (databaseId) => {
    return useAppSelector(state => {
        if (!databaseId) {
            return undefined
        }
        return databaseGlobalSelectors.selectById(state, databaseId)
    })
}


export const useForms: (options: { databaseId?: string | undefined }) => FormDefinition[] = ({databaseId}) => {
    return useAppSelector(state => {
        let allForms = formGlobalSelectors.selectAll(state)
        if (databaseId) {
            allForms = allForms.filter(f => f.databaseId === databaseId)
        }
        return allForms
    })
}


export function useApiClient(): ClientDefinition {
    return useMemo(() => {
        return new client()
    }, [])
}
