import axios from "axios";
import {DiscoveryDocument, IssuerOrDiscovery} from "../types/types";

export default async function resolveDiscoveryAsync(issuerOrDiscovery: IssuerOrDiscovery): Promise<DiscoveryDocument | null> {
    let issuer: string
    if (typeof issuerOrDiscovery === "string") {
        issuer = issuerOrDiscovery
    } else {
        issuer = issuerOrDiscovery.issuer
    }
    const metadataEndpoint = `${issuer}/.well-known/openid-configuration`;
    return axios.get<DiscoveryDocument>(metadataEndpoint)
        .then(value => value.data)
        .catch(err => {
            return null
        })
}
