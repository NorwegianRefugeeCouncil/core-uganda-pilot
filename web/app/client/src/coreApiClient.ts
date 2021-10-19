import {map, OperatorFunction} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {
    Case,
    CaseList,
    CaseListOptions,
    CaseType,
    CaseTypeList,
    CaseTypeListOptions,
    Comment,
    CommentList,
    CommentListOptions,
    Country,
    CountryList,
    CountryListOptions,
    IdentificationDocument,
    IdentificationDocumentList,
    IdentificationDocumentListOptions,
    IdentificationDocumentType,
    IdentificationDocumentTypeList,
    IdentificationDocumentTypeListOptions,
    Individual,
    IndividualList,
    IndividualListOptions,
    Membership,
    MembershipList,
    MembershipListOptions,
    Nationality,
    NationalityList,
    NationalityListOptions,
    Party,
    PartyAttributeDefinition,
    PartyAttributeDefinitionList,
    PartyAttributeDefinitionListOptions,
    PartyList,
    PartyListOptions,
    PartyType,
    PartyTypeList,
    PartyTypeListOptions,
    Relationship,
    RelationshipList,
    RelationshipListOptions,
    RelationshipType,
    RelationshipTypeList,
    RelationshipTypeListOptions,
    Team,
    TeamList,
    TeamListOptions
} from './types/models';
import {ajax} from 'rxjs/ajax';
import {XMLHttpRequest} from 'xhr2';

// needed for rxjs/ajax compatibility outside the browser
global.XMLHttpRequest = global.XMLHttpRequest ? global.XMLHttpRequest : XMLHttpRequest;

export type AjaxRequestOptions = {
    url: string,
    headers: Headers,
    method: string,
    async: boolean,
    timeout: number,
    crossDomain: boolean,
    withCredentials: boolean,
    body: any
}

export function prepareRequestOptions(req: Request): AjaxRequestOptions {
    let url = `${req._client._scheme}://${req._client._host}`
    if (req._path) {
        url = url + req._path;
    }

    let headers: Headers = {};

    copyHeaders(req._client._headers, headers);
    copyHeaders(req._headers, headers);

    if (!headers['Content-Type']) {
        headers['Content-Type'] = ['application/json'];
    }

    if (!headers['Accept']) {
        headers['Accept'] = ['application/json'];
    }

    const queryParams = req._params
    if (queryParams != null && Object.keys(queryParams as Object).length > 0) {
        url = appendQueryParams(url, queryParams)
    }
    const pathParams = req._pathParams
    if (pathParams != null && Object.keys(pathParams).length > 0) {
        url = replacePathParams(url, pathParams)
    }

    return {
        url: url,
        headers: headers,
        method: req._verb,
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        body: req._body
    } as AjaxRequestOptions
}

function copyHeaders(from: Headers, to: Headers) {
    for (let headerKey in from) {
        if (from.hasOwnProperty(headerKey)) {
            for (let headerValue of from[headerKey]) {
                if (!to[headerKey]) {
                    to[headerKey] = [];
                }
                to[headerKey].push(headerValue);
            }
        }
    }
}

function replacePathParams(s: string, params: URLValues): string {
    let tmp = s
    for (const [key, value] of Object.entries(params)) {
        tmp = tmp.replace(`:${key}`, value)
    }
    return tmp
}

function appendQueryParams(s: string, params: URLValues): string {
    let paramStrings = []
    for (const [key, value] of Object.entries(params)) {
        paramStrings.push(`${key}=${value}`)
    }
    return paramStrings.length > 0 ? `${s}?${paramStrings.join("&")}` : s
}

function urlValuesNotEmpty(u: URLValues): boolean {
    if (u == null) return false
    return Object.keys(u).length > 0
}

// TODO
function isErrorResponse(data: any): boolean {
    return false;
}

class Client {
    readonly _scheme: string
    readonly _host: string
    _headers: Headers

    constructor(host: string, scheme: string, headers: Headers) {
        this._host = host
        this._scheme = scheme
        this._headers = headers
    }

    get(): Request {
        const r = new Request(this).get()
        return r
    }

    post(): Request {
        return new Request(this).post();
    }

    put(): Request {
        return new Request(this).put();
    }

    delete(): Request {
        return new Request(this).delete();
    }

