import { Observable } from 'rxjs';
import {
  Case,
  CaseListOptions,
  Comment,
  CaseType,
  CaseTypeListOptions,
  CommentListOptions,
  Party,
  PartySearchOptions,
  PartyListOptions,
  PartyList,
  Country,
  CountryListOptions,
  CountryList,
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
  PartyAttributeDefinition,
  PartyAttributeDefinitionList,
  PartyAttributeDefinitionListOptions,
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
  TeamListOptions, CommentList
} from './types/models';
import { ajax, AjaxResponse } from 'rxjs/ajax';
import { XMLHttpRequest } from 'xhr2';

// needed for rxjs/ajax compatibility outside the browser
global.XMLHttpRequest = global.XMLHttpRequest ? global.XMLHttpRequest : XMLHttpRequest;

// todo: should come from environment
const shouldAddAuthHeader = true;

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

class TeamClient {
  readonly httpClient = new HttpClient<Team>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/teams'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(t: Team) {
    return this.httpClient.post(this.endpoint, t)
  }

  Update(t: Team) {
    return this.httpClient.put(this.endpoint + '/' + t.id, t)
  }

  List(lo: TeamListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<TeamList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class RelationshipTypeClient {
  readonly httpClient = new HttpClient<RelationshipType>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/relationshiptypes'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(r: RelationshipType) {
    return this.httpClient.post(this.endpoint, r)
  }

  Update(r: RelationshipType) {
    return this.httpClient.put(this.endpoint + '/' + r.id, r)
  }

  List(lo: RelationshipTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<RelationshipTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class RelationshipClient {
  readonly httpClient = new HttpClient<Relationship>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/relationships'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(r: Relationship) {
    return this.httpClient.post(this.endpoint, r)
  }

  Update(r: Relationship) {
    return this.httpClient.put(this.endpoint + '/' + r.id, r)
  }

  List(lo: RelationshipListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<RelationshipList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + '/' + id)
  }
}

class PartyTypeClient {
  readonly httpClient = new HttpClient<PartyType>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/partytypes'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(p: PartyType) {
    return this.httpClient.post(this.endpoint, p)
  }

  Update(p: PartyType) {
    return this.httpClient.put(this.endpoint + '/' + p.id, p)
  }

  List(lo: PartyTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<PartyTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class PartyAttributeDefinitionClient {
  readonly httpClient = new HttpClient<PartyAttributeDefinition>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/attributes'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(a: PartyAttributeDefinition) {
    return this.httpClient.post(this.endpoint, a)
  }

  Update(a: PartyAttributeDefinition) {
    return this.httpClient.put(this.endpoint + '/' + a.id, a)
  }

  List(lo: PartyAttributeDefinitionListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<PartyAttributeDefinitionList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class NationalityClient {
  readonly httpClient = new HttpClient<Nationality>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/nationalities'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(n: Nationality) {
    return this.httpClient.post(this.endpoint, n)
  }

  Update(n: Nationality) {
    return this.httpClient.put(this.endpoint + '/' + n.id, n)
  }

  List(lo: NationalityListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<NationalityList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class MembershipClient {
  readonly httpClient = new HttpClient<Membership>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/memberships'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(m: Membership) {
    return this.httpClient.post(this.endpoint, m)
  }

  Update(m: Membership) {
    return this.httpClient.put(this.endpoint + '/' + m.id, m)
  }

  List(lo: MembershipListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<MembershipList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class IndividualClient {
  readonly httpClient = new HttpClient<Individual>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/individuals'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(i: Individual) {
    return this.httpClient.post(this.endpoint, i)
  }

  Update(i: Individual) {
    return this.httpClient.put(this.endpoint + '/' + i.id, i)
  }

  List(lo: IndividualListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<IndividualList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class IdentificationDocumentTypeClient {
  readonly httpClient = new HttpClient<IdentificationDocumentType>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/identificationdocumenttypes'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(i: IdentificationDocumentType) {
    return this.httpClient.post(this.endpoint, i)
  }

  Update(i: IdentificationDocumentType) {
    return this.httpClient.put(this.endpoint + '/' + i.id, i)
  }

  List(lo: IdentificationDocumentTypeListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<IdentificationDocumentTypeList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class IdentificationDocumentClient {
  readonly httpClient = new HttpClient<IdentificationDocument>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/identificationdocuments'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(i: IdentificationDocument) {
    return this.httpClient.post(this.endpoint, i)
  }

  Update(i: IdentificationDocument) {
    return this.httpClient.put(this.endpoint + '/' + i.id, i)
  }

  List(lo: IdentificationDocumentListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<IdentificationDocumentList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + '/' + id)
  }
}

class CountryClient {
  readonly httpClient = new HttpClient<Country>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/countries'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(c: Country) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: Country) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c)
  }

  List(lo: CountryListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<CountryList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }
}

class PartyClient {
  readonly httpClient = new HttpClient<Party>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/iam/v1/parties'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(p: Party) {
    return this.httpClient.post(this.endpoint, p)
  }

  Update(p: Party) {
    return this.httpClient.put(this.endpoint + '/' + p.id, p)
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

export class IAMClient {
  static Parties() {
    return new PartyClient
  }

  static Countries() {
    return new CountryClient
  }

  static IdentificationDocuments() {
    return new IdentificationDocumentClient
  }

  static IdentificationDocumentTypes() {
    return new IdentificationDocumentTypeClient
  }

  static Individuals() {
    return new IndividualClient
  }

  static Memberships() {
    return new MembershipClient
  }

  static Nationalities() {
    return new NationalityClient
  }

  static PartyAttributeDefinitions() {
    return new PartyAttributeDefinitionClient
  }

  static PartyTypes() {
    return new PartyTypeClient
  }

  static Relationships() {
    return new RelationshipClient
  }

  static RelationshipTypes() {
    return new RelationshipTypeClient
  }

  static Teams() {
    return new TeamClient
  }
}

// -- CMS ---------------------------

class CaseClient {
  readonly httpClient = new HttpClient<Case>(shouldAddAuthHeader)
  readonly endpoint = 'http://localhost:9000/apis/cms/v1/cases'

  Get(id: string) {
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(c: Case) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: Case) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c)
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
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(c: CaseType) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: CaseType) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c)
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
    return this.httpClient.get(this.endpoint + '/' + id)
  }

  Create(c: Comment) {
    return this.httpClient.post(this.endpoint, c)
  }

  Update(c: Comment) {
    return this.httpClient.put(this.endpoint + '/' + c.id, c)
  }

  List(lo?: CommentListOptions) {
    const query = new URLSearchParams(lo as Record<string, string>)
    return this.httpClient.getCustom<CommentList>(
      query ? this.endpoint : this.endpoint + `?${query}`
    )
  }

  Delete(id: string) {
    return this.httpClient.delete(this.endpoint + '/' + id)
  }
}

export class CMSClient {
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

