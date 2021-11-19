import React, {FC} from "react";

type FormNameProps = {
    formName: string
    setFormName: (name: string) => void
}

export const FormName: FC<FormNameProps> = props => {
    const {formName, setFormName} = props
    return <div>
        <div className={"form-group mb-2"}>
            <label className={"form-label"} htmlFor={"formName"}>Form Name</label>
            <input className={"form-control"}
                   id="formName"
                   type={"text"}
                   value={formName ? formName : ""}
                   onChange={event => setFormName(event.target.value)}/>
        </div>
    </div>
}
