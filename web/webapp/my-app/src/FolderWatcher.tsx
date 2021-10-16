import {useLocation} from "react-router-dom"
import {Fragment, useEffect} from "react";
import {setSelectedFolder, store} from "./store";

export function FolderWatcher() {
    const location = useLocation()
    useEffect(() => {
        const qry = new URLSearchParams(location.search)
        const folderId = qry.get("folderId")
        store.dispatch(setSelectedFolder({id: folderId ? folderId : undefined}))
    }, [location])
    return <Fragment/>
}
