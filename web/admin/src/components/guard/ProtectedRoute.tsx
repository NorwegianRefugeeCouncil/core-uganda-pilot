import {ComponentType, FC, Fragment} from "react";
import {useIsAuthenticated} from "../../app/hooks";
import {Route} from "react-router-dom";

export type ProtectedRouteProps = {
    component: ComponentType
    path: string
}

export const ProtectedRoute: FC<ProtectedRouteProps> = props => {
    const {component, path} = props
    const isAuthenticated = useIsAuthenticated()
    if (!isAuthenticated) {
        return <Fragment/>
    }
    return <Route path={path} component={component}/>
}
