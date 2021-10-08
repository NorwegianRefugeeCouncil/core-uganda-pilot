import { IAMClient, CMSClient } from './coreApiClient'
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

const defaultHost = "localhost:9000"
const defaultScheme = "http"

const cmsClient = new CMSClient(defaultHost, defaultScheme)
const iamClient = new IAMClient(defaultHost, defaultScheme)

// -- CMS ---------------------------

describe("CMSClient - Comments", () => {
  let testId: string
  let testComment = {
    caseId: "TESTCASEID"
  } as Comment

  it("Should return 200 on Create", (done) => {
    cmsClient.Comments().Create(testComment)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    cmsClient.Comments().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotComment = response.response as Comment
          expect(gotComment.caseId).to.equal(testComment.caseId)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testComment = {
      id: testId,
      caseId: "TESTCASEID2"
    } as Comment

    cmsClient.Comments().Update({
      id: testId,
      caseId: "TESTCASEID2"
    } as Comment)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    cmsClient.Comments().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotComment = response.response as Comment
          expect(gotComment.caseId).to.equal(testComment.caseId)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    cmsClient.Comments().List({} as CommentListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCommentList = response.response as CommentList
          let listContainsComment = false
          gotCommentList.items.forEach((comment) => {
            if (comment.id == testId) {
              expect(comment.caseId).to.equal(testComment.caseId)
              listContainsComment = true
            }
          })
          expect(listContainsComment).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Delete", (done) => {
    cmsClient.Comments().Delete(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})

describe("CMSClient - Case Types", () => {
  let testId: string
  let testCaseType = {
    name: "TESTCASETYPENAME",
    partyTypeId: "TESTPARTYTYPEID",
    teamId: "TESTPARTYTYPEID",
    form: {} as Form,
    intakeCaseType: false
  } as CaseType

  it("Should return 200 on Create", (done) => {
    cmsClient.CaseTypes().Create(testCaseType)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    cmsClient.CaseTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCaseType = response.response as CaseType
          expect(gotCaseType.name).to.equal(testCaseType.name)
          expect(gotCaseType.partyTypeId).to.equal(testCaseType.partyTypeId)
          expect(gotCaseType.teamId).to.equal(testCaseType.teamId)
          expect(gotCaseType.intakeCaseType).to.equal(testCaseType.intakeCaseType)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testCaseType = {
      name: "TESTCASETYPENAME2",
      partyTypeId: "TESTPARTYTYPEID2",
      teamId: "TESTPARTYTYPEID2",
      form: {} as Form,
      intakeCaseType: true
    } as CaseType

    cmsClient.CaseTypes().Update({
      name: "TESTCASETYPENAME2",
      partyTypeId: "TESTPARTYTYPEID2",
      teamId: "TESTPARTYTYPEID2",
      form: {} as Form,
      intakeCaseType: true
    } as CaseType)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    cmsClient.CaseTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCaseType = response.response as CaseType
          expect(gotCaseType.name).to.equal(testCaseType.name)
          expect(gotCaseType.partyTypeId).to.equal(testCaseType.partyTypeId)
          expect(gotCaseType.teamId).to.equal(testCaseType.teamId)
          expect(gotCaseType.intakeCaseType).to.equal(testCaseType.intakeCaseType)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    cmsClient.CaseTypes().List({} as CaseTypeListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCaseTypeList = response.response as CaseTypeList
          let listContainsCaseType = false
          gotCaseTypeList.items.forEach((caseType) => {
            if (caseType.id == testId) {
              expect(caseType.name).to.equal(testCaseType.name)
              expect(caseType.partyTypeId).to.equal(testCaseType.partyTypeId)
              expect(caseType.teamId).to.equal(testCaseType.teamId)
              expect(caseType.intakeCaseType).to.equal(testCaseType.intakeCaseType)
              listContainsCaseType = true
            }
          })
          expect(listContainsCaseType).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})

describe("CMSClient - Cases", () => {
  let testId: string
  let testCase = {
    caseTypeId: "TESTCASETYPEID",
    partyId: "TESTPARTYID",
    teamId: "TESTTEAMID",
    creatorId: "TESTCREATORID",
    parentId: "TESTPARENTID",
    intakeCase: false,
    form: {} as Form,
    formData: {} as {[key: string]: string[]},
    done: false
  } as Case

  it("Should return 200 on Create", (done) => {
    cmsClient.Cases().Create(testCase)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    cmsClient.Cases().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCase = response.response as Case
          expect(gotCase.caseTypeId).to.equal(testCase.caseTypeId)
          expect(gotCase.partyId).to.equal(testCase.partyId)
          expect(gotCase.teamId).to.equal(testCase.teamId)
          expect(gotCase.creatorId).to.equal(testCase.creatorId)
          expect(gotCase.parentId).to.equal(testCase.parentId)
          expect(gotCase.intakeCase).to.equal(testCase.intakeCase)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testCase = {
      caseTypeId: "TESTCASETYPEID2",
      partyId: "TESTPARTYID2",
      teamId: "TESTTEAMID2",
      creatorId: "TESTCREATORID2",
      parentId: "TESTPARENTID2",
      intakeCase: true,
      form: {} as Form,
      formData: {} as {[key: string]: string[]},
      done: true
    } as Case

    cmsClient.Cases().Update({
      caseTypeId: "TESTCASETYPEID2",
      partyId: "TESTPARTYID2",
      teamId: "TESTTEAMID2",
      creatorId: "TESTCREATORID2",
      parentId: "TESTPARENTID2",
      intakeCase: true,
      form: {} as Form,
      formData: {} as {[key: string]: string[]},
      done: true
    } as Case)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    cmsClient.Cases().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCase = response.response as Case
          expect(gotCase.caseTypeId).to.equal(testCase.caseTypeId)
          expect(gotCase.partyId).to.equal(testCase.partyId)
          expect(gotCase.teamId).to.equal(testCase.teamId)
          expect(gotCase.creatorId).to.equal(testCase.creatorId)
          expect(gotCase.parentId).to.equal(testCase.parentId)
          expect(gotCase.intakeCase).to.equal(testCase.intakeCase)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    cmsClient.Cases().List({} as CaseListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotCaseList = response.response as CaseList
          let listContainsCase = false
          gotCaseList.items.forEach((kase) => {
            if (kase.id == testId) {
              expect(kase.caseTypeId).to.equal(testCase.caseTypeId)
              expect(kase.partyId).to.equal(testCase.partyId)
              expect(kase.teamId).to.equal(testCase.teamId)
              expect(kase.creatorId).to.equal(testCase.creatorId)
              expect(kase.parentId).to.equal(testCase.parentId)
              expect(kase.intakeCase).to.equal(testCase.intakeCase)
              listContainsCase = true
            }
          })
          expect(listContainsCase).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})

// -- IAM ---------------------------

describe("IAMClient - Teams", () => {
  let testId: string
  let testTeam = {
    name: "TESTTEAMNAME"
  } as Team

  it("Should return 200 on Create", (done) => {
    iamClient.Teams().Create(testTeam)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    iamClient.Teams().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotTeam = response.response as Team
          expect(gotTeam.name).to.equal(testTeam.name)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testTeam = {
        name: "TESTTEAMNAME2"
    } as Team

    iamClient.Teams().Update({
      name: "TESTTEAMNAME2"
    } as Team)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    iamClient.Teams().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotTeam = response.response as Team
          expect(gotTeam.name).to.equal(testTeam.name)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    iamClient.Teams().List({} as TeamListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotTeamList = response.response as TeamList
          let listContainsTeam = false
          gotTeamList.items.forEach((team) => {
            if (team.id == testId) {
              expect(team.name).to.equal(testTeam.name)
              listContainsTeam = true
            }
          })
          expect(listContainsTeam).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})

describe("IAMClient - Relationship Types", () => {
  let testId: string
  let testRelationshipType = {
    isDirectional: true,
    name: "TESTRELATIONSHIPTYPE",
    firstPartyRole: "TESTFIRSTPARTYROLE",
    secondPartyRole: "TESTSECONDPARTYROLE"
  } as RelationshipType

  it("Should return 200 on Create", (done) => {
    iamClient.RelationshipTypes().Create(testRelationshipType)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    iamClient.RelationshipTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationshipType = response.response as RelationshipType
          expect(gotRelationshipType.isDirectional).to.equal(testRelationshipType.isDirectional)
          expect(gotRelationshipType.name).to.equal(testRelationshipType.name)
          expect(gotRelationshipType.firstPartyRole).to.equal(testRelationshipType.firstPartyRole)
          expect(gotRelationshipType.secondPartyRole).to.equal(testRelationshipType.secondPartyRole)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testRelationshipType = {
      isDirectional: false,
      name: "TESTRELATIONSHIPTYPE2",
      firstPartyRole: "TESTNONDIRECTIONAL",
      secondPartyRole: "TESTNONDIRECTIONAL"
    } as RelationshipType

    iamClient.RelationshipTypes().Update({
      isDirectional: false,
      name: "TESTRELATIONSHIPTYPE2",
      firstPartyRole: "TESTNONDIRECTIONAL",
      secondPartyRole: "TESTNONDIRECTIONAL"
    } as RelationshipType)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    iamClient.RelationshipTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationshipType = response.response as RelationshipType
          expect(gotRelationshipType.isDirectional).to.equal(testRelationshipType.isDirectional)
          expect(gotRelationshipType.name).to.equal(testRelationshipType.name)
          expect(gotRelationshipType.firstPartyRole).to.equal(testRelationshipType.firstPartyRole)
          expect(gotRelationshipType.secondPartyRole).to.equal(testRelationshipType.secondPartyRole)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    iamClient.RelationshipTypes().List({} as RelationshipTypeListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationshipTypeList = response.response as RelationshipTypeList
          let listContainsRelationshipType = false
          gotRelationshipTypeList.items.forEach((relationshipType) => {
            if (relationshipType.id == testId) {
              expect(relationshipType.isDirectional).to.equal(testRelationshipType.isDirectional)
              expect(relationshipType.name).to.equal(testRelationshipType.name)
              expect(relationshipType.firstPartyRole).to.equal(testRelationshipType.firstPartyRole)
              expect(relationshipType.secondPartyRole).to.equal(testRelationshipType.secondPartyRole)
              listContainsRelationshipType = true
            }
          })
          expect(listContainsRelationshipType).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})

describe("IAMClient - Relationships", () => {
  let testId: string
  let testRelationship = {
    relationshipTypeId: "TESTRELATIONSHIPTYPEID",
    firstParty: "TESTFIRSTPARTY",
    secondParty: "TESTSECONDPARTY"
  } as Relationship

  it("Should return 200 on Create", (done) => {
    iamClient.Relationships().Create(testRelationship)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    iamClient.Relationships().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationship = response.response as Relationship
          expect(gotRelationship.relationshipTypeId).to.equal(testRelationship.relationshipTypeId)
          expect(gotRelationship.firstParty).to.equal(testRelationship.firstParty)
          expect(gotRelationship.secondParty).to.equal(testRelationship.secondParty)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testRelationship = {
      relationshipTypeId: "TESTRELATIONSHIPTYPEID2",
      firstParty: "TESTFIRSTPARTY2",
      secondParty: "TESTSECONDPARTY2"
    } as Relationship

    iamClient.Relationships().Update({
      relationshipTypeId: "TESTRELATIONSHIPTYPEID2",
      firstParty: "TESTFIRSTPARTY2",
      secondParty: "TESTSECONDPARTY2"
    } as Relationship)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    iamClient.Relationships().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationship = response.response as Relationship
          expect(gotRelationship.relationshipTypeId).to.equal(testRelationship.relationshipTypeId)
          expect(gotRelationship.firstParty).to.equal(testRelationship.firstParty)
          expect(gotRelationship.secondParty).to.equal(testRelationship.secondParty)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    iamClient.Relationships().List({} as RelationshipListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationshipList = response.response as RelationshipList
          let listContainsRelationship = false
          gotRelationshipList.items.forEach((relationship) => {
            if (relationship.id == testId) {
              expect(relationship.relationshipTypeId).to.equal(testRelationship.relationshipTypeId)
              expect(relationship.firstParty).to.equal(testRelationship.firstParty)
              expect(relationship.secondParty).to.equal(testRelationship.secondParty)
              listContainsRelationship = true
            }
          })
          expect(listContainsRelationship).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})

describe("IAMClient - Party Types", () => {
  let testId: string
  let testPartyType = {
    name: "TESTNAME",
    isBuiltIn: false
  } as PartyType

  it("Should return 200 on Create", (done) => {
    iamClient.PartyTypes().Create(testPartyType)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    iamClient.PartyTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotPartyType = response.response as PartyType
          expect(gotPartyType.name).to.equal(testPartyType.name)
          expect(gotPartyType.isBuiltIn).to.equal(testPartyType.isBuiltIn)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testPartyType = {
      name: "TESTNAME2",
      isBuiltIn: false
    } as PartyType

    iamClient.PartyTypes().Update({
      name: "TESTNAME2",
      isBuiltIn: false
    } as PartyType)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    iamClient.PartyTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotPartyType = response.response as PartyType
          expect(gotPartyType.name).to.equal(testPartyType.name)
          expect(gotPartyType.isBuiltIn).to.equal(testPartyType.isBuiltIn)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  // Party Type List test omitted since party type list is not implemented
})

/*describe("IAMClient - Party Attribute Definitions", () => {
  let testId: string
  let testAttribute = {
    countryId: "TESTCOUNTRYID",
    partyTypeIds: ["TESTPARTYTYPE"],
    isPii: false,
    formControl: {} as Control
  } as PartyAttributeDefinition

  it("Should return 200 on Create", (done) => {
    iamClient.PartyAttributeDefinitions().Create(testAttribute)
      .subscribe((response) => {
        if (response.status === 200) {
          testId = response.response.id
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get", (done) => {
    iamClient.PartyAttributeDefinitions().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotPartyAttributeDefinition = response.response as PartyAttributeDefinition
          expect(gotPartyAttributeDefinition.countryId).to.equal(testAttribute.countryId)
          testAttribute.partyTypeIds.forEach((id) => {
            expect(gotPartyAttributeDefinition.partyTypeIds).to.include(id)
          })
          expect(gotPartyAttributeDefinition.isPii).to.equal(testAttribute.isPii)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Update", (done) => {
    testAttribute = {
      countryId: "TESTCOUNTRYID2",
      partyTypeIds: ["TESTPARTYTYPE2"],
      isPii: false,
      formControl: {} as Control
    } as PartyAttributeDefinition

    iamClient.PartyAttributeDefinitions().Update({
      countryId: "TESTCOUNTRYID2",
      partyTypeIds: ["TESTPARTYTYPE2"],
      isPii: false,
      formControl: {} as Control
    } as PartyAttributeDefinition)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on Get after Update", (done) => {
    iamClient.PartyTypes().Get(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotPartyType = response.response as PartyType
          expect(gotPartyType.name).to.equal(testPartyType.name)
          expect(gotPartyType.isBuiltIn).to.equal(testPartyType.isBuiltIn)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })

  it("Should return 200 on List", (done) => {
    iamClient.Relationships().List({} as RelationshipListOptions)
      .subscribe((response) => {
        if (response.status === 200) {
          let gotRelationshipList = response.response as RelationshipList
          let listContainsRelationship = false
          gotRelationshipList.items.forEach((relationship) => {
            if (relationship.id == testId) {
              expect(relationship.relationshipTypeId).to.equal(testRelationship.relationshipTypeId)
              expect(relationship.firstParty).to.equal(testRelationship.firstParty)
              expect(relationship.secondParty).to.equal(testRelationship.secondParty)
              listContainsRelationship = true
            }
          })
          expect(listContainsRelationship).to.equal(true)
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})*/
