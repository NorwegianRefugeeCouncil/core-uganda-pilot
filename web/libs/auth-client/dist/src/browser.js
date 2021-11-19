"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateRandom = exports.openAuthSessionAsync = exports.maybeCompleteAuthSession = void 0;
const types_1 = require("./types/types");
let popupWindow;
const listenerMap = new Map();
const getHandle = () => "NrcCoreBrowserRedirectHandle";
function maybeCompleteAuthSession() {
    var _a;
    const handle = window.sessionStorage.getItem(getHandle());
    if (!handle) {
        console.log("NO HANDLE");
        return { type: "failed", message: "No auth session is currently in progress" };
    }
    const url = window.location.href;
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
    const parent = (_a = window.opener) !== null && _a !== void 0 ? _a : window.parent;
    if (!parent) {
        throw new Error("ERR_WEB_BROWSER_REDIRECT: The window cannot complete the redirect request because the invoking window doesn't have a reference to its parent");
    }
    parent.postMessage({ url, sender: handle });
    return { type: "success", message: "attempting to complete auth" };
}
exports.maybeCompleteAuthSession = maybeCompleteAuthSession;
const dismissPopup = () => {
    console.log("DISMISSING POPUP");
    if (!popupWindow) {
        console.log("CANNOT DISMISS POPUP: no popup window");
        return;
    }
    popupWindow.close();
    if (listenerMap.has(popupWindow)) {
        const { listener, interval } = listenerMap.get(popupWindow);
        clearInterval(interval);
        window.removeEventListener('message', listener);
        listenerMap.delete(popupWindow);
        const handle = window.sessionStorage.getItem(getHandle());
        if (handle) {
            window.sessionStorage.removeItem(getHandle());
        }
    }
    popupWindow = null;
};
function getPopupFeaturesString() {
    const popupWidth = 380;
    const popupHeight = 480;
    const left = window.screen.width / 2 - popupWidth / 2;
    const top = window.screen.height / 2 - popupHeight / 2;
    return `menubar=no,location=no,resizable=no,scrollbars=no,status=no,width=${popupWidth},height=${popupHeight},top=${top},left=${left}`;
}
function openAuthSessionAsync(args) {
    return __awaiter(this, void 0, void 0, function* () {
        const { url } = args;
        const state = yield getStateFromUrlOrGenerateAsync(url);
        window.sessionStorage.setItem(getHandle(), state);
        if (popupWindow == null || popupWindow.closed) {
            const popupWidth = 380;
            const popupHeight = 480;
            const left = window.screen.width / 2 - popupWidth / 2;
            const top = window.screen.height / 2 - popupHeight / 2;
            let popupFeaturesString = `menubar=no,location=no,resizable=no,scrollbars=no,status=no,width=${popupWidth},height=${popupHeight},top=${top},left=${left}`;
            popupWindow = window.open(url, "Core Login", popupFeaturesString);
            if (popupWindow) {
                try {
                    popupWindow.focus();
                }
                catch (e) {
                }
            }
            else {
                throw new Error("ERR_WEB_BROWSER_LOCKED: Popup window was blocked by the browser of failed to open.");
            }
            if (!popupWindow) {
                throw new Error("ERR_NO_POPUP: Failed to open login popup window");
            }
        }
        return new Promise((resolve) => __awaiter(this, void 0, void 0, function* () {
            const listener = (event) => {
                console.log("EVENT RECEIVED", event);
                if (!event.isTrusted) {
                    console.log("EVENT NOT TRUSTED");
                    return;
                }
                if (event.origin !== window.location.origin) {
                    console.log("ORIGIN MISMATCH", event.origin, window.location.origin);
                    return;
                }
                const { data } = event;
                const handle = window.sessionStorage.getItem(getHandle());
                console.log("COMPARING SENDER", data.sender, handle);
                if (data.sender === handle) {
                    dismissPopup();
                    resolve({ type: types_1.WebBrowserResultType.SUCCESS, url: data.url });
                }
            };
            window.addEventListener('message', listener, false);
            const interval = setInterval(() => {
                if (popupWindow === null || popupWindow === void 0 ? void 0 : popupWindow.closed) {
                    console.log("POPUP WINDOW CLOSED");
                    resolve({ type: types_1.WebBrowserResultType.DISMISS });
                    clearInterval(interval);
                    dismissPopup();
                }
            }, 1000);
            if (popupWindow) {
                listenerMap.set(popupWindow, { listener, interval });
            }
        }));
    });
}
exports.openAuthSessionAsync = openAuthSessionAsync;
function getStateFromUrlOrGenerateAsync(inputUrl) {
    return __awaiter(this, void 0, void 0, function* () {
        const url = new URL(inputUrl);
        if (url.searchParams.has('state') && typeof url.searchParams.get('state') === 'string') {
            return url.searchParams.get('state');
        }
        return yield generateStateAsync();
    });
}
function generateStateAsync() {
    return __awaiter(this, void 0, void 0, function* () {
        if (!isSubtleCryptoAvailable()) {
            throw new Error('ERR_WEB_BROWSER_CRYPTO: The current envorionment does not support Crypto');
        }
        const encoder = new TextEncoder();
        const data = generateRandom(10);
        const buffer = encoder.encode(data);
        const hashedData = yield crypto.subtle.digest('SHA-256', buffer);
        return btoa(String.fromCharCode(...new Uint8Array(hashedData)));
    });
}
function generateRandom(size) {
    let arr = new Uint8Array(size);
    if (arr.byteLength !== arr.length) {
        arr = new Uint8Array(arr.buffer);
    }
    const array = new Uint8Array(arr.length);
    if (isCryptoAvailable()) {
        window.crypto.getRandomValues(array);
    }
    else {
        for (let i = 0; i < size; i += 1) {
            array[i] = (Math.random() * CHARSET.length) | 0;
        }
    }
    return bufferToString(array);
}
exports.generateRandom = generateRandom;
function bufferToString(buffer) {
    const state = [];
    for (let i = 0; i < buffer.byteLength; i += 1) {
        const index = buffer[i] % CHARSET.length;
        state.push(CHARSET[index]);
    }
    return state.join('');
}
const CHARSET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
function isCryptoAvailable() {
    return !!(window === null || window === void 0 ? void 0 : window.crypto);
}
function isSubtleCryptoAvailable() {
    if (!isCryptoAvailable())
        return false;
    return !!window.crypto.subtle;
}
