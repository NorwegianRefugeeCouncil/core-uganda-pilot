import {FC, useCallback} from "react";
import {Database} from "../types/types";
import {useDatabases} from "../app/hooks";

type DatabasePickerProps = {
    databaseId: string | undefined
    databases: Database[]
    setDatabaseId: (databaseId: string) => void
}

export const DatabasePicker: FC<DatabasePickerProps> = props => {
    const {databases, databaseId, setDatabaseId} = props
    return <div>
        <select
            placeholder={"Select Database"}
            onChange={e => setDatabaseId(e.target.value)}
            value={databaseId ? databaseId : ""}
            className="form-select" aria-label="Select Database">
            <option disabled={true} value={""}>Select Database</option>
            {databases.map(d => {
                return (
                    <option
                        value={d.id}>{d.name}
                    </option>
                );
            })}
        </select>
    </div>

}

type DatabasePickerContainerProps = {
    databaseId: string | undefined
    setDatabaseId?: (databaseId: string) => void
    setDatabase?: (database: Database | undefined) => void
}

const DatabasePickerContainer: FC<DatabasePickerContainerProps> = props => {
    const databases = useDatabases()

    const {databaseId, setDatabaseId, setDatabase} = props

    const setDbCallback = useCallback((databaseId: string) => {
        if (setDatabaseId) {
            setDatabaseId(databaseId)
        }
        const database = databases.find(d => d.id === databaseId)
        if (setDatabase) {
            setDatabase(database)
        }
    }, [databases, setDatabase, setDatabaseId])

    return <DatabasePicker
        databaseId={databaseId}
        setDatabaseId={setDbCallback}
        databases={databases}
    />
}

export default DatabasePickerContainer
