import React, {Fragment, FC, useState} from "react";
import {useForm} from "react-hook-form";
import {databaseActions} from "../../reducers/database";
import {Folder} from "../../types/types";
import {Redirect} from "react-router-dom"
import {useDatabaseFromQueryParam, useFolderFromQueryParam} from "../../app/hooks";
import client from "../../app/client";

type FormData = {
    name: string
}

export const FolderEditor: FC = props => {

    const {register, handleSubmit,} = useForm<FormData>()

    const parentFolder = useFolderFromQueryParam("parentId")
    const database = useDatabaseFromQueryParam("databaseId")
    const [folder, setFolder] = useState<Folder | undefined>(undefined)

    const onSubmit = (data: FormData) => {
        if (!database?.id) {
            return
        }
        client.createFolder({
            object: {
                name: data.name,
                databaseId: database?.id,
                parentId: parentFolder?.id
            }
        }).then(resp => {
            if (resp.response) {
                databaseActions.addOne(resp.response)
                setFolder(resp.response)
            } else {

            }
        })
    }

    if (!database) {
        return <Fragment>Database not found</Fragment>
    }

    if (folder) {
        return <Redirect to={`/browse/folders/${folder.id}`}/>
    }

    return <div className={"flex-grow-1 bg-dark text-white pt-3"}>
        <div className={"container"}>
            <div className={"row"}>
                <div className={"col"}>
                    <h3>Create New Folder</h3>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <div className={"form-group mb-2"}>
                            <label htmlFor={"name"}>Folder Name</label>
                            <input className={"form-control"} {...register("name")}/>
                        </div>
                        <button
                            className={"btn btn-primary"}>Create New Folder</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
}
