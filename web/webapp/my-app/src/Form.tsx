import React, {Fragment, useEffect} from "react";
import {useSelectedDatabase, useSelectedForm, useRecords, useSelectedFolders} from "./utils";
import {Header} from "./Header";
import {Link} from "react-router-dom"
import {fetchDatabases, fetchFolders, fetchForms, fetchRecords, store} from "./store";

export function Form() {

    const database = useSelectedDatabase()
    const selectedFolders = useSelectedFolders()
    const form = useSelectedForm()
    const records = useRecords(database?.id, form?.id)

    useEffect(() => {
        store.dispatch(fetchDatabases())
        store.dispatch(fetchForms())
        store.dispatch(fetchFolders({options: {}}))
    }, [])

    useEffect(() => {
        if (!database?.name) {
            return
        }
        if (!form?.name) {
            return
        }
        store.dispatch(fetchRecords({options: {databaseId: database.id, formId: form.id}}))
    }, [database, form])

    if (!database) {
        return <Fragment/>
    }

    return <Fragment>
        <Header form={form} database={database} folders={selectedFolders} title={form?.name}/>
        <div className={"p-3 d-flex flex-row"}>
            <Link
                to={`/databases/${database?.id}/forms/${form?.id}/add`}
                className={"btn btn-primary me-2"}>
                <i className={"bi bi-plus-circle-fill"}/> Add Record
            </Link>
        </div>
        <table className={"table table-bordered"}>
            <thead>
            <tr>
                {form?.fields.map(f => <th key={f.id}>{f.name}</th>)}
            </tr>
            </thead>
            <tbody>
            {records.map(r =>
                <tr key={r.id}>
                    {form?.fields.map(f => <td key={f.name}>{r.values[f.id]}</td>)}
                </tr>
            )}
            </tbody>
        </table>
    </Fragment>

}
