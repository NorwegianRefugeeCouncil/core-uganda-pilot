import React, {FC, useState} from "react";
import {useForm} from "react-hook-form";
import {defaultClient} from "../../data/client";
import {databaseActions} from "../../reducers/database";
import {Database} from "../../types/types";
import {Redirect} from "react-router-dom"

type FormData = {
    name: string
}

export const DatabaseEditor: FC = props => {

    const {register, handleSubmit} = useForm<FormData>()
    const [database, setDatabase] = useState<Database | undefined>(undefined)

    const onSubmit = (data: FormData) => {
        defaultClient.createDatabase({object: {name: data.name}}).then(resp => {
            if (resp.response) {
                databaseActions.addOne(resp.response)
                setDatabase(resp.response)
            } else {

            }
        })
    }
    if (database) {
        return <Redirect to={`/browse/databases/${database.id}`}/>
    }
    return <div className={"flex-grow-1 bg-dark text-white pt-3"}>
        <div className={"container"}>
            <div className={"row"}>
                <div className={"col"}>
                    <h3>Create New Database</h3>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <div className={"form-group mb-2"}>
                            <label htmlFor={"name"}>Database Name</label>
                            <input className={"form-control"} {...register("name")}/>
                        </div>
                        <button className={"btn btn-primary"}>Create New Database</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
}
