"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CodedError = void 0;
class CodedError extends Error {
    constructor(code, message) {
        super(message);
        this.code = code;
    }
}
exports.CodedError = CodedError;
