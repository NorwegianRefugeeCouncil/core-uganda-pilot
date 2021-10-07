import { IAMClient, CMSClient } from './coreApiClient'
import {
  Case,
  CaseListOptions,
  Comment,
  CaseType,
  CaseTypeListOptions,
  CommentListOptions,
  Form,
  Party,
  PartySearchOptions,
  PartyListOptions,
  PartyList,
  Country,
  CountryListOptions, CountryList, CommentList
} from './types/models';
import { expect } from 'chai'

const noop = () => {}

describe("CMSClient - Comments", () => {
  let testId: string
  let testComment: Comment

  before(() => {
    testComment = {
      caseId: "TESTCASEID"
    } as Comment
  })

  it("Should return 200 on Create and params match", (done) => {
    CMSClient.Comments().Create(testComment)
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
    CMSClient.Comments().Get(testId)
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

  before(() => {
    testComment = {
      id: testId,
      caseId: "TESTCASEID2"
    } as Comment
  })

  it("Should return 200 on Update", (done) => {
    CMSClient.Comments().Update({
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
    CMSClient.Comments().Get(testId)
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
    CMSClient.Comments().List({} as CommentListOptions)
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
    CMSClient.Comments().Delete(testId)
      .subscribe((response) => {
        if (response.status === 200) {
          done()
        } else {
          done(expect(response.status).to.equal(200));
        }
      })
  })
})
/*describe("CMSClient - Case Types", () => {
  it("Should return 200 on Create and params match", () => {
    CMSClient.CaseTypes().Create({
      id: "TESTID",
      name: "TESTNAME",
      partyTypeId: "TESTPARTYTYPEID",
      teamId: "TESTTEAMID",
      form: {} as Form,
      intakeCaseType: false
    } as CaseType)
      .subscribe((response) => {
        expect(response.status).to.equal(200)
      })
  })

  it("Should return 200 on Get", () => {
    CMSClient.CaseTypes().Get("TESTID")
      .subscribe((response) => {
        expect(response.status).to.equal(200)
      })
  })

  it("Should return 200 on Update", () => {
    CMSClient.CaseTypes().Update({
      id: "TESTID",
      name: "TESTNAME2",
      partyTypeId: "TESTPARTYTYPEID",
      teamId: "TESTTEAMID",
      form: {} as Form,
      intakeCaseType: false
    } as CaseType)
      .subscribe((response) => {
        expect(response.status).to.equal(200)
      })
  })

  it("Should return 200 on Get after Update", () => {
    CMSClient.CaseTypes().Get("TESTID")
      .subscribe((response) => {
        expect(response.status).to.equal(200)
      })
  })

  it("Should return 200 on List", () => {
    CMSClient.CaseTypes().List({} as CaseTypeListOptions)
      .subscribe((response) => {
        expect(response.status).to.equal(200)
      })
  })
})*/