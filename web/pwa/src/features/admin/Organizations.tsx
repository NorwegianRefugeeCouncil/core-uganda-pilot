import React, {FC, useEffect} from "react";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {fetchOrganizations, organizationGlobalSelectors} from "../../reducers/organizations";
import {Link, useRouteMatch} from "react-router-dom";

export const Organizations: FC = props => {

    const dispatch = useAppDispatch()
    const match = useRouteMatch()
    useEffect(() => {
        dispatch(fetchOrganizations())
    }, [dispatch])

    const organizations = useAppSelector(organizationGlobalSelectors.selectAll)

    return (
        <div className={"d-flex flex-row w-100 h-100 flex-grow-1"}>
            <ul className={"list-group h-100 border-end pt-2 rounded-0"}>
                <li className={"list-group-item border-0 rounded-0"}>
                    <Link
                        className={"btn btn-light"}
                        to={`${match.url}/organizations`}>Organizations</Link>
                </li>
            </ul>
            <div className={"flex-grow-1"}>
                <div className={"p-2"}>
                    <Link to={`${match.url}/add`} className={"btn btn-light"}>Add Organization</Link>
                </div>
                <table className={"table w-100 rounded-0"}>
                    <thead>
                    <tr>
                        <th>Key</th>
                        <th>Name</th>
                    </tr>
                    </thead>
                    <tbody>
                    {organizations.map(o => {
                        return <tr>
                            <td><Link to={`${match.url}/${o.id}`}>{o.key}</Link></td>
                            <td>{o.name}</td>
                        </tr>
                    })}
                    </tbody>
                </table>
            </div>
        </div>

    )

}
