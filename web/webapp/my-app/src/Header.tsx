import React, {FC, Fragment} from "react";
import {Database, Folder, FormDefinition} from "./client";
import {Link} from "react-router-dom";
import {FormIntf} from "./store";

type FormHeaderProps = {
    form?: FormIntf
    database: Database
    folders?: Folder[]
    title?: string
    additionalBreadcrumb?: string
}

type BreadCrumb = {
    label: string
    link?: string
}

export const Header: FC<FormHeaderProps> = props => {
    const {form, database, title, additionalBreadcrumb} = props
    const folders = props.folders ? props.folders : []

    const breadCrumbs: BreadCrumb[] = []

    breadCrumbs.push({
        link: `/databases`,
        label: "Databases"
    })

    breadCrumbs.push({
        link: `/databases/${database?.id}`,
        label: database?.name
    })

    for (let i = folders.length - 1; i >= 0; i--) {
        const folder = folders[i]
        breadCrumbs.push({
            link: `/databases/${database.id}?folderId=${folder.id}`,
            label: folder.name
        })
    }

    if (form) {
        let link = `/databases/${database.id}/forms/${form.id}`;
        if (folders.length > 0) {
            link = link + `?folderId=${folders[0].id}`
        }
        breadCrumbs.push({
            link: link,
            label: form.name
        })
    }

    if (additionalBreadcrumb) {
        breadCrumbs.push({label: additionalBreadcrumb})
    }

    return <Fragment>
        <header className="bg-white py-3 shadow-sm">
            <div className="container-fluid">
                <div className="row">
                    <div className="col">
                        <div className={"px-2"}>
                            <nav aria-label="breadcrumb">
                                <small className={"text-muted"}>
                                    <ol className="breadcrumb mb-0">
                                        <li className={"me-2"}><i className="bi bi-box"/></li>
                                        {breadCrumbs.map((b, i) => {
                                            return <li key={i} className={"text-uppercase breadcrumb-item fw-bold"}>
                                                {b.link
                                                    ? <Link to={b.link}>{b.label}</Link>
                                                    : <span>{b.label}</span>}
                                            </li>
                                        })}
                                    </ol>
                                </small>
                            </nav>
                            <h5 className={"mt-2"}>{title}</h5>
                        </div>
                    </div>
                </div>
            </div>
        </header>
    </Fragment>

}
