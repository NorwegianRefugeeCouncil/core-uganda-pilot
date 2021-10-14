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
import { of, Subject } from 'rxjs';
import { switchMap } from 'rxjs/operators'
import { IAMClient, Headers, Response } from './coreApiClient';

const defaultHost = "localhost:9000"
const defaultScheme = "http"
const defaultHeaders: Headers = {
  'X-Authenticated-User-Subject': ['stephen.kabagambe@email.com']
}

describe('IAMClient - Parties', () => {
  let subject
  let iamClient
  let partyClient
  before(() => {
    subject = new Subject<string>()
    iamClient = new IAMClient("http", "localhost:9000", {
      'X-Authenticated-User-Subject': ['stephen.kabagambe@email.com']
    })
    partyClient = iamClient.Parties()
    partyClient.execute = () => {
      return s => {
        console.log("Custom PartyClient Execute", s)
        s.pipe(switchMap(req => {}))
        return of(new Response({ id: 'abcdef' }));
      };
    }
  })
  it('Should return expected value for get request', (done) => {
    subject.pipe(partyClient.Get()).subscribe((party) => {
      console.log(party)
      done()
    })
    subject.next('78b494dc-7461-42f5-bf2d-1c9695e63ba8')
  })
})

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
//
// }
//
// const client = new Client()
// const myClient = new PartyClient2(client);
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