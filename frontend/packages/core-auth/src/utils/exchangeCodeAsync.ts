import {
    AccessTokenRequestConfig, DiscoveryDocument,
} from "../types/types";
import {AccessTokenRequest} from "../types/request";
import {TokenResponse} from "../types/response";

export default function exchangeCodeAsync(
    config: AccessTokenRequestConfig,
    discovery: Pick<DiscoveryDocument, 'token_endpoint'>
): Promise<TokenResponse> {
    const request = new AccessTokenRequest(config);
    return request.performAsync(discovery);
}
