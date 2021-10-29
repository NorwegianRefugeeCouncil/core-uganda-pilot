import React, {useEffect} from 'react';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import {NavBarContainer} from "./features/navbar/navbar";
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom";
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
import {FolderEditor} from "./features/folders/FolderEditor";
import {RecordBrowser} from "./features/browser/RecordBrowser";
import {useSqlDB} from "./app/db";

function App() {

    const dispatch = useAppDispatch()
    const {database, saveDatabase} = useSqlDB()

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

                        <Route path={`/edit/forms/:formId/record`} component={RecordEditorContainer}/>

                        <Route path={`/edit/forms`} component={FormerContainer}/>

                        <Route path={`/add/folders`} component={FolderEditor}/>

                        <Route path={`/edit/databases`} component={DatabaseEditor}/>

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

                        <Route path={`/browse/databases`} component={DatabasesContainer}/>

                        <Route path={`/browse/records/:recordId`} component={RecordBrowser}/>

                        <Route path={`/admin`} component={AdminPanel}/>

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

