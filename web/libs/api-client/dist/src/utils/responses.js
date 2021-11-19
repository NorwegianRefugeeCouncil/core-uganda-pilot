"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.clientResponse = void 0;
function errorResponse(request, r) {
    return {
        request: request,
        response: undefined,
        status: r.request,
        statusCode: r.status,
        error: r.data,
        success: false,
    };
}
function successResponse(request, r) {
    return {
        request: request,
        response: r.data,
        status: r.statusText,
        statusCode: r.status,
        error: undefined,
        success: true,
    };
}
function clientResponse(r, request, expectedStatusCode) {
    return r.status !== expectedStatusCode
        ? errorResponse(request, r)
        : successResponse(request, r);
}
exports.clientResponse = clientResponse;
