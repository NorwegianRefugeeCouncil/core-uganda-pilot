import { AuthRequestConfig, AuthRequestPromptOptions, AuthSessionResult, DiscoveryDocument, IssuerOrDiscovery, PromptMethod } from "./types/types";
import { AuthRequest } from "./authrequest";
export declare function useDiscovery(issuerOrDiscovery: IssuerOrDiscovery): DiscoveryDocument | null;
export declare function useLoadedAuthRequest(config: AuthRequestConfig, discovery: DiscoveryDocument | null, AuthRequestInstance: typeof AuthRequest): AuthRequest | null;
export declare function useAuthRequestResult(request: AuthRequest | null, discovery: DiscoveryDocument | null): [AuthSessionResult | null, PromptMethod];
export declare function useAuthRequest(config: AuthRequestConfig, discovery: DiscoveryDocument | null): [
    AuthRequest | null,
    AuthSessionResult | null,
    (options?: AuthRequestPromptOptions) => Promise<AuthSessionResult>
];
