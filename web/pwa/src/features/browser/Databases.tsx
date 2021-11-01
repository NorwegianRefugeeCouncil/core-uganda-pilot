import React, {FC, Fragment, useEffect} from "react";
import {Database} from "../../types/types";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {databaseGlobalSelectors, fetchDatabases,} from "../../reducers/database";
import {Link} from "react-router-dom";
import {DatabaseRow} from "./DatabaseRow";

export type DatabaseBrowserProps = {
    databases: Database[]
}

export const Databases: FC<DatabaseBrowserProps> = props => {


    const {databases} = props

    const menuEntries = databases.map(database => {
        return <DatabaseRow key={database.id} database={database}/>
    })

    return <Fragment>
        <div className={"py-3"}>
            <Link className={"btn btn-primary"} to={`/edit/databases`}>Add Database</Link>
        </div>
        <div className={"list-group shadow"} style={{cursor: "pointer"}}>
            {menuEntries}
        </div>
    </Fragment>
}

export type DatabaseBrowserContainerProps = {}

export const DatabasesContainer: FC<DatabaseBrowserContainerProps> = props => {

    const dispatch = useAppDispatch()

    useEffect(() => {
        dispatch(fetchDatabases())
    }, [dispatch])

    const databases = useAppSelector(databaseGlobalSelectors.selectAll)
    return <div className={"flex-grow-1 bg-light"}>
        <div className={"container"}>
            <div className={"row"}>
                <div className={"col"}>
                    <Databases databases={databases}/>
                </div>
            </div>
        </div>
    </div>
}
