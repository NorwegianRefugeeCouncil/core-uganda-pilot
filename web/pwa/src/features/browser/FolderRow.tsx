import React, {FC} from "react";
import {Folder} from "core-js-api-client";
import {Link} from "react-router-dom";

export const FolderRow : FC<{folder: Folder}> = ({folder}) => {
    return <Link
        to={`/browse/folders/${folder.id}`}
        key={folder.id}
        className={"list-group-item list-group-item-action py-4 fw-bold"}>
        <i className={"bi bi-folder me-3"}/>
        <span>{folder.name}</span>
    </Link>
}
