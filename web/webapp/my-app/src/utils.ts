import {useParams, useLocation} from "react-router-dom";
import {useEffect, useState} from "react";
import {FormIntf, setSelectedDatabase, setSelectedForm, store} from "./store";
import {Database, Database as DB, Folder, FormDefinition, LocalRecord, Record} from "./client";

export function useSelectedDatabase() {
    const {databaseId} = useParams<{ databaseId: string }>()
    const [database, setDatabase] = useState<DB | undefined>(undefined)
    useEffect(() => {
        store.dispatch(setSelectedDatabase({id: databaseId}))
    }, [databaseId])
    useEffect(() => {
        const dbSub = store.selectedDatabase$.subscribe(s => setDatabase(s))
        return () => {
            dbSub.unsubscribe()
        }
    }, [])
    return database
}

export function useSelectedForm() {
    const {formId} = useParams<{ formId: string }>()
    const [form, setForm] = useState<FormDefinition | undefined>(undefined)
    useEffect(() => {
        if (form?.id !== formId ? formId : undefined) {
            store.dispatch(setSelectedForm({id: formId}))
        }
    }, [formId, form?.id])
    useEffect(() => {
        const sub = store.selectedForm$.subscribe(s => {
            if (form?.id !== s?.id) {
                setForm(s)
            }
        })
        return () => {
            sub.unsubscribe()
        }
    }, [form?.id])
    return form
}

export function useSelectedFormIntf() {
    const [subForm, setSubForm] = useState<FormIntf | undefined>(undefined)
    useEffect(() => {
        const sub = store.selectedSubForm$.subscribe(s => {
            if (subForm?.id !== s?.id) {
                setSubForm(s)
            }
        })
        return () => {
            sub.unsubscribe()
        }
    }, [subForm?.id])
    return subForm
}

export function useDatabaseForms(databaseId: string | undefined, folderId: string | undefined) {
    const [forms, setForms] = useState<FormDefinition[]>([])
    useEffect(() => {
        if (!databaseId) {
            setForms([])
            return
        }
        let fId = ""
        if (folderId) {
            fId = folderId
        }
        const dbSub = store.forms$(databaseId, fId).subscribe(s => {
            setForms(s)
        })
        return () => {
            dbSub.unsubscribe()
        }
    }, [databaseId, folderId])
    return forms
}

export function useDatabases() {
    const [databases, setDatabases] = useState<Database[]>([])
    useEffect(() => {
        const dbSub = store.state$.subscribe(s => {
            setDatabases(s.databases)
        })
        return () => {
            dbSub.unsubscribe()
        }
    }, [])
    return databases
}

export function useFolders(databaseId: string | undefined, parentFolderId: string | undefined) {

    const [folders, setFolders] = useState<Folder[]>([])

    useEffect(() => {
        const stateSub = store.state$.subscribe(s => {

            if (!databaseId) {
                setFolders([])
                return
            }

            const folderArr: Folder[] = []

            if (parentFolderId) {
                for (let folder of s.folders) {
                    if (folder.databaseId === databaseId && folder.parentId === parentFolderId) {
                        folderArr.push(folder)
                    }
                }
            } else {
                for (let folder of s.folders) {
                    if (folder.databaseId === databaseId && !folder.parentId) {
                        folderArr.push(folder)
                    }
                }
            }


            setFolders(folderArr)

        })
        return () => {
            stateSub.unsubscribe()
        }
    }, [databaseId, parentFolderId])
    return folders
}

export function useSelectedFolders() {
    const [folders, setFolders] = useState<Folder[]>([])
    useEffect(() => {
        const stateSub = store.state$.subscribe(s => {
            const folderArr: Folder[] = []
            let walk = s.selectedFolderId
            while (walk) {
                const folder = s.folders[s.folderIds[walk]]
                if (folder) {
                    folderArr.push(folder)
                    walk = folder.parentId
                } else {
                    break
                }
            }
            setFolders(folderArr)

        })
        return () => {
            stateSub.unsubscribe()
        }
    }, [])
    return folders
}

export function useRecords(databaseId?: string, formId?: string) {

    const [records, setRecords] = useState<Record[]>([])

    useEffect(() => {
        if (!databaseId) {
            return
        }
        if (!formId) {
            return
        }
        const dbSub = store.records$(databaseId, formId).subscribe(s => setRecords(s))
        return () => {
            dbSub.unsubscribe()
        }
    }, [databaseId, formId])

    return records
}


export function usePostRecordSuccess() {
    const [postRecordSuccess, setPostRecordSuccess] = useState<boolean>(false)
    useEffect(() => {
        const dbSub = store.postRecordSuccess$.subscribe(s => setPostRecordSuccess(s))
        return () => {
            dbSub.unsubscribe()
        }
    }, [])
    return postRecordSuccess
}


export function useSelectedRecord(): LocalRecord | undefined {
    const [selectedRecord, setSelectedRecord] = useState<LocalRecord | undefined>(undefined)
    useEffect(() => {
        const dbSub = store.selectedRecord.subscribe(s => setSelectedRecord(s))
        return () => {
            dbSub.unsubscribe()
        }
    }, [])
    return selectedRecord
}
