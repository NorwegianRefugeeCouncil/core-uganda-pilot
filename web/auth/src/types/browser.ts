import {listenerMapEntry, WebBrowserAuthSessionResult, WebBrowserResultType} from "./types";
import {handler} from "../utils/helpers";


export default class Browser {
    private readonly CHARSET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    private listenerMap: Map<Window, listenerMapEntry>;
    private popupWindow: Window | null;
    private readonly handle: string;
    private webBrowserAuthSessionResult: WebBrowserAuthSessionResult;

    constructor() {
        this.popupWindow = null
        this.listenerMap = new Map<Window, listenerMapEntry>()
        this.handle = "NrcCoreBrowserRedirectHandle"
        this.webBrowserAuthSessionResult = {type: WebBrowserResultType.DISMISS}
    }

    getPopup() {
        return this.popupWindow;
    }

    getListenerMap() {
        return this.listenerMap;
    }
    getHandle() {
        return this.handle;
    }

    maybeCompleteAuthSession() {
        const handle = window.sessionStorage.getItem(this.handle)

        if (!handle) {
            return {type: "failed", message: "No auth session is currently in progress"}
        }
        const url = window.location.href

        /** TODO
         * if (skipRedirectCheck !== true) {
              const redirectUrl = window.localStorage.getItem(getRedirectUrlHandle(handle));
              // Compare the original redirect url against the current url with it's query params removed.
              const currentUrl = window.location.origin + window.location.pathname;
              if (!compareUrls(redirectUrl, currentUrl)) {
                return {
                  type: 'failed',
                  message: `Current URL "${currentUrl}" and original redirect URL "${redirectUrl}" do not match.`,
                };
              }
            }
         */

        const parent = window.opener ?? window.parent
        if (!parent) {
            throw new Error("ERR_WEB_BROWSER_REDIRECT: The window cannot complete the redirect request because the invoking window doesn't have a reference to its parent")
        }
        parent.postMessage({url, sender: handle}, url)
        return {type: "success", message: "attempting to complete auth"}
    }

    dismissPopup = () => {
        if (!this.popupWindow) {
            return
        }
        this.popupWindow.close()
        if (this.listenerMap.has(this.popupWindow)) {
            const {listener, interval} = this.listenerMap.get(this.popupWindow) as listenerMapEntry
            clearInterval(interval)
            window.removeEventListener('message', listener)
            this.listenerMap.delete(this.popupWindow)
            const handle = window.sessionStorage.getItem(this.handle)
            if (handle) {
                window.sessionStorage.removeItem(this.handle)
            }
        }
        this.popupWindow = null
    }

    createPopup(url: string) {
        const popupWidth = 380
        const popupHeight = 480
        const left = window.screen.width / 2 - popupWidth / 2
        const top = window.screen.height / 2 - popupHeight / 2
        const popupFeaturesString =
            `menubar=no,location=no,resizable=no,scrollbars=no,status=no,`+
            `width=${popupWidth},height=${popupHeight},top=${top},left=${left}`;
        this.popupWindow = window.open(url, "Core Login", popupFeaturesString)
        if (this.popupWindow) {
            try {
                this.popupWindow.focus()
            } catch (e) {
            }
        } else {
            throw new Error("ERR_WEB_BROWSER_LOCKED: Popup window was blocked by the browser of failed to open.")
        }
    }

    async openAuthSessionAsync(args: { url: string }) {
        const {url} = args
        const state = await this.getStateFromUrlOrGenerateAsync(url)

        window.sessionStorage.setItem(this.handle, state)

        const popup= this.getPopup();
        if (popup == null || popup.closed) {
            this.createPopup(url)
        }

        return new Promise<WebBrowserAuthSessionResult>((resolve) => handler(resolve, this));
    }

    async generateStateAsync(): Promise<string> {
        if (!this.isSubtleCryptoAvailable()) {
            throw new Error('ERR_WEB_BROWSER_CRYPTO: The current environment does not support Crypto',);
        }
        const encoder = new TextEncoder();
        const data = this.generateRandom(10);
        const buffer = encoder.encode(data);
        const hashedData = await window.crypto.subtle.digest('SHA-256', buffer);
        return btoa(String.fromCharCode(...new Uint8Array(hashedData)));
    }

    isCryptoAvailable(): boolean {
        return !!(window?.crypto as any);
    }

    isSubtleCryptoAvailable(): boolean {
        if (!this.isCryptoAvailable()) return false;
        return !!(window.crypto.subtle as any);
    }

    bufferToString(buffer: Uint8Array): string {
        const state: string[] = [];
        for (let i = 0; i < buffer.byteLength; i += 1) {
            const index = buffer[i] % this.CHARSET.length;
            state.push((this.CHARSET)[index]);
        }
        return state.join('');
    }

    generateRandom(size: number): string {
        let arr = new Uint8Array(size);
        if (arr.byteLength !== arr.length) {
            arr = new Uint8Array(arr.buffer);
        }
        const array = new Uint8Array(arr.length);
        if (this.isCryptoAvailable()) {
            window.crypto.getRandomValues(array);
        } else {
            for (let i = 0; i < size; i += 1) {
                array[i] = (Math.random() * this.CHARSET.length) | 0;
            }
        }
        return this.bufferToString(array);
    }

    async getStateFromUrlOrGenerateAsync(inputUrl: string): Promise<string> {
        const url = new URL(inputUrl);
        if (url.searchParams.has('state') && typeof url.searchParams.get('state') === 'string') {
            return url.searchParams.get('state')!;
        }
        return this.generateStateAsync();
    }

}
