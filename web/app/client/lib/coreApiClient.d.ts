import { OperatorFunction } from 'rxjs';
import { Party } from './types/models';
export declare type AjaxRequestOptions = {
    url: string;
    headers: Headers;
    method: string;
    async: boolean;
    timeout: number;
    crossDomain: boolean;
    withCredentials: boolean;
    body: any;
};
export declare function prepareRequestOptions(req: Request): AjaxRequestOptions;
declare class Client {
    readonly _scheme: string;
    readonly _host: string;
    _headers: Headers;
    constructor(host: string, scheme: string, headers: Headers);
    get(): Request;
    post(): Request;
    put(): Request;
    delete(): Request;
    do(): OperatorFunction<Request, Response>;
}
interface URLValues {
    [key: string]: string;
}
export interface Headers {
    [key: string]: string[];
}
export declare class Response {
    readonly body: any;
    private readonly _body;
    constructor(body: any);
    as<T>(): T;
}
export declare class Request {
    _client: Client;
    _error: Error;
    _path: string;
    _verb: string;
    _body: any;
    _params: URLValues;
    _pathParams: URLValues;
    _headers: Headers;
    constructor(client: Client);
    verb(verb: string): Request;
    get(): Request;
    put(): Request;
    post(): Request;
    delete(): Request;
    path(path: string): Request;
    body(body: any): Request;
    params(params: URLValues): Request;
    pathParam(key: string, value: string): Request;
    headers(headers: Headers): Request;
}
declare class PartyClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Party>;
    Update(): OperatorFunction<Party, Party>;
}
export declare class IAMClient {
    private readonly _scheme;
    private readonly _host;
    private readonly _headers;
    private readonly _client;
    constructor(scheme: string, host: string, headers: Headers);
    Parties(): PartyClient;
}
export {};
