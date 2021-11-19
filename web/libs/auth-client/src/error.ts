export class CodedError extends Error {
    code: string;
    info?: any;
    constructor(code: string, message: string) {
        super(message);
        this.code = code;
    }
}
