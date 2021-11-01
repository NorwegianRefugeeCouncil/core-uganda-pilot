import {FC} from "react";
import {NavLink} from "react-router-dom"

export const NavBar: FC = props => {
    return (
        <nav className="navbar navbar-expand-sm navbar-light bg-light">
            <div className="container-fluid">
                <a className="navbar-brand" href="#">Core Admin</a>
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"
                        aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"/>
                </button>
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav">
                        <li className="nav-item">
                            <NavLink className={"nav-link"} activeClassName={"active"} to={"/organizations"}>Organizations</NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink className={"nav-link"} activeClassName={"active"} to={"/clients"}>Clients</NavLink>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    )
}
