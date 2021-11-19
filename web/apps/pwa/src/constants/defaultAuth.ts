import {AuthWrapperProps} from "@core/auth-client";

const defaultAuth: AuthWrapperProps = {
    clientId: 'client-id',
    issuer: 'https://localhost/hydra',
    redirectUriSuffix: 'app',
    scopes: ['openid', 'profile', 'offline_access']
}
export default defaultAuth
