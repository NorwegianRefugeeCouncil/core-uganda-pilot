import { OperatorFunction } from 'rxjs';
import { Case, CaseList, CaseType, CaseTypeList, Comment, CommentList, Country, CountryList, IdentificationDocument, IdentificationDocumentList, IdentificationDocumentType, IdentificationDocumentTypeList, Individual, IndividualList, Membership, MembershipList, Nationality, NationalityList, Party, PartyAttributeDefinition, PartyAttributeDefinitionList, PartyType, PartyTypeList, Relationship, RelationshipList, RelationshipType, RelationshipTypeList, Team, TeamList } from './types/models';
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
declare class CaseClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Case>;
    Create(): OperatorFunction<Case, Case>;
    Update(): OperatorFunction<Case, Case>;
    List(): OperatorFunction<any, CaseList>;
}
declare class CaseTypeClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, CaseType>;
    Create(): OperatorFunction<CaseType, CaseType>;
    Update(): OperatorFunction<CaseType, CaseType>;
    List(): OperatorFunction<any, CaseTypeList>;
}
declare class CommentClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Comment>;
    Create(): OperatorFunction<Comment, Comment>;
    Update(): OperatorFunction<Comment, Comment>;
    List(): OperatorFunction<any, CommentList>;
    Delete(): OperatorFunction<string, any>;
}
export declare class CMSClient {
    private readonly _scheme;
    private readonly _host;
    private readonly _headers;
    private readonly _client;
    constructor(scheme: string, host: string, headers: Headers);
    Cases(): CaseClient;
    CaseTypes(): CaseTypeClient;
    Comments(): CommentClient;
}
declare class CountryClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Country>;
    Create(): OperatorFunction<Country, Country>;
    Update(): OperatorFunction<Country, Country>;
    List(): OperatorFunction<any, CountryList>;
}
declare class IdentificationDocumentClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, IdentificationDocument>;
    Create(): OperatorFunction<IdentificationDocument, IdentificationDocument>;
    Update(): OperatorFunction<IdentificationDocument, IdentificationDocument>;
    List(): OperatorFunction<any, IdentificationDocumentList>;
    Delete(): OperatorFunction<string, any>;
}
declare class IdentificationDocumentTypeClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, IdentificationDocumentType>;
    Create(): OperatorFunction<IdentificationDocumentType, IdentificationDocumentType>;
    Update(): OperatorFunction<IdentificationDocumentType, IdentificationDocumentType>;
    List(): OperatorFunction<any, IdentificationDocumentTypeList>;
}
declare class IndividualClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Individual>;
    Create(): OperatorFunction<Individual, Individual>;
    Update(): OperatorFunction<Individual, Individual>;
    List(): OperatorFunction<any, IndividualList>;
}
declare class MembershipClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Membership>;
    Create(): OperatorFunction<Membership, Membership>;
    Update(): OperatorFunction<Membership, Membership>;
    List(): OperatorFunction<any, MembershipList>;
}
declare class NationalityClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Nationality>;
    Create(): OperatorFunction<Nationality, Nationality>;
    Update(): OperatorFunction<Nationality, Nationality>;
    List(): OperatorFunction<any, NationalityList>;
}
declare class PartyAttributeDefinitionClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, PartyAttributeDefinition>;
    Create(): OperatorFunction<PartyAttributeDefinition, PartyAttributeDefinition>;
    Update(): OperatorFunction<PartyAttributeDefinition, PartyAttributeDefinition>;
    List(): OperatorFunction<any, PartyAttributeDefinitionList>;
}
declare class PartyClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Party>;
    Update(): OperatorFunction<Party, Party>;
}
declare class PartyTypeClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, PartyType>;
    Create(): OperatorFunction<PartyType, PartyType>;
    Update(): OperatorFunction<PartyType, PartyType>;
    List(): OperatorFunction<any, PartyTypeList>;
}
declare class RelationshipClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Relationship>;
    Create(): OperatorFunction<Relationship, Relationship>;
    Update(): OperatorFunction<Relationship, Relationship>;
    List(): OperatorFunction<any, RelationshipList>;
    Delete(): OperatorFunction<string, any>;
}
declare class RelationshipTypeClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, RelationshipType>;
    Create(): OperatorFunction<RelationshipType, RelationshipType>;
    Update(): OperatorFunction<RelationshipType, RelationshipType>;
    List(): OperatorFunction<any, RelationshipTypeList>;
}
declare class TeamClient {
    private client;
    execute: () => OperatorFunction<Request, Response>;
    constructor(client: Client);
    Get(): OperatorFunction<string, Team>;
    Create(): OperatorFunction<Team, Team>;
    Update(): OperatorFunction<Team, Team>;
    List(): OperatorFunction<any, TeamList>;
}
export declare class IAMClient {
    private readonly _scheme;
    private readonly _host;
    private readonly _headers;
    private readonly _client;
    constructor(scheme: string, host: string, headers: Headers);
    Countries(): CountryClient;
    IdentificationDocuments(): IdentificationDocumentClient;
    IdentificationDocumentTypes(): IdentificationDocumentTypeClient;
    Individuals(): IndividualClient;
    Memberships(): MembershipClient;
    Nationalities(): NationalityClient;
    PartyAttributeDefinitions(): PartyAttributeDefinitionClient;
    Parties(): PartyClient;
    PartyTypes(): PartyTypeClient;
    Relationships(): RelationshipClient;
    RelationshipTypes(): RelationshipTypeClient;
    Teams(): TeamClient;
}
export {};
