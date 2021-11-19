import SQL, {Database, SqlJsStatic} from "sql.js";
import {deflate, inflate} from "pako";
import {createContext, FC, Fragment, useContext, useEffect, useState} from "react";
import log from "loglevel";

export type sqlContext = {
    sqlStatic: SqlJsStatic | undefined
    pending: boolean
    error: any
    useDatabase: () => dbProps
}

export type dbProps = {
    database: Database
    saveDatabase: () => void
}

export const SQLContext = createContext<sqlContext>({
    sqlStatic: undefined,
    pending: false,
    error: undefined,
    useDatabase: () => {
        throw new Error("sql static not initialized")
    }
})

export const SQLContextProvider: FC = props => {

    const [sqlStatic, setSqlStatic] = useState<SqlJsStatic | undefined>(undefined)
    const [pending, setPending] = useState(false)
    const [error, setError] = useState(undefined as any)

    useEffect(() => {
        log.debug("loading sql js context")
        setPending(true)
        SQL.default({
            locateFile(url: string): string {
                log.debug("locating sql.js resource at url", url)
                return `${process.env.PUBLIC_URL}/${url}`
            }
        }).then((sql) => {
            log.debug("loaded static sql context")
            setPending(false)
            setError(undefined)
            setSqlStatic(sql)
        }).catch(err => {
            setPending(false)
            setError(err)
        })
    }, [])

    return <SQLContext.Provider value={{
        sqlStatic: sqlStatic,
        error: error,
        pending: pending,
        useDatabase: () => {
            if (!sqlStatic) {
                throw new Error("sql static not initialized")
            }
            const db = getDb(sqlStatic)
            return {
                database: db,
                saveDatabase: () => {
                    saveDb(db)
                }
            }
        }
    }}>
        {sqlStatic ? props.children : <Fragment/>}
    </SQLContext.Provider>

}

export const useSQLContext = () => {
    return useContext(SQLContext)
}

export const useSqlDB = () => {
    const ctx = useSQLContext()
    return ctx.useDatabase()
}

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
