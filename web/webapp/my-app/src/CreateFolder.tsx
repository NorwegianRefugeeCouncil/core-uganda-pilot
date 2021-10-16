import React, {FC, Fragment, useEffect, useState} from "react";
import {Header} from "./Header";
import {useSelectedDatabase, useSelectedFolders} from "./utils";
import {Redirect} from "react-router-dom";
import {createFolder, createFolderReset, createFolderSetName, store} from "./store";

type CreateFolderProps = {}

export const CreateFolder: FC<CreateFolderProps> = props => {

    const database = useSelectedDatabase()
    const selectedFolders = useSelectedFolders()

    const [createFolderPending, setCreateFolderPending] = useState(false)
    const [createFolderSuccess, setCreateFolderSuccess] = useState(false)
    const [createFolderError, setCreateFolderError] = useState(undefined)
    const [createFolderId, setCreateFolderId] = useState("")
    const [createFolderName, setCreateFolderName] = useState("")

    useEffect(() => {
        const sub = store.state$.subscribe(s => {
            setCreateFolderPending(s.createFolderPending)
            setCreateFolderSuccess(s.createFolderSuccess)
            setCreateFolderError(s.createFolderError)
            setCreateFolderId(s.createFolderId)
            setCreateFolderName(s.createFolderName)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])

    useEffect(() => {
        if (createFolderSuccess) {
            store.dispatch(createFolderReset())
        }
    }, [createFolderSuccess])

    if (!database) {
        return <Fragment/>
    }

    const handleCreateFolder = () => {
        return store.dispatch(createFolder({
            folder: {
                name: createFolderName,
                databaseId: database?.id,
                parentId: selectedFolders.length > 0 ? selectedFolders[0].id : "",
            }
        }));
    }

    return <Fragment>
        <Header
            folders={selectedFolders}
            database={database}
            title={"Create Folder"}
        />
        <main>
            <div className="container">
                <div className="row mt-3">
                    <div className="col">

                        {createFolderPending ? <p>Pending</p> : <></>}

                        {createFolderError
                            ? <p>Error: {JSON.stringify(createFolderError)}</p>
                            : <></>
                        }

                        {
                            createFolderSuccess
                                ? <Redirect to={`/databases/${database?.id}?folderId=${createFolderId}`}/>
                                : <></>
                        }

                        <div className="form-group">
                            <label className={"form-label"} htmlFor={"folderName"}>
                                Enter Folder name:
                            </label>
                            <input
                                className="form-control"
                                type={"text"}
                                value={createFolderName}
                                id={"folderName"}
                                name={"folderName"}
                                onChange={event => {
                                    store.dispatch(createFolderSetName({
                                        name: event.target.value,
                                    }))
                                }}/>
                        </div>

                        <button className={"btn btn-primary mt-3"}
                                disabled={createFolderName.length === 0}
                                onClick={handleCreateFolder}>Save
                        </button>

                    </div>
                </div>
            </div>
        </main>
    </Fragment>

}


