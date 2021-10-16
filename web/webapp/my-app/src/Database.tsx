import React, {FC, Fragment, useEffect} from "react";
import {Link} from "react-router-dom";
import {useSelectedDatabase, useDatabaseForms, useFolders, useSelectedFolders} from "./utils";
import {fetchDatabases, fetchFolders, fetchForms, store} from "./store";
import {Header} from "./Header";

enum EntryType {
    Form = "form",
    Folder = "folder"
}

type MenuEntry = {
    type: EntryType
    label: string
    link: string
    id: string
}

export function Database() {

    const database = useSelectedDatabase()
    const selectedFolders = useSelectedFolders()

    const selectedFolderId = selectedFolders.length > 0
        ? selectedFolders[0].id
        : undefined

    const title = selectedFolders.length > 0
        ? selectedFolders[0].name
        : database?.name

    const createFolderLink = selectedFolderId
        ? `/databases/${database?.id}/folders/new?folderId=${selectedFolderId}`
        : `/databases/${database?.id}/folders/new`

    const createFormLink = selectedFolderId
        ? `/databases/${database?.id}/forms/new?folderId=${selectedFolderId}`
        : `/databases/${database?.id}/forms/new`

    const forms = useDatabaseForms(database?.id, selectedFolderId)

    const folders = useFolders(database?.id, selectedFolderId)

    useEffect(() => {
        store.dispatch(fetchDatabases())
        store.dispatch(fetchForms())
        store.dispatch(fetchFolders({options: {}}))
    }, [])

    const formMenuEntries: MenuEntry[] = forms.map(f => ({
        type: EntryType.Form,
        label: f.name,
        id: f.id,
        link: f.folderId
            ? `/databases/${f.databaseId}/forms/${f.id}?folderId=${f.folderId}`
            : `/databases/${f.databaseId}/forms/${f.id}`
    }))

    const folderEntries: MenuEntry[] = folders.map(f => ({
        type: EntryType.Folder,
        label: f.name,
        id: f.id,
        link: `/databases/${f.databaseId}?folderId=${f.id}`
    }))

    const allEntries = [...formMenuEntries, ...folderEntries]

    if (!database) {
        return <Fragment/>
    }

    return (
        <Fragment>
            <Header database={database} title={title} folders={selectedFolders}/>
            <main>
                <div className={"container mt-3"}>
                    <div className={"row"}>
                        <div className={"col"}>
                            <div className={'d-flex flex-row'}>
                                <Link
                                    className={"btn btn-primary"}
                                    to={createFormLink}>
                                    <i className={"bi bi-plus-circle-fill"}/> Add Form
                                </Link>
                                <Link
                                    className={"btn btn-primary ms-2"}
                                    to={createFolderLink}>
                                    <i className={"bi bi-plus-circle-fill"}/> Add Folder
                                </Link>
                            </div>
                        </div>
                    </div>
                    <div className={"row mt-2"}>
                        <div className={"col"}>
                            <ul className={"list-group"}>
                                {allEntries.map(f => <FormEntry key={f.id} entry={f}/>)}
                            </ul>
                        </div>
                    </div>
                </div>
            </main>
        </Fragment>
    );
}

export const FormEntry: FC<{ entry: MenuEntry }> = props => {
    const {entry} = props
    return <Link
        to={entry.link}
        className={"list-group-item p-3"}>
        <span>{entry.label}</span>
    </Link>
}
