import React, {Fragment, useEffect, useRef, useState} from "react";
import {BehaviorSubject} from "rxjs";
import {Redirect} from "react-router-dom";
import {createDatabase, initialState, createDatabaseReset, createDatabaseSetName, State, store} from "./store";


export function CreateDatabase({databaseName = ""}) {

    const dbName$ = useRef(new BehaviorSubject(""))

    useEffect(() => {
        dbName$.current.next(databaseName);
    }, [databaseName]);

    const [state, setState] = useState<State>(initialState)

    useEffect(() => {
        store.dispatch(createDatabaseReset())
        const sub = store.state$.subscribe(s => {
            setState(s)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])

    return (
        <Fragment>
            <header className="App-header bg-white py-3">
                <div className="container">
                    <div className="row">
                        <div className="col">
                            <h2>Create Database</h2>
                        </div>
                    </div>
                </div>
            </header>
            <main>
                <div className="container">
                    <div className="row mt-3">
                        <div className="col">

                            {state.createDatabasePending ? <p>Pending</p> : <></>}

                            {state.createDatabaseError ?
                                <p>Error: {JSON.stringify(state.createDatabaseError)}</p> : <></>}

                            {state.createDatabaseSuccess ?
                                <Redirect to={`/databases/${state.createDatabaseId}`}/> : <></>}

                            <div className="form-group">
                                <label className={"form-label"} htmlFor={"databaseName"}>
                                    Enter database name:
                                </label>
                                <input
                                    className="form-control"
                                    type={"text"}
                                    value={state.createDatabaseName}
                                    id={"databaseName"}
                                    name={"databaseName"}
                                    onChange={event => {
                                        store.dispatch(createDatabaseSetName({name: event.target.value}))
                                    }}/>
                            </div>

                            <button className={"btn btn-primary mt-3"}
                                    disabled={state.createDatabaseName.length === 0}
                                    onClick={() => store.dispatch(createDatabase({database: {name: state.createDatabaseName}}))}>Save
                            </button>

                        </div>
                    </div>
                </div>
            </main>
        </Fragment>
    );
}
