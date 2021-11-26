import { AuthRequestConfig, AuthSessionResult, CodeChallengeMethod, DiscoveryDocument, Prompt, ResponseType } from "../types/types";
export declare class AuthRequest {
    state: string;
    url: string;
    codeVerifier: string;
    codeChallenge: string;
    responseType: ResponseType;
    clientId: string;
    extraParams: Record<string, string>;
    usePKCE: boolean;
    codeChallengeMethod: CodeChallengeMethod;
    redirectUri: string;
    scopes: string[];
    prompt: Prompt;
    constructor(request: AuthRequestConfig);
    promptAsync(discoveryDocument: DiscoveryDocument): Promise<AuthSessionResult>;
    parseReturnUrl(url: string): AuthSessionResult;
    makeAuthUrlAsync(discovery: DiscoveryDocument): Promise<string>;
    private ensureCodeIsSetupAsync;
    private getAuthRequestConfigAsync;
}
//# sourceMappingURL=authrequest.d.ts.map