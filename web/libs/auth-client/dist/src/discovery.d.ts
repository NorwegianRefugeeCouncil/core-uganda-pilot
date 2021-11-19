import { DiscoveryDocument, IssuerOrDiscovery } from "./types/types";
export declare function resolveDiscoveryAsync(issuerOrDiscovery: IssuerOrDiscovery): Promise<DiscoveryDocument | null>;
