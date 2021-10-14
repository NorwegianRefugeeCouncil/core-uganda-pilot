import {
  Form,
  Case,
  CaseList,
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
  TeamListOptions,
  CommentList,
  CaseTypeList, Control
} from './types/models';
import { expect } from 'chai'
import { of, Subject, map, OperatorFunction } from 'rxjs';
import { switchMap } from 'rxjs/operators'
import { IAMClient, Headers, Response, Request, prepareRequestOptions, AjaxRequestOptions } from './coreApiClient';

const defaults = {
    global: {
        scheme: 'testscheme',
        host: 'testhost',
        headers: {
            'X-Authenticated-User-Subject': ['test@user.email']
        }
    },
    parties: {
        party: {
            id: 'TESTPARTYID',
            partyTypeIds: [
                'TESTPARTYTYPEID'
            ],
            attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
        },
        partyGetRequestOptions: {
            url: 'http://localhost:9000/apis/iam/v1/parties/TESTPARTYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partyUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/parties/TESTPARTYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        }
    }
}

describe('IAMClient - Parties', () => {
  let subject
  let iamClient
  let partyClient
  let partyId
  let party

  // sets up test variables and injects mock execute function
  const init = (execute: () => OperatorFunction<Request, Response>) => {
      partyId = defaults.parties.party.id
      subject = new Subject<string>()
      iamClient = new IAMClient(
          defaults.global.scheme,
          defaults.global.host,
          defaults.global.headers
      )
      partyClient = iamClient.Parties()
      partyClient.execute = execute
  }

  // GET
  before(() => {
    init((): OperatorFunction<Request, Response> => {
        return source => {
            return source.pipe(
                switchMap(req => {
                    let ro = prepareRequestOptions(req)
                    expect(ro).to.deep.equal(defaults.parties.partyGetRequestOptions)
                    return of(new Response({ id: 'abcdef' }))
                })
            )
        }
    })
  })
  it('should correctly build GET request', (done) => {
    subject.pipe(partyClient.Get()).subscribe((party) => {done()})
    subject.next(partyId)
  })
})
