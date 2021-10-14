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
  for(const [key, value] of Object.entries(params)){
    tmp = tmp.replace(`:${key}`, value)
  }
  return tmp
}

function appendQueryParams(s: string, params: URLValues): string {
  let paramStrings = []
  for(const [key, value] of Object.entries(params)){
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
  private readonly _scheme: string
  private readonly _host: string
  private _headers: Headers

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

class Request {
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

class PartyClient {

  private client: Client;

  execute: () => OperatorFunction<Request, Response>;

  public constructor(client: Client) {
    this.client = client
    this.execute = client.do
  }

  Get(): OperatorFunction<string, Party> {
    return id$ => id$.pipe(
      map(id => this.client.get().path('/apis/iam/v1/parties/:id').pathParam('id', id)),
      this.execute(),
      responseAs<Party>()
    )
  }

  Update(): OperatorFunction<Party, Party> {
    return party$ => party$.pipe(
      map(party => this.client.put().path('/apis/iam/v1/parties/:id').pathParam('id', party.id)),
      this.execute(),
      responseAs<Party>()
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

  public Parties(){
    return new PartyClient(this._client)
  }
}

