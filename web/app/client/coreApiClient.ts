import { map, Observable, of, OperatorFunction, Subject } from 'rxjs';
import { switchMap } from 'rxjs/operators';
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
  PartySearchOptions,
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
import { ajax, AjaxResponse } from 'rxjs/ajax';
import { XMLHttpRequest } from 'xhr2';

// needed for rxjs/ajax compatibility outside the browser
global.XMLHttpRequest = global.XMLHttpRequest ? global.XMLHttpRequest : XMLHttpRequest;

// todo: should come from environment
const shouldAddAuthHeader = true;

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

// TODO
function isErrorResponse(data: any): boolean {
  return false;
}

class Client {
  private _host: string;
  private _headers: Headers;

  verb(verb: string): Request {
    return new Request(this);
  }

  get(): Request {
    return this.verb('GET');
  }

  post(): Request {
    return this.verb('POST');
  }

  put(): Request {
    return this.verb('PUT');
  }

  delete(): Request {
    return this.verb('DELETE');
  }

  do(): OperatorFunction<Request, Response> {
    return source => {
      return source.pipe(
        switchMap(req => {

          let url = this._host;
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

          //TODO : append query parameters

          //TODO: replace path params

          return ajax(
            {
              url: url,
              headers: headers,
              method: req._verb,
              async: true,
              timeout: 0,
              crossDomain: true,
              withCredentials: false,
              body: req._body
            }
          );
        }),
        map(ajaxResponse => {
          if (ajaxResponse.status > 399) {
            if (isErrorResponse(ajaxResponse.response)) {
              return new Response(ajaxResponse.response);
            } else {
              return new Response({ error: 'server error', status: 500 });
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

interface Headers {
  [key: string]: string[];
}

class Response {
  public constructor(readonly body: any) {

  }

  as<T>(): T {
    return this.body as T;
  }

}


// http:google.com/:pathParam?queryParam=1
class Request {
  public _client: Client;
  public _error: Error;
  public _path: string;
  public _verb: string;
  public _body: any;
  public _params: URLValues;
  public _pathParams: URLValues;
  public _headers: Headers;

  public constructor(client: Client) {
    this._client = client;
  }

  public verb(verb: string): Request {
    this._verb = verb;
    return this;
  }

  public get(): Request {
    return this.verb('GET');
  }

  public put(): Request {
    return this.verb('PUT');
  }

  public post(): Request {
    return this.verb('POST');
  }

  public delete(): Request {
    return this.verb('DELETE');
  }

  public path(path: string): Request {
    this._path = path;
    return this;
  }

  public body(body: any): Request {
    this._body = body;
    return this;
  }

  public params(params: URLValues): Request {
    this._params = params;
    return this;
  }

  public pathParam(key: string, value: string): Request {
    if (!this._pathParams) {
      this._pathParams = {};
    }
    this._pathParams[key] = value;
    return this;
  }

  public headers(headers: Headers): Request {
    this._headers = headers;
    return this;
  }

}

function responseAs<T>(): OperatorFunction<Response, T> {
  return source => {
    return source.pipe(map(resp => {
      if (isErrorResponse(resp.body)) {
        throw new Error(resp.body);
      } else {
        return resp.body;
      }
    }));
  };
}

class PartyClient2 {

  private client: Client;

  execute: () => OperatorFunction<Request, Response>;

  public constructor(client: Client) {
    this.execute = client.do;
  }

  get(): OperatorFunction<string, Party> {
    return id$ => id$.pipe(
      map(id => this.client.get().path('/apis/iam/v1/parties/:id').pathParam('id', id)),
      this.execute(),
      responseAs<Party>()
    );
  }

  put(): OperatorFunction<Party, Party> {
    return party$ => party$.pipe(
      map(party => this.client.put().path('/apis/iam/v1/parties/:id').pathParam('id', party.id)),
      this.execute(),
      responseAs<Party>()
    );
  }

}

//
// function myTest() {
//
//   var client = new Client();
//   var partyClient = new PartyClient2(client);
//   partyClient.execute = () => {
//     return s => {
//       return of(new Response({ id: 'abcdef' }));
//     };
//   };
//   partyClient.execute = () => {
//     return s => {
//       return of(new Response({ status: 500, error: 'abcdef' }));
//     };
//   };
//   partyClient.execute = () => {
//     return s => {
//       s.subscribe(r => {
//
//       });
//     };
//   };
//
// }


// const myClient = new PartyClient2();
// const subject = new Subject<string>();
// const party$ = subject.pipe(myClient.get()).subscribe({
//     next: party => {
//
//     },
//     error: err => {
//
//     },
//     complete: () => {
//     }
//   }
// );
//
// subject.next('mynewid');


class HttpClient<T> {
  headers: { 'X-Authenticated-User-Subject'?: string; };

  constructor(shouldAddAuthPassthroughHeader: boolean) {
    this.headers = {};
    if (shouldAddAuthPassthroughHeader) {
      this.headers = {
        'X-Authenticated-User-Subject': 'stephen.kabagambe@email.com'
      };
    }
  }

  get(url: string): Observable<AjaxResponse<T>> {
    return ajax(
      {
        url: url,
        headers: this.headers,
        method: 'GET',
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        responseType: 'json'
      }
    );
  }

  getCustom<R>(url: string): Observable<AjaxResponse<R>> {
    return ajax(
      {
        url: url,
        headers: this.headers,
        method: 'GET',
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        responseType: 'json'
      }
    );
  }

  put(url: string, body: T): Observable<AjaxResponse<T>> {
    return ajax(
      {
        url: url,
        body: body,
        headers: this.headers,
        method: 'PUT',
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        responseType: 'json'
      }
    );
  }

  post(url: string, body: T): Observable<AjaxResponse<T>> {
    return ajax(
      {
        url: url,
        body: body,
        headers: this.headers,
        method: 'POST',
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        responseType: 'json'
      }
    );
  }

  postCustom<B, R>(url: string, body: B): Observable<AjaxResponse<R>> {
    return ajax(
      {
        url: url,
        body: body,
        headers: this.headers,
        method: 'POST',
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        responseType: 'json'
      }
    );
  }

  delete(url: string): Observable<AjaxResponse<T>> {
    return ajax(
      {
        url: url,
        headers: this.headers,
        method: 'DELETE',
        async: true,
        timeout: 0,
        crossDomain: true,
        withCredentials: false,
        responseType: 'json'
      }
    );
  }
}

// -- IAM ---------------------------

class TeamClient {
  readonly httpClient = new HttpClient<Team>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/teams`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(t: Team) {
    return this.httpClient.post(this.endpoint, t);
  }

  Update(t: Team) {
    return this.httpClient.put(this.endpoint + '/' + t.id, t);
  }

  List(lo: TeamListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<TeamList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class RelationshipTypeClient {
  readonly httpClient = new HttpClient<RelationshipType>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/relationshiptypes`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(r: RelationshipType) {
    return this.httpClient.post(this.endpoint, r);
  }

  Update(r: RelationshipType) {
    return this.httpClient.put(this.endpoint + '/' + r.id, r);
  }

  List(lo: RelationshipTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<RelationshipTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class RelationshipClient {
  readonly httpClient = new HttpClient<Relationship>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/relationships`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(r: Relationship) {
    return this.httpClient.post(this.endpoint, r);
  }

  Update(r: Relationship) {
    return this.httpClient.put(this.endpoint + '/' + r.id, r);
  }

  List(lo: RelationshipListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<RelationshipList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + '/' + id);
  }
}

class PartyTypeClient {
  readonly httpClient = new HttpClient<PartyType>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/partytypes`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(p: PartyType) {
    return this.httpClient.post(this.endpoint, p);
  }

  Update(p: PartyType) {
    return this.httpClient.put(this.endpoint + '/' + p.id, p);
  }

  List(lo: PartyTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<PartyTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class PartyAttributeDefinitionClient {
  readonly httpClient = new HttpClient<PartyAttributeDefinition>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/attributes`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(a: PartyAttributeDefinition) {
    return this.httpClient.post(this.endpoint, a);
  }

  Update(a: PartyAttributeDefinition) {
    return this.httpClient.put(this.endpoint + '/' + a.id, a);
  }

  List(lo: PartyAttributeDefinitionListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<PartyAttributeDefinitionList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class NationalityClient {
  readonly httpClient = new HttpClient<Nationality>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/nationalities`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(n: Nationality) {
    return this.httpClient.post(this.endpoint, n);
  }

  Update(n: Nationality) {
    return this.httpClient.put(this.endpoint + '/' + n.id, n);
  }

  List(lo: NationalityListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<NationalityList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class MembershipClient {
  readonly httpClient = new HttpClient<Membership>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/memberships`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(m: Membership) {
    return this.httpClient.post(this.endpoint, m);
  }

  Update(m: Membership) {
    return this.httpClient.put(this.endpoint + '/' + m.id, m);
  }

  List(lo: MembershipListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<MembershipList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class IndividualClient {
  readonly httpClient = new HttpClient<Individual>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/individuals`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(i: Individual) {
    return this.httpClient.post(this.endpoint, i);
  }

  Update(i: Individual) {
    return this.httpClient.put(this.endpoint + '/' + i.id, i);
  }

  List(lo: IndividualListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<IndividualList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class IdentificationDocumentTypeClient {
  readonly httpClient = new HttpClient<IdentificationDocumentType>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/identificationdocumenttypes`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(i: IdentificationDocumentType) {
    return this.httpClient.post(this.endpoint, i);
  }

  Update(i: IdentificationDocumentType) {
    return this.httpClient.put(this.endpoint + '/' + i.id, i);
  }

  List(lo: IdentificationDocumentTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<IdentificationDocumentTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class IdentificationDocumentClient {
  readonly httpClient = new HttpClient<IdentificationDocument>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/identificationdocuments`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(i: IdentificationDocument) {
    return this.httpClient.post(this.endpoint, i);
  }

  Update(i: IdentificationDocument) {
    return this.httpClient.put(this.endpoint + '/' + i.id, i);
  }

  List(lo: IdentificationDocumentListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<IdentificationDocumentList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + '/' + id);
  }
}

class CountryClient {
  readonly httpClient = new HttpClient<Country>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/countries`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(c: Country) {
    return this.httpClient.post(this.endpoint, c);
  }

  Update(c: Country) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c);
  }

  List(lo: CountryListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<CountryList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class PartyClient {
  readonly httpClient = new HttpClient<Party>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/iam/v1/parties`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(p: Party) {
    return this.httpClient.post(this.endpoint, p);
  }

  Update(p: Party) {
    return this.httpClient.put(this.endpoint + '/' + p.id, p);
  }

  List(lo: PartyListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<PartyList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }

  Search(so: PartySearchOptions) {
    return this.httpClient.postCustom<PartySearchOptions, PartyList>(this.endpoint + '/search', so);
  }
}

export class IAMClient {
  host: string;
  scheme: string;

  constructor(host: string, scheme: string) {
    this.host = host;
    this.scheme = scheme;
  }

  Parties() {
    return new PartyClient(this.host, this.scheme);
  }

  Countries() {
    return new CountryClient(this.host, this.scheme);
  }

  IdentificationDocuments() {
    return new IdentificationDocumentClient(this.host, this.scheme);
  }

  IdentificationDocumentTypes() {
    return new IdentificationDocumentTypeClient(this.host, this.scheme);
  }

  Individuals() {
    return new IndividualClient(this.host, this.scheme);
  }

  Memberships() {
    return new MembershipClient(this.host, this.scheme);
  }

  Nationalities() {
    return new NationalityClient(this.host, this.scheme);
  }

  PartyAttributeDefinitions() {
    return new PartyAttributeDefinitionClient(this.host, this.scheme);
  }

  PartyTypes() {
    return new PartyTypeClient(this.host, this.scheme);
  }

  Relationships() {
    return new RelationshipClient(this.host, this.scheme);
  }

  RelationshipTypes() {
    return new RelationshipTypeClient(this.host, this.scheme);
  }

  Teams() {
    return new TeamClient(this.host, this.scheme);
  }
}

// -- CMS ---------------------------

class CaseClient {
  readonly httpClient = new HttpClient<Case>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/cms/v1/cases`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(c: Case) {
    return this.httpClient.post(this.endpoint, c);
  }

  Update(c: Case) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c);
  }

  List(lo: CaseListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<CaseList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class CaseTypeClient {
  readonly httpClient = new HttpClient<CaseType>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/cms/v1/casetypes`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(c: CaseType) {
    return this.httpClient.post(this.endpoint, c);
  }

  Update(c: CaseType) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c);
  }

  List(lo: CaseTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<CaseTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }
}

class CommentClient {
  readonly httpClient = new HttpClient<Comment>(shouldAddAuthHeader);
  endpoint: string;

  constructor(host: string, scheme: string) {
    this.endpoint = `${scheme}://${host}/apis/cms/v1/comments`;
  }

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id);
  }

  Create(c: Comment) {
    return this.httpClient.post(this.endpoint, c);
  }

  Update(c: Comment) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c);
  }

  List(lo?: CommentListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>);
    return this.httpClient.getCustom<CommentList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    );
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + '/' + id);
  }
}

export class CMSClient {
  host: string;
  scheme: string;

  constructor(host: string, scheme: string) {
    this.host = host;
    this.scheme = scheme;
  }

  Cases() {
    return new CaseClient(this.host, this.scheme);
  }

  CaseTypes() {
    return new CaseTypeClient(this.host, this.scheme);
  }

  Comments() {
    return new CommentClient(this.host, this.scheme);
  }
}
