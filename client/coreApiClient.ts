import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import {
  Case, CaseListOptions, Comment,
  CaseType, CaseTypeListOptions, CommentListOptions, Party, PartySearchOptions, PartyListOptions, PartyList
} from './types/models';
import { ajax, AjaxResponse } from 'rxjs/ajax';
import { XMLHttpRequest } from 'xhr2';

// needed for rxjs/ajax compatibility outside the browser
global.XMLHttpRequest = global.XMLHttpRequest ? global.XMLHttpRequest : XMLHttpRequest;

// todo: should come from environment
const shouldAddAuthHeader = true;
const noop = () => {}

class HttpClient<T> {
  headers: { "X-Authenticated-User-Subject"?: string; }

  constructor(shouldAddAuthPassthroughHeader: boolean) {
    this.headers = {}
    if (shouldAddAuthPassthroughHeader) {
      this.headers = {
        'X-Authenticated-User-Subject': 'stephen.kabagambe@email.com'
      }
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
    )
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
    )
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
    )
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
    )
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
    )
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
    )
  }
}

// -- IAM ---------------------------

class PartyClient {
  readonly httpClient = new HttpClient<Party>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/parties'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + id)
  }

  Create(p: Party) {
    return this.httpClient.post(this.endpoint, p)
  }

  Update(p: Party) {
    return this.httpClient.put(this.endpoint + p.id, p)
  }

  List(lo: PartyListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<PartyList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }

  Search(so: PartySearchOptions) {
    return this.httpClient.postCustom<PartySearchOptions, PartyList>(this.endpoint + '/search', so)
  }
}

class IAMClient {
  static Parties() {
    return new PartyClient
  }
}

// -- CMS ---------------------------

class CaseClient {
  readonly httpClient = new HttpClient<Case>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/cms/v1/cases'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + id)
  }

  Create(c: Case) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: Case) {
    return this.httpClient.put(this.endpoint + c.id, c)
  }

  List(lo: CaseListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.get(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class CaseTypeClient {
  readonly httpClient = new HttpClient<CaseType>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/cms/v1/casetypes'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + id)
  }

  Create(c: CaseType) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: CaseType) {
    return this.httpClient.put(this.endpoint + c.id, c)
  }

  List(lo: CaseTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.get(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class CommentClient {
  readonly httpClient = new HttpClient<Comment>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/cms/v1/comments'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + id)
  }

  Create(c: Comment) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: Comment) {
    return this.httpClient.put(this.endpoint + c.id, c)
  }

  List(lo: CommentListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.get(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + id)
  }
}

class CMSClient {
  static Cases() {
    return new CaseClient
  }

  static CaseTypes() {
    return new CaseTypeClient
  }

  static Comments() {
    return new CommentClient
  }
}

CMSClient.Comments().List(new CommentListOptions())
  .pipe(map((response) => {
    console.log(response);
  }, catchError(error => {
    console.log('error: ', error);
    return of(error);
  }))
).subscribe(noop);

