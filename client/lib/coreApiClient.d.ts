import { Observable } from 'rxjs';
import { Case, CaseList, CaseListOptions, Comment, CaseType, CaseTypeListOptions, CommentListOptions, Party, PartySearchOptions, PartyListOptions, PartyList, Country, CountryListOptions, CountryList, IdentificationDocument, IdentificationDocumentList, IdentificationDocumentListOptions, IdentificationDocumentType, IdentificationDocumentTypeList, IdentificationDocumentTypeListOptions, Individual, IndividualList, IndividualListOptions, Membership, MembershipList, MembershipListOptions, Nationality, NationalityList, NationalityListOptions, PartyAttributeDefinition, PartyAttributeDefinitionList, PartyAttributeDefinitionListOptions, PartyType, PartyTypeList, PartyTypeListOptions, Relationship, RelationshipList, RelationshipListOptions, RelationshipType, RelationshipTypeList, RelationshipTypeListOptions, Team, TeamList, TeamListOptions, CommentList, CaseTypeList } from './types/models';
import { AjaxResponse } from 'rxjs/ajax';
declare class HttpClient<T> {
    headers: {
        "X-Authenticated-User-Subject"?: string;
    };
    constructor(shouldAddAuthPassthroughHeader: boolean);
    get(url: string): Observable<AjaxResponse<T>>;
    getCustom<R>(url: string): Observable<AjaxResponse<R>>;
    put(url: string, body: T): Observable<AjaxResponse<T>>;
    post(url: string, body: T): Observable<AjaxResponse<T>>;
    postCustom<B, R>(url: string, body: B): Observable<AjaxResponse<R>>;
    delete(url: string): Observable<AjaxResponse<T>>;
}
declare class TeamClient {
    readonly httpClient: HttpClient<Team>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Team>>;
    Create(t: Team): Observable<AjaxResponse<Team>>;
    Update(t: Team): Observable<AjaxResponse<Team>>;
    List(lo: TeamListOptions): Observable<AjaxResponse<TeamList>>;
}
declare class RelationshipTypeClient {
    readonly httpClient: HttpClient<RelationshipType>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<RelationshipType>>;
    Create(r: RelationshipType): Observable<AjaxResponse<RelationshipType>>;
    Update(r: RelationshipType): Observable<AjaxResponse<RelationshipType>>;
    List(lo: RelationshipTypeListOptions): Observable<AjaxResponse<RelationshipTypeList>>;
}
declare class RelationshipClient {
    readonly httpClient: HttpClient<Relationship>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Relationship>>;
    Create(r: Relationship): Observable<AjaxResponse<Relationship>>;
    Update(r: Relationship): Observable<AjaxResponse<Relationship>>;
    List(lo: RelationshipListOptions): Observable<AjaxResponse<RelationshipList>>;
    Delete(id: string): Observable<AjaxResponse<Relationship>>;
}
declare class PartyTypeClient {
    readonly httpClient: HttpClient<PartyType>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<PartyType>>;
    Create(p: PartyType): Observable<AjaxResponse<PartyType>>;
    Update(p: PartyType): Observable<AjaxResponse<PartyType>>;
    List(lo: PartyTypeListOptions): Observable<AjaxResponse<PartyTypeList>>;
}
declare class PartyAttributeDefinitionClient {
    readonly httpClient: HttpClient<PartyAttributeDefinition>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<PartyAttributeDefinition>>;
    Create(a: PartyAttributeDefinition): Observable<AjaxResponse<PartyAttributeDefinition>>;
    Update(a: PartyAttributeDefinition): Observable<AjaxResponse<PartyAttributeDefinition>>;
    List(lo: PartyAttributeDefinitionListOptions): Observable<AjaxResponse<PartyAttributeDefinitionList>>;
}
declare class NationalityClient {
    readonly httpClient: HttpClient<Nationality>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Nationality>>;
    Create(n: Nationality): Observable<AjaxResponse<Nationality>>;
    Update(n: Nationality): Observable<AjaxResponse<Nationality>>;
    List(lo: NationalityListOptions): Observable<AjaxResponse<NationalityList>>;
}
declare class MembershipClient {
    readonly httpClient: HttpClient<Membership>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Membership>>;
    Create(m: Membership): Observable<AjaxResponse<Membership>>;
    Update(m: Membership): Observable<AjaxResponse<Membership>>;
    List(lo: MembershipListOptions): Observable<AjaxResponse<MembershipList>>;
}
declare class IndividualClient {
    readonly httpClient: HttpClient<Individual>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Individual>>;
    Create(i: Individual): Observable<AjaxResponse<Individual>>;
    Update(i: Individual): Observable<AjaxResponse<Individual>>;
    List(lo: IndividualListOptions): Observable<AjaxResponse<IndividualList>>;
}
declare class IdentificationDocumentTypeClient {
    readonly httpClient: HttpClient<IdentificationDocumentType>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<IdentificationDocumentType>>;
    Create(i: IdentificationDocumentType): Observable<AjaxResponse<IdentificationDocumentType>>;
    Update(i: IdentificationDocumentType): Observable<AjaxResponse<IdentificationDocumentType>>;
    List(lo: IdentificationDocumentTypeListOptions): Observable<AjaxResponse<IdentificationDocumentTypeList>>;
}
declare class IdentificationDocumentClient {
    readonly httpClient: HttpClient<IdentificationDocument>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<IdentificationDocument>>;
    Create(i: IdentificationDocument): Observable<AjaxResponse<IdentificationDocument>>;
    Update(i: IdentificationDocument): Observable<AjaxResponse<IdentificationDocument>>;
    List(lo: IdentificationDocumentListOptions): Observable<AjaxResponse<IdentificationDocumentList>>;
    Delete(id: string): Observable<AjaxResponse<IdentificationDocument>>;
}
declare class CountryClient {
    readonly httpClient: HttpClient<Country>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Country>>;
    Create(c: Country): Observable<AjaxResponse<Country>>;
    Update(c: Country): Observable<AjaxResponse<Country>>;
    List(lo: CountryListOptions): Observable<AjaxResponse<CountryList>>;
}
declare class PartyClient {
    readonly httpClient: HttpClient<Party>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Party>>;
    Create(p: Party): Observable<AjaxResponse<Party>>;
    Update(p: Party): Observable<AjaxResponse<Party>>;
    List(lo: PartyListOptions): Observable<AjaxResponse<PartyList>>;
    Search(so: PartySearchOptions): Observable<AjaxResponse<PartyList>>;
}
export declare class IAMClient {
    host: string;
    scheme: string;
    constructor(host: string, scheme: string);
    Parties(): PartyClient;
    Countries(): CountryClient;
    IdentificationDocuments(): IdentificationDocumentClient;
    IdentificationDocumentTypes(): IdentificationDocumentTypeClient;
    Individuals(): IndividualClient;
    Memberships(): MembershipClient;
    Nationalities(): NationalityClient;
    PartyAttributeDefinitions(): PartyAttributeDefinitionClient;
    PartyTypes(): PartyTypeClient;
    Relationships(): RelationshipClient;
    RelationshipTypes(): RelationshipTypeClient;
    Teams(): TeamClient;
}
declare class CaseClient {
    readonly httpClient: HttpClient<Case>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Case>>;
    Create(c: Case): Observable<AjaxResponse<Case>>;
    Update(c: Case): Observable<AjaxResponse<Case>>;
    List(lo: CaseListOptions): Observable<AjaxResponse<CaseList>>;
}
declare class CaseTypeClient {
    readonly httpClient: HttpClient<CaseType>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<CaseType>>;
    Create(c: CaseType): Observable<AjaxResponse<CaseType>>;
    Update(c: CaseType): Observable<AjaxResponse<CaseType>>;
    List(lo: CaseTypeListOptions): Observable<AjaxResponse<CaseTypeList>>;
}
declare class CommentClient {
    readonly httpClient: HttpClient<Comment>;
    endpoint: string;
    constructor(host: string, scheme: string);
    Get(id: string): Observable<AjaxResponse<Comment>>;
    Create(c: Comment): Observable<AjaxResponse<Comment>>;
    Update(c: Comment): Observable<AjaxResponse<Comment>>;
    List(lo?: CommentListOptions): Observable<AjaxResponse<CommentList>>;
    Delete(id: string): Observable<AjaxResponse<Comment>>;
}
export declare class CMSClient {
    host: string;
    scheme: string;
    constructor(host: string, scheme: string);
    Cases(): CaseClient;
    CaseTypes(): CaseTypeClient;
    Comments(): CommentClient;
}
export {};
