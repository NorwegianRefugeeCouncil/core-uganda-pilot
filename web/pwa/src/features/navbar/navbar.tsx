import React, {FC, Fragment, useEffect, useState} from "react";
import {useAppSelector} from "../../app/hooks";
import {Link, NavLink, useLocation} from "react-router-dom"
import {Database, Folder} from "../../types/types";
import {databaseGlobalSelectors} from "../../reducers/database";
import {folderGlobalSelectors, selectParents} from "../../reducers/folder";
import {FormInterface, selectFormOrSubFormById, selectRootForm, selectSubFormParents} from "../../reducers/form";
import {getLogger} from "loglevel"

const log = getLogger("navbar")
log.setLevel(log.levels.TRACE)

export type BreadCrumb = {
    label: string
    link?: string
    icon?: string
}

export type Path = BreadCrumb[]


export type NavBarProps = {
    title: string
    database?: Database
    forms: FormInterface[]
    folders: Folder[]
    additionalItem: string
    secondaryTitle: string
}


function getPath(props: NavBarProps) {

    const {database, forms, folders, additionalItem} = props

    const path: Path = []

    path.push({
        label: "databases",
        link: `/browse/databases`,
        icon: "bi-box"
    })

    if (!database) {
        if (additionalItem) {
            path.push({
                label: additionalItem,
            })
        }
        return path
    }

    path.push({
        label: database.name,
        link: `/browse/databases/${database.id}`,
        icon: "bi-box"
    })

    for (let folder of folders) {
        path.push({
            label: folder.name,
            link: `/browse/folders/${folder.id}`,
        })
    }

    for (let form of forms) {
        path.push({
            label: form.name,
            link: `/browse/forms/${form.id}`,
        })
    }

    if (additionalItem) {
        path.push({
            label: additionalItem,
        })
    }

    return path
}

export const NavBar: FC<NavBarProps> = props => {

    const path = getPath(props)
    const {title, secondaryTitle} = props

    return (
        <Fragment>
            <nav className="navbar navbar-expand-sm navbar-light bg-light">
                <div className="container-fluid">
                    <Link className="navbar-brand" to="/">Core</Link>
                    <button className="navbar-toggler" type="button" data-bs-toggle="collapse"
                            data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent"
                            aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"/>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarSupportedContent">
                        <div className="d-flex flex-grow-1"/>
                        <ul className="navbar-nav me-auto mb-2 mb-lg-0">
                            <li className="nav-item">
                                <Link className="nav-link" to={"/browse/databases"}>
                                    <i className="bi bi-box"/> Databases
                                </Link>
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>
            <header className="bg-white pt-2 pb-0 shadow-sm">
                <div className="container-fluid">
                    <div className="row">
                        <div className="col">
                            <div className={"px-2"}>
                                <nav aria-label="breadcrumb">
                                    <small className={"text-muted"}>
                                        <ol className="breadcrumb mb-0">
                                            {path.map((b, i) => {
                                                return <li key={i}
                                                           className={"text-uppercase breadcrumb-item fw-bold small"}>
                                                    {b.link
                                                        ? <NavLink activeClassName={"active"}
                                                                   className={"text-decoration-none"}
                                                                   to={b.link}>
                                                            {b.label}
                                                        </NavLink>
                                                        : <span>{b.label}</span>}
                                                </li>
                                            })}
                                        </ol>
                                    </small>
                                </nav>
                                <h5 className={"mt-2"}>{title} {secondaryTitle ?
                                    <span className={"text-secondary"}>{secondaryTitle}</span> : <Fragment/>}</h5>

                            </div>
                        </div>
                    </div>
                </div>
            </header>
        </Fragment>
    )
}

