import React, {Fragment, useEffect, useState} from "react";
import {Link} from "react-router-dom";
import {fetchDatabases, fetchForms, initialState, State, store} from "./store";
import {useDatabases} from "./utils";

export function Databases() {

    useEffect(() => {
        store.dispatch(fetchDatabases())
        store.dispatch(fetchForms())
    }, [])

    const databases = useDatabases()

    return (
        <Fragment>
            <header className="App-header bg-white py-3">
                <div className="container">
                    <div className="row">
                        <div className="col">
                            <h2>Databases</h2>
                        </div>
                    </div>
                </div>
            </header>
            <main className="bg-light">
                <div className="container">
                    <div className="row">
                        <div className="col">
                            <Link to="/databases/new" className="btn btn-primary my-3">
                                Create database
                            </Link>
                        </div>
                    </div>
                    <div className="row">
                        <div className="col">

                            <ul className="list-group">
                                {databases?.map(d =>
                                    <Link
                                        key={d.id}
                                        className="list-group-item py-4 text-black text-decoration-none fw-bold"
                                        to={`/databases/${d.id}`}>
                                        <i className="bi bi-box"/> {d.name}
                                    </Link>)}
                            </ul>


                        </div>
                    </div>
                </div>
            </main>
        </Fragment>
    );

}
