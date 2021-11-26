import { WebBrowserAuthSessionResult } from "../types/types";
export declare function maybeCompleteAuthSession(): {
    type: string;
    message: string;
};
export declare function openAuthSessionAsync(args: {
    url: string;
}): Promise<WebBrowserAuthSessionResult>;
export declare function generateRandom(size: number): string;
//# sourceMappingURL=browser.d.ts.map