const databasesRegex = new RegExp(/^\/browse\/databases$/)
const databaseRegex = new RegExp(/^\/browse\/databases\/([a-f0-9]{8}-[a-fA-F0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})$/)
const folderRegex = new RegExp(/^\/browse\/folders\/([a-f0-9]{8}-[a-fA-F0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})$/)
const formRegex = new RegExp(/^\/browse\/forms\/([a-f0-9]{8}-[a-fA-F0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})(?:\?parentRecordId=(?:[a-f0-9]{8}-[a-fA-F0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}))?$/)
const recordRegex = new RegExp(/^\/edit\/forms\/([a-f0-9]{8}-[a-fA-F0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})\/record$/)

export const NavBarContainer: React.FC = () => {

    const location = useLocation()
    const [databaseId, setDatabaseId] = useState<string>("")
    const [folderId, setFolderId] = useState<string>("")
    const [formId, setFormId] = useState("")
    const [title, setTitle] = useState("")
    const [additionalItem, setAdditionalItem] = useState("")
    const [secondaryTitle, setSecondaryTitle] = useState("")

    const database = useAppSelector(state => databaseGlobalSelectors.selectById(state, databaseId))
    const folder = useAppSelector(state => folderGlobalSelectors.selectById(state, folderId))
    const folderChain = useAppSelector(state => selectParents(state, folder?.id, true))
    const form = useAppSelector(state => selectFormOrSubFormById(state, formId))
    const formChain = useAppSelector(state => selectSubFormParents(state, form?.id, true))
    const rootForm = useAppSelector(state => {
        if (form) {
            return selectRootForm(state, form.id)
        }
    })

    useEffect(() => {
        const browseDatabasesMatch = location.pathname.match(databasesRegex)
        const browseDatabaseMatch = location.pathname.match(databaseRegex)
        const browseFolderMatch = location.pathname.match(folderRegex)
        const browseFormMatch = location.pathname.match(formRegex)
        const addRecordMatch = location.pathname.match(recordRegex)

        setAdditionalItem("")

        if (browseDatabasesMatch) {
            setDatabaseId("")
            setFormId("")
            setFolderId("")
            setSecondaryTitle("")
            setTitle("Databases")
            log.debug("route match: Browse Database List")
        } else if (browseDatabaseMatch) {
            const matchDatabaseId = browseDatabaseMatch[1];
            log.debug(`route match: Browse Database ${matchDatabaseId}`)

            setDatabaseId(matchDatabaseId)

            setFormId("")
            setFolderId("")
            setSecondaryTitle("Database")

            if (database) {
                setTitle(database.name)
                log.debug("route match: database found", database)
            }

        } else if (browseFolderMatch) {
            const matchFolderId = browseFolderMatch[1];
            log.debug(`route match: Browse Folder ${matchFolderId}`)

            setFolderId(matchFolderId)

            if (folder) {

                log.debug("route match: folder found", folder)

                setDatabaseId(folder.databaseId)
                setTitle(folder.name)
                setSecondaryTitle("Folder")

            }
        } else if (browseFormMatch) {
            let matchFormId = browseFormMatch[1];

            log.debug(`route match: Browse Form ${matchFormId}`)

            setFormId(matchFormId)

            if (form && rootForm) {

                log.debug("route match: root form found", rootForm)
                log.debug("route match: form found", form)

                setTitle(form.name)
                setDatabaseId(rootForm.databaseId)
                setFolderId(rootForm.folderId)
                setSecondaryTitle("Form")

            }

        } else if (addRecordMatch) {
            let matchFormId = addRecordMatch[1];

            log.debug(`route match: Add Record on form ${matchFormId}`)

            setFormId(matchFormId)
            setAdditionalItem("Add record")

            if (form && rootForm) {

                log.debug("route: root form found", rootForm)
                log.debug("route: form found", form)

                setFormId(form.id)
                setDatabaseId(rootForm.databaseId)
                setFolderId(rootForm.folderId)
                setAdditionalItem("Add Record")
            }

        }

    }, [database, folder, form, location, rootForm])

    return <NavBar
        database={database}
        folders={folderChain}
        forms={formChain}
        title={title}
        additionalItem={additionalItem}
        secondaryTitle={secondaryTitle}
    />
}
