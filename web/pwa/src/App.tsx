import React, {useEffect} from 'react';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import * as SQL from 'sql.js'
import {NavBarContainer} from "./features/navbar/navbar";
import {BrowserRouter, Switch, Route, Redirect} from "react-router-dom";
import {useAppDispatch} from "./app/hooks";
import {fetchDatabases} from "./reducers/database";
import {fetchForms} from "./reducers/form";
import {fetchFolders} from "./reducers/folder";
import {FolderBrowserContainer} from "./features/browser/FolderBrowser";
import {FormBrowserContainer} from "./features/browser/FormBrowser";
import {DatabasesContainer} from "./features/browser/Databases";
import {RecordEditorContainer} from "./features/recorder/RecordEditor";
import {FormerContainer} from "./features/former/Former";
import {AdminPanel} from "./features/admin/AdminPanel";
import {useAuth} from 'oidc-react';
import {DatabaseEditor} from "./features/databases/DatabaseEditor";
import {Database, SqlJsStatic} from "sql.js";
import {inflate, deflate} from "pako"
import {FolderEditor} from "./features/folders/FolderEditor";
import log from "loglevel"
import {RecordBrowser} from "./features/browser/RecordBrowser";

function App() {

    useEffect(() => {
        SQL.default({
            locateFile(url: string): string {
                log.debug("locating sql file at url", url)
                return `/${url}`
            }
        }).then((sql) => {

            const db = getDb(sql)

            try {
                db.run(`CREATE TABLE IF NOT EXISTS users  (id int primary key);`)
                const stmt = db.prepare("select * from users")
                stmt.bind({0: "id"})
                while (stmt.step()) {
                    console.log(stmt.getAsObject())
                }
            } catch (err) {
                console.log(err.toString())
            }

            saveDb(db)

        })
    }, [])


    const dispatch = useAppDispatch()

    useEffect(() => {
        dispatch(fetchDatabases())
        dispatch(fetchForms())
        dispatch(fetchFolders())
    }, [dispatch])

    const auth = useAuth()

    if (auth.isLoading) {
        return <div>Loading...</div>
    }

    return (
        <BrowserRouter>
            <div className="App">
                <div
                    style={{maxHeight: "100vh", maxWidth: "100vw"}}
                    className={"d-flex flex-column wh-100 vh-100"}>

                    <NavBarContainer/>
                    <Switch>

                        <Route path={`/edit/forms/:formId/record`} render={p => {
                            return <RecordEditorContainer/>
                        }}/>

                        <Route path={`/edit/forms`} render={p => {
                            return <FormerContainer/>
                        }}/>

                        <Route path={`/add/folders`} render={p => {
                            return <FolderEditor/>
                        }}/>

                        <Route path={`/edit/databases`} render={p => {
                            return <DatabaseEditor/>
                        }}/>

                        <Route path={`/browse/databases/:databaseId`} render={p => {
                            const {databaseId} = p.match.params
                            return <FolderBrowserContainer databaseId={databaseId}/>
                        }}/>

                        <Route path={`/browse/folders/:folderId`} render={p => {
                            const {folderId} = p.match.params
                            return <FolderBrowserContainer folderId={folderId}/>
                        }}/>

                        <Route path={`/browse/forms/:formId`} render={p => {
                            const search = new URLSearchParams(p.location.search)
                            const parentRecordId = search.get("parentRecordId")
                            const {formId} = p.match.params
                            return <FormBrowserContainer
                                parentRecordId={parentRecordId ? parentRecordId : ""}
                                formId={formId ? formId : ""}/>
                        }}/>

                        <Route path={`/browse/databases`}>
                            <DatabasesContainer/>
                        </Route>

                        <Route path={`/browse/records/:recordId`} component={RecordBrowser}/>

                        <Route path={`/admin`}>
                            <AdminPanel/>
                        </Route>

                        <Route path="/">
                            <Redirect to="/browse/databases"/>
                        </Route>
                    </Switch>

                </div>
            </div>


        </BrowserRouter>
    );
}

export default App;

/**
 * Convert an Uint8Array into a string.
 *
 * @returns {String}
 */
function DecodeUint8arr(uint8array: Uint8Array) {
    return uint8array.join(",")
}

/**
 * Convert a string into a Uint8Array.
 *
 * @returns {Uint8Array}
 */
function EncodeUint8arr(myString: string) {
    return Uint8Array.from(myString.split(","), v => parseInt(v))
}

function saveDb(db: Database) {
    const rawSqlBytes = db.export()
    const compressedSqlBytes = deflate(rawSqlBytes)
    const compressedSqlStr = DecodeUint8arr(compressedSqlBytes)
    localStorage.setItem("db", compressedSqlStr)
}


function getDb(sql: SqlJsStatic): Database {
    let db: Database
    const compressedStrFromStore = localStorage.getItem("db")
    if (compressedStrFromStore !== null) {
        console.log("loading file from local storage")
        const compressedBytesFromStore = EncodeUint8arr(compressedStrFromStore)
        const decompressedBytesFromStore = inflate(compressedBytesFromStore)
        db = new sql.Database(decompressedBytesFromStore)
    } else {
        console.log("creating new database")
        db = new sql.Database()
    }
    return db
}
