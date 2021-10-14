import { Observable, OperatorFunction } from 'rxjs';
import { Case, CaseList, CaseListOptions, Comment, CaseType, CaseTypeListOptions, CommentListOptions, Party, PartySearchOptions, PartyListOptions, PartyList, Country, CountryListOptions, CountryList, IdentificationDocument, IdentificationDocumentList, IdentificationDocumentListOptions, IdentificationDocumentType, IdentificationDocumentTypeList, IdentificationDocumentTypeListOptions, Individual, IndividualList, IndividualListOptions, Membership, MembershipList, MembershipListOptions, Nationality, NationalityList, NationalityListOptions, PartyAttributeDefinition, PartyAttributeDefinitionList, PartyAttributeDefinitionListOptions, PartyType, PartyTypeList, PartyTypeListOptions, Relationship, RelationshipList, RelationshipListOptions, RelationshipType, RelationshipTypeList, RelationshipTypeListOptions, Team, TeamList, TeamListOptions, CommentList, CaseTypeList } from './types/models';
import { AjaxResponse } from 'rxjs/ajax';

export declare interface Headers {
    [key: string]: string[]
}

export declare interface URLValues {
    [key: string]: string
}

declare class Request {
    constructor(client: Client)
    verb(verb: string): Request
    get(): Request
    put(): Request
    post(): Request
    delete(): Request
    path(path: string): Request
    body(body: any): Request
    params(params: URLValues): Request
    pathParam(key: string, value: string): Request
    headers(headers: Headers): Request
}

declare class Response {
    constructor(body: any)
    as<T>(): T
}

declare class Client {
    constructor(host: string, scheme: string, headers: Headers)
    verb(verb: string): Request
    get(): Request
    post(): Request
    put(): Request
    delete(): Request
    do(): OperatorFunction<Request, Response>
}

declare class PartyClient {
    constructor(client: Client)
    Get(): OperatorFunction<string, Party>
    Update(): OperatorFunction<Party, Party>
}

export declare class IAMClient {
    constructor(scheme: string, host: string, headers: Headers)
    Parties(): PartyClient
}

export {};
