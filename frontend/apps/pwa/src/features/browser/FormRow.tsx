import React, {FC} from "react";
import {FormDefinition} from "core-js-api-client";
import {Link} from "react-router-dom";

export const FormRow : FC<{form: FormDefinition}> = ({form}) => {

    return <Link
        to={`/browse/forms/${form.id}`}
        key={form.id}
        className={"list-group-item list-group-item-action py-4 fw-bold"}>
        <i className={"bi bi-file-earmark-text me-3"}/>
        <span>{form.name}</span>
    </Link>

}