    do(): OperatorFunction<Request, Response> {
        return source => {
            return source.pipe(
                switchMap(req => {
                    const tmpReq = req
                    const ro = prepareRequestOptions(tmpReq as Request)
                    return ajax(
                        {
                            url: ro.url,
                            headers: ro.headers,
                            method: ro.method,
                            async: ro.async,
                            timeout: ro.timeout,
                            crossDomain: ro.crossDomain,
                            withCredentials: ro.withCredentials,
                            body: ro.body
                        }
                    );
                }),
                map(ajaxResponse => {
                    if (ajaxResponse.status > 399) {
                        if (isErrorResponse(ajaxResponse.response)) {
                            return new Response(ajaxResponse.response);
                        } else {
                            return new Response({error: 'server error', status: 500});
                        }
                    }
                    return new Response(ajaxResponse.response);
                })
            );
        };
    }
}

interface URLValues {
    [key: string]: string;
}

export interface Headers {
    [key: string]: string[];
}

export class Response {
    private readonly _body: any

    public constructor(readonly body: any) {
        this._body = body
    }

    as<T>(): T {
        return this._body as T
    }
}

export class Request {
    public _client: Client
    public _error: Error
    public _path: string
    public _verb: string
    public _body: any
    public _params: URLValues
    public _pathParams: URLValues
    public _headers: Headers

    public constructor(client: Client) {
        this._client = client
    }

    public verb(verb: string): Request {
        this._verb = verb
        return this
    }

    public get(): Request {
        return this.verb('GET')
    }

    public put(): Request {
        return this.verb('PUT')
    }

    public post(): Request {
        return this.verb('POST')
    }

    public delete(): Request {
        return this.verb('DELETE')
    }

    public path(path: string): Request {
        this._path = path
        return this
    }

    public body(body: any): Request {
        this._body = body
        return this
    }

    public params(params: URLValues): Request {
        this._params = params
        return this
    }

    public pathParam(key: string, value: string): Request {
        if (!this._pathParams) {
            this._pathParams = {}
        }
        this._pathParams[key] = value
        return this
    }

    public headers(headers: Headers): Request {
        this._headers = headers
        return this
    }

}

function responseAs<T>(): OperatorFunction<Response, T> {
    return source => {
        return source.pipe(map(resp => {
            if (isErrorResponse(resp.body)) {
                throw new Error(resp.body)
            } else {
                return resp.body
            }
        }))
    }
}

// --- CMS ---

function buildCMSPath(module: string, endpoint: string = ""): string {
    return `/apis/cms/v1/${module}${endpoint}`
}

class CaseClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Case> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildCMSPath('cases', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Case>()
        )
    }

    Create(): OperatorFunction<Case, Case> {
        return c$ => c$.pipe(
            map(c => this.client.post().body(c).path(buildCMSPath(('cases')))),
            this.execute(),
            responseAs<Case>()
        )
    }

    Update(): OperatorFunction<Case, Case> {
        return c$ => c$.pipe(
            map(c => this.client.put().body(c).path(buildCMSPath('cases', '/:id')).pathParam('id', c.id)),
            this.execute(),
            responseAs<Case>()
        )
    }

    List(): OperatorFunction<CaseListOptions, CaseList> {
        return clo$ => clo$.pipe(
            map(clo => clo ?
                this.client.get().path(buildCMSPath('cases')).params(clo as unknown as URLValues)
                : this.client.get().path(buildCMSPath('cases'))
            ),
            this.execute(),
            responseAs<CaseList>()
        )
    }
}

class CaseTypeClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, CaseType> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildCMSPath('casetypes', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<CaseType>()
        )
    }

    Create(): OperatorFunction<CaseType, CaseType> {
        return c$ => c$.pipe(
            map(c => this.client.post().body(c).path(buildCMSPath(('casetypes')))),
            this.execute(),
            responseAs<CaseType>()
        )
    }

    Update(): OperatorFunction<CaseType, CaseType> {
        return c$ => c$.pipe(
            map(c => this.client.put().body(c).path(buildCMSPath('casetypes', '/:id')).pathParam('id', c.id)),
            this.execute(),
            responseAs<CaseType>()
        )
    }

    List(): OperatorFunction<CaseTypeListOptions, CaseTypeList> {
        return ctlo$ => ctlo$.pipe(
            map(ctlo => ctlo ?
                this.client.get().path(buildCMSPath('casetypes')).params(ctlo as unknown as URLValues)
                : this.client.get().path(buildCMSPath('casetypes'))
            ),
            this.execute(),
            responseAs<CaseTypeList>()
        )
    }
}

class CommentClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Comment> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildCMSPath('comments', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Comment>()
        )
    }

    Create(): OperatorFunction<Comment, Comment> {
        return c$ => c$.pipe(
            map(c => this.client.post().body(c).path(buildCMSPath(('comments')))),
            this.execute(),
            responseAs<Comment>()
        )
    }

    Update(): OperatorFunction<Comment, Comment> {
        return c$ => c$.pipe(
            map(c => this.client.put().body(c).path(buildCMSPath('comments', '/:id')).pathParam('id', c.id)),
            this.execute(),
            responseAs<Comment>()
        )
    }

    List(): OperatorFunction<CommentListOptions, CommentList> {
        return clo$ => clo$.pipe(
            map(clo => clo ?
                this.client.get().path(buildCMSPath('comments')).params(clo as unknown as URLValues)
                : this.client.get().path(buildCMSPath('comments'))
            ),
            this.execute(),
            responseAs<CommentList>()
        )
    }

    Delete(): OperatorFunction<string, any> {
        return id$ => id$.pipe(
            map(id => this.client.delete().path(buildCMSPath('comments', '/:id')).pathParam('id', id)),
            this.execute()
        )
    }
}

export class CMSClient {
    private readonly _scheme: string
    private readonly _host: string
    private readonly _headers: Headers
    private readonly _client: Client

    public constructor(scheme: string, host: string, headers: Headers) {
        this._host = host
        this._scheme = scheme
        this._headers = headers
        this._client = new Client(this._host, this._scheme, this._headers)
    }

    public Cases() {
        return new CaseClient(this._client)
    }

    public CaseTypes() {
        return new CaseTypeClient(this._client)
    }

    public Comments() {
        return new CommentClient(this._client)
    }
}

// --- IAM ---

function buildIAMPath(module: string, endpoint: string = ""): string {
    return `/apis/iam/v1/${module}${endpoint}`
}

class CountryClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Country> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('countries', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Country>()
        )
    }

    Create(): OperatorFunction<Country, Country> {
        return country$ => country$.pipe(
            map(country => this.client.post().body(country).path(buildIAMPath(('countries')))),
            this.execute(),
            responseAs<Country>()
        )
    }

    Update(): OperatorFunction<Country, Country> {
        return country$ => country$.pipe(
            map(country => this.client.put().body(country).path(buildIAMPath('countries', '/:id')).pathParam('id', country.id)),
            this.execute(),
            responseAs<Country>()
        )
    }

    List(): OperatorFunction<CountryListOptions, CountryList> {
        return clo$ => clo$.pipe(
            map(clo => clo ?
                this.client.get().path(buildIAMPath('countries')).params(clo as URLValues)
                : this.client.get().path(buildIAMPath('countries'))
            ),
            this.execute(),
            responseAs<CountryList>()
        )
    }
}

class IdentificationDocumentClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, IdentificationDocument> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('identificationdocuments', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<IdentificationDocument>()
        )
    }

    Create(): OperatorFunction<IdentificationDocument, IdentificationDocument> {
        return id$ => id$.pipe(
            map(id => this.client.post().body(id).path(buildIAMPath(('identificationdocuments')))),
            this.execute(),
            responseAs<IdentificationDocument>()
        )
    }

    Update(): OperatorFunction<IdentificationDocument, IdentificationDocument> {
        return id$ => id$.pipe(
            map(id => this.client.put().body(id).path(buildIAMPath('identificationdocuments', '/:id')).pathParam('id', id.id)),
            this.execute(),
            responseAs<IdentificationDocument>()
        )
    }

    List(): OperatorFunction<IdentificationDocumentListOptions, IdentificationDocumentList> {
        return idlo$ => idlo$.pipe(
            map(idlo => idlo ?
                this.client.get().path(buildIAMPath('identificationdocuments')).params(idlo as URLValues)
                : this.client.get().path(buildIAMPath('identificationdocuments'))
            ),
            this.execute(),
            responseAs<IdentificationDocumentList>()
        )
    }

    Delete(): OperatorFunction<string, any> {
        return id$ => id$.pipe(
            map(id => this.client.delete().path(buildIAMPath('identificationdocuments', '/:id')).pathParam('id', id)),
            this.execute()
        )
    }
}

class IdentificationDocumentTypeClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, IdentificationDocumentType> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('identificationdocumenttypes', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<IdentificationDocumentType>()
        )
    }

    Create(): OperatorFunction<IdentificationDocumentType, IdentificationDocumentType> {
        return idt$ => idt$.pipe(
            map(idt => this.client.post().body(idt).path(buildIAMPath(('identificationdocumenttypes')))),
            this.execute(),
            responseAs<IdentificationDocumentType>()
        )
    }

    Update(): OperatorFunction<IdentificationDocumentType, IdentificationDocumentType> {
        return idt$ => idt$.pipe(
            map(idt => this.client.put().body(idt).path(buildIAMPath('identificationdocumenttypes', '/:id')).pathParam('id', idt.id)),
            this.execute(),
            responseAs<IdentificationDocumentType>()
        )
    }

    List(): OperatorFunction<IdentificationDocumentTypeListOptions, IdentificationDocumentTypeList> {
        return idtlo$ => idtlo$.pipe(
            map(idtlo => idtlo ?
                this.client.get().path(buildIAMPath('identificationdocumenttypes')).params(idtlo as URLValues)
                : this.client.get().path(buildIAMPath('identificationdocumenttypes'))
            ),
            this.execute(),
            responseAs<IdentificationDocumentTypeList>()
        )
    }
}

class IndividualClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Individual> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('individuals', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Individual>()
        )
    }

    Create(): OperatorFunction<Individual, Individual> {
        return ind$ => ind$.pipe(
            map(ind => this.client.post().body(ind).path(buildIAMPath(('individuals')))),
            this.execute(),
            responseAs<Individual>()
        )
    }

    Update(): OperatorFunction<Individual, Individual> {
        return ind$ => ind$.pipe(
            map(ind => this.client.put().body(ind).path(buildIAMPath('individuals', '/:id')).pathParam('id', ind.id)),
            this.execute(),
            responseAs<Individual>()
        )
    }

    List(): OperatorFunction<IndividualListOptions, IndividualList> {
        return ilo$ => ilo$.pipe(
            map(ilo => ilo ?
                this.client.get().path(buildIAMPath('individuals')).params(ilo as URLValues)
                : this.client.get().path(buildIAMPath('individuals'))
            ),
            this.execute(),
            responseAs<IndividualList>()
        )
    }
}

class MembershipClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Membership> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('memberships', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Membership>()
        )
    }

    Create(): OperatorFunction<Membership, Membership> {
        return mem$ => mem$.pipe(
            map(mem => this.client.post().body(mem).path(buildIAMPath(('memberships')))),
            this.execute(),
            responseAs<Membership>()
        )
    }

    Update(): OperatorFunction<Membership, Membership> {
        return mem$ => mem$.pipe(
            map(mem => this.client.put().body(mem).path(buildIAMPath('memberships', '/:id')).pathParam('id', mem.id)),
            this.execute(),
            responseAs<Membership>()
        )
    }

    List(): OperatorFunction<MembershipListOptions, MembershipList> {
        return mlo$ => mlo$.pipe(
            map(mlo => mlo ?
                this.client.get().path(buildIAMPath('memberships')).params(mlo as URLValues)
                : this.client.get().path(buildIAMPath('memberships'))
            ),
            this.execute(),
            responseAs<MembershipList>()
        )
    }
}

class NationalityClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Nationality> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('nationalities', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Nationality>()
        )
    }

    Create(): OperatorFunction<Nationality, Nationality> {
        return nat$ => nat$.pipe(
            map(nat => this.client.post().body(nat).path(buildIAMPath(('nationalities')))),
            this.execute(),
            responseAs<Nationality>()
        )
    }

    Update(): OperatorFunction<Nationality, Nationality> {
        return nat$ => nat$.pipe(
            map(nat => this.client.put().body(nat).path(buildIAMPath('nationalities', '/:id')).pathParam('id', nat.id)),
            this.execute(),
            responseAs<Nationality>()
        )
    }

    List(): OperatorFunction<NationalityListOptions, NationalityList> {
        return nlo$ => nlo$.pipe(
            map(nlo => nlo ?
                this.client.get().path(buildIAMPath('nationalities')).params(nlo as URLValues)
                : this.client.get().path(buildIAMPath('nationalities'))
            ),
            this.execute(),
            responseAs<NationalityList>()
        )
    }
}

class PartyAttributeDefinitionClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, PartyAttributeDefinition> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('attributes', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<PartyAttributeDefinition>()
        )
    }

    Create(): OperatorFunction<PartyAttributeDefinition, PartyAttributeDefinition> {
        return pad$ => pad$.pipe(
            map(pad => this.client.post().body(pad).path(buildIAMPath(('attributes')))),
            this.execute(),
            responseAs<PartyAttributeDefinition>()
        )
    }

    Update(): OperatorFunction<PartyAttributeDefinition, PartyAttributeDefinition> {
        return pad$ => pad$.pipe(
            map(pad => this.client.put().body(pad).path(buildIAMPath('attributes', '/:id')).pathParam('id', pad.id)),
            this.execute(),
            responseAs<PartyAttributeDefinition>()
        )
    }

    List(): OperatorFunction<PartyAttributeDefinitionListOptions, PartyAttributeDefinitionList> {
        return padlo$ => padlo$.pipe(
            map(padlo => padlo ?
                this.client.get().path(buildIAMPath('attributes')).params(padlo as URLValues)
                : this.client.get().path(buildIAMPath('attributes'))
            ),
            this.execute(),
            responseAs<PartyAttributeDefinitionList>()
        )
    }
}

class PartyClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Party> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('parties', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Party>()
        )
    }

    Create(): OperatorFunction<Party, Party> {
        return party$ => party$.pipe(
            map(party => this.client.post().body(party).path(buildIAMPath('parties'))),
            this.execute(),
            responseAs<Party>()
        )
    }

    Update(): OperatorFunction<Party, Party> {
        return party$ => party$.pipe(
            map(party => this.client.put().body(party).path(buildIAMPath('parties', '/:id')).pathParam('id', party.id)),
            this.execute(),
            responseAs<Party>()
        )
    }

    List(): OperatorFunction<PartyListOptions, PartyList> {
        return plo$ => plo$.pipe(
            map(plo => plo ?
                this.client.get().path(buildIAMPath('parties')).params(plo as unknown as URLValues)
                : this.client.get().path(buildIAMPath('parties'))
            ),
            this.execute(),
            responseAs<PartyList>()
        )
    }
}

class PartyTypeClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, PartyType> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('partytypes', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<PartyType>()
        )
    }

    Create(): OperatorFunction<PartyType, PartyType> {
        return pt$ => pt$.pipe(
            map(pt => this.client.post().body(pt).path(buildIAMPath(('partytypes')))),
            this.execute(),
            responseAs<PartyType>()
        )
    }

    Update(): OperatorFunction<PartyType, PartyType> {
        return pt$ => pt$.pipe(
            map(pt => this.client.put().body(pt).path(buildIAMPath('partytypes', '/:id')).pathParam('id', pt.id)),
            this.execute(),
            responseAs<PartyType>()
        )
    }

    List(): OperatorFunction<PartyTypeListOptions, PartyTypeList> {
        return ptlo$ => ptlo$.pipe(
            map(ptlo => ptlo ?
                this.client.get().path(buildIAMPath('partytypes')).params(ptlo as URLValues)
                : this.client.get().path(buildIAMPath('partytypes'))
            ),
            this.execute(),
            responseAs<PartyTypeList>()
        )
    }
}

class RelationshipClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Relationship> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('relationships', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Relationship>()
        )
    }

    Create(): OperatorFunction<Relationship, Relationship> {
        return r$ => r$.pipe(
            map(r => this.client.post().body(r).path(buildIAMPath(('relationships')))),
            this.execute(),
            responseAs<Relationship>()
        )
    }

    Update(): OperatorFunction<Relationship, Relationship> {
        return r$ => r$.pipe(
            map(r => this.client.put().body(r).path(buildIAMPath('relationships', '/:id')).pathParam('id', r.id)),
            this.execute(),
            responseAs<Relationship>()
        )
    }

    List(): OperatorFunction<RelationshipListOptions, RelationshipList> {
        return rlo$ => rlo$.pipe(
            map(rlo => rlo ?
                this.client.get().path(buildIAMPath('relationships')).params(rlo as URLValues)
                : this.client.get().path(buildIAMPath('relationships'))
            ),
            this.execute(),
            responseAs<RelationshipList>()
        )
    }

    Delete(): OperatorFunction<string, any> {
        return id$ => id$.pipe(
            map(id => this.client.delete().path(buildIAMPath('relationships', '/:id')).pathParam('id', id)),
            this.execute()
        )
    }
}

class RelationshipTypeClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, RelationshipType> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('relationshiptypes', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<RelationshipType>()
        )
    }

    Create(): OperatorFunction<RelationshipType, RelationshipType> {
        return rt$ => rt$.pipe(
            map(rt => this.client.post().body(rt).path(buildIAMPath(('relationshiptypes')))),
            this.execute(),
            responseAs<RelationshipType>()
        )
    }

    Update(): OperatorFunction<RelationshipType, RelationshipType> {
        return rt$ => rt$.pipe(
            map(rt => this.client.put().body(rt).path(buildIAMPath('relationshiptypes', '/:id')).pathParam('id', rt.id)),
            this.execute(),
            responseAs<RelationshipType>()
        )
    }

    List(): OperatorFunction<RelationshipTypeListOptions, RelationshipTypeList> {
        return rtlo$ => rtlo$.pipe(
            map(rtlo => rtlo ?
                this.client.get().path(buildIAMPath('relationshiptypes')).params(rtlo as URLValues)
                : this.client.get().path(buildIAMPath('relationshiptypes'))
            ),
            this.execute(),
            responseAs<RelationshipTypeList>()
        )
    }
}

class TeamClient {

    private client: Client;

    execute: () => OperatorFunction<Request, Response>;

    public constructor(client: Client) {
        this.client = client
        this.execute = client.do
    }

    Get(): OperatorFunction<string, Team> {
        return id$ => id$.pipe(
            map(id => this.client.get().path(buildIAMPath('teams', '/:id')).pathParam('id', id)),
            this.execute(),
            responseAs<Team>()
        )
    }

    Create(): OperatorFunction<Team, Team> {
        return t$ => t$.pipe(
            map(t => this.client.post().body(t).path(buildIAMPath(('teams')))),
            this.execute(),
            responseAs<Team>()
        )
    }

    Update(): OperatorFunction<Team, Team> {
        return t$ => t$.pipe(
            map(t => this.client.put().body(t).path(buildIAMPath('teams', '/:id')).pathParam('id', t.id)),
            this.execute(),
            responseAs<Team>()
        )
    }

    List(): OperatorFunction<TeamListOptions, TeamList> {
        return tlo$ => tlo$.pipe(
            map(tlo => tlo ?
                this.client.get().path(buildIAMPath('teams')).params(tlo as URLValues)
                : this.client.get().path(buildIAMPath('teams'))
            ),
            this.execute(),
            responseAs<TeamList>()
        )
    }
}

export class IAMClient {
    private readonly _scheme: string
    private readonly _host: string
    private readonly _headers: Headers
    private readonly _client: Client

    public constructor(scheme: string, host: string, headers: Headers) {
        this._host = host
        this._scheme = scheme
        this._headers = headers
        this._client = new Client(this._host, this._scheme, this._headers)
    }

    public Countries() {
        return new CountryClient(this._client)
    }

    public IdentificationDocuments() {
        return new IdentificationDocumentClient(this._client)
    }

    public IdentificationDocumentTypes() {
        return new IdentificationDocumentTypeClient(this._client)
    }

    public Individuals() {
        return new IndividualClient(this._client)
    }

    public Memberships() {
        return new MembershipClient(this._client)
    }

    public Nationalities() {
        return new NationalityClient(this._client)
    }

    public PartyAttributeDefinitions() {
        return new PartyAttributeDefinitionClient(this._client)
    }

    public Parties() {
        return new PartyClient(this._client)
    }

    public PartyTypes() {
        return new PartyTypeClient(this._client)
    }

    public Relationships() {
        return new RelationshipClient(this._client)
    }

    public RelationshipTypes() {
        return new RelationshipTypeClient(this._client)
    }

    public Teams() {
        return new TeamClient(this._client)
    }
}
