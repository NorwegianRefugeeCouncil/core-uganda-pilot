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
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
exports.__esModule = true;
exports.generateRandom = exports.openAuthSessionAsync = exports.maybeCompleteAuthSession = void 0;
var types_1 = require("./types/types");
var popupWindow;
var listenerMap = new Map();
var getHandle = function () { return "NrcCoreBrowserRedirectHandle"; };
function maybeCompleteAuthSession() {
    var _a;
    var handle = window.sessionStorage.getItem(getHandle());
    if (!handle) {
        console.log("NO HANDLE");
        return { type: "failed", message: "No auth session is currently in progress" };
    }
    var url = window.location.href;
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
    var parent = (_a = window.opener) !== null && _a !== void 0 ? _a : window.parent;
    if (!parent) {
        throw new Error("ERR_WEB_BROWSER_REDIRECT: The window cannot complete the redirect request because the invoking window doesn't have a reference to its parent");
    }
    parent.postMessage({ url: url, sender: handle });
    return { type: "success", message: "attempting to complete auth" };
}
exports.maybeCompleteAuthSession = maybeCompleteAuthSession;
var dismissPopup = function () {
    console.log("DISMISSING POPUP");
    if (!popupWindow) {
        console.log("CANNOT DISMISS POPUP: no popup window");
        return;
    }
    popupWindow.close();
    if (listenerMap.has(popupWindow)) {
        var _a = listenerMap.get(popupWindow), listener = _a.listener, interval = _a.interval;
        clearInterval(interval);
        window.removeEventListener('message', listener);
        listenerMap["delete"](popupWindow);
        var handle = window.sessionStorage.getItem(getHandle());
        if (handle) {
            window.sessionStorage.removeItem(getHandle());
        }
    }
    popupWindow = null;
};
function getPopupFeaturesString() {
    var popupWidth = 380;
    var popupHeight = 480;
    var left = window.screen.width / 2 - popupWidth / 2;
    var top = window.screen.height / 2 - popupHeight / 2;
    return "menubar=no,location=no,resizable=no,scrollbars=no,status=no,width=" + popupWidth + ",height=" + popupHeight + ",top=" + top + ",left=" + left;
}
function openAuthSessionAsync(args) {
    return __awaiter(this, void 0, void 0, function () {
        var url, state, popupWidth, popupHeight, left, top_1, popupFeaturesString;
        var _this = this;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    url = args.url;
                    return [4 /*yield*/, getStateFromUrlOrGenerateAsync(url)];
                case 1:
                    state = _a.sent();
                    window.sessionStorage.setItem(getHandle(), state);
                    if (popupWindow == null || popupWindow.closed) {
                        popupWidth = 380;
                        popupHeight = 480;
                        left = window.screen.width / 2 - popupWidth / 2;
                        top_1 = window.screen.height / 2 - popupHeight / 2;
                        popupFeaturesString = "menubar=no,location=no,resizable=no,scrollbars=no,status=no,width=" + popupWidth + ",height=" + popupHeight + ",top=" + top_1 + ",left=" + left;
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
                    return [2 /*return*/, new Promise(function (resolve) { return __awaiter(_this, void 0, void 0, function () {
                            var listener, interval;
                            return __generator(this, function (_a) {
                                listener = function (event) {
                                    console.log("EVENT RECEIVED", event);
                                    if (!event.isTrusted) {
                                        console.log("EVENT NOT TRUSTED");
                                        return;
                                    }
                                    if (event.origin !== window.location.origin) {
                                        console.log("ORIGIN MISMATCH", event.origin, window.location.origin);
                                        return;
                                    }
                                    var data = event.data;
                                    var handle = window.sessionStorage.getItem(getHandle());
                                    console.log("COMPARING SENDER", data.sender, handle);
                                    if (data.sender === handle) {
                                        dismissPopup();
                                        resolve({ type: types_1.WebBrowserResultType.SUCCESS, url: data.url });
                                    }
                                };
                                window.addEventListener('message', listener, false);
                                interval = setInterval(function () {
                                    if (popupWindow === null || popupWindow === void 0 ? void 0 : popupWindow.closed) {
                                        console.log("POPUP WINDOW CLOSED");
                                        resolve({ type: types_1.WebBrowserResultType.DISMISS });
                                        clearInterval(interval);
                                        dismissPopup();
                                    }
                                }, 1000);
                                if (popupWindow) {
                                    listenerMap.set(popupWindow, { listener: listener, interval: interval });
                                }
                                return [2 /*return*/];
                            });
                        }); })];
            }
        });
    });
}
exports.openAuthSessionAsync = openAuthSessionAsync;
function getStateFromUrlOrGenerateAsync(inputUrl) {
    return __awaiter(this, void 0, void 0, function () {
        var url;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    url = new URL(inputUrl);
                    if (url.searchParams.has('state') && typeof url.searchParams.get('state') === 'string') {
                        return [2 /*return*/, url.searchParams.get('state')];
                    }
                    return [4 /*yield*/, generateStateAsync()];
                case 1: return [2 /*return*/, _a.sent()];
            }
        });
    });
}
function generateStateAsync() {
    return __awaiter(this, void 0, void 0, function () {
        var encoder, data, buffer, hashedData;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!isSubtleCryptoAvailable()) {
                        throw new Error('ERR_WEB_BROWSER_CRYPTO: The current envorionment does not support Crypto');
                    }
                    encoder = new TextEncoder();
                    data = generateRandom(10);
                    buffer = encoder.encode(data);
                    return [4 /*yield*/, crypto.subtle.digest('SHA-256', buffer)];
                case 1:
                    hashedData = _a.sent();
                    return [2 /*return*/, btoa(String.fromCharCode.apply(String, new Uint8Array(hashedData)))];
            }
        });
    });
}
function generateRandom(size) {
    var arr = new Uint8Array(size);
    if (arr.byteLength !== arr.length) {
        arr = new Uint8Array(arr.buffer);
    }
    var array = new Uint8Array(arr.length);
    if (isCryptoAvailable()) {
        window.crypto.getRandomValues(array);
    }
    else {
        for (var i = 0; i < size; i += 1) {
            array[i] = (Math.random() * CHARSET.length) | 0;
        }
    }
    return bufferToString(array);
}
exports.generateRandom = generateRandom;
function bufferToString(buffer) {
    var state = [];
    for (var i = 0; i < buffer.byteLength; i += 1) {
        var index = buffer[i] % CHARSET.length;
        state.push(CHARSET[index]);
    }
    return state.join('');
}
var CHARSET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
function isCryptoAvailable() {
    return !!(window === null || window === void 0 ? void 0 : window.crypto);
}
function isSubtleCryptoAvailable() {
    if (!isCryptoAvailable())
        return false;
    return !!window.crypto.subtle;
}
