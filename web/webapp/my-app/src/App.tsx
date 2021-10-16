import React from 'react';
import './App.css';
import {BrowserRouter as Router, Link, Route, Switch} from "react-router-dom";
import {Databases} from "./Databases";
import {CreateDatabase} from "./CreateDatabase";
import {Database} from "./Database";
import {NavBar} from "./NavBar";
import {FormEditor} from "./FormEditor";
import {Form} from "./Form";
import {RecordEditorContainer} from "./RecordEditor";
import {FolderWatcher} from "./FolderWatcher";
import {CreateFolder} from "./CreateFolder";


export default function App() {
    return (
        <Router>
            <FolderWatcher/>
            <main>
                <nav>
                    <NavBar/>
                    <ul>
                        <li><Link to="/">Databases</Link></li>
                        <li><Link to="/forms">Forms</Link></li>
                    </ul>
                </nav>
                <Switch>
                    <Route path="/databases/:databaseId/folders/new" component={CreateFolder}/>
                    <Route path="/databases/:databaseId/forms/new" component={FormEditor}/>
                    <Route path="/databases/:databaseId/forms/:formId/add" component={RecordEditorContainer}/>
                    <Route path="/databases/:databaseId/forms/:formId" component={Form}/>
                    <Route path="/databases/new" component={CreateDatabase}/>
                    <Route path="/databases/:databaseId" component={Database}/>
                    <Route path="/" component={Databases}/>
                </Switch>
            </main>
        </Router>
    )
        ;
}
