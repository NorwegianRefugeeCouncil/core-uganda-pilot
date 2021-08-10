import { URL } from '../helpers';
import CasesOverviewPage from '../pages/casesOverview.page';

const mockText = 'Test';
const mockUpdatedText = 'Updated Test';

describe('Case Page', function () {
    describe.skip('Prepares the db', () => {
        it('should add the test casetype', () => {
            cy.fixture('test_casetype')
                .then((tmpl) => {
                    cy.request('POST', URL.CASETYPES, tmpl);
                })
                .then((res) => {
                    expect(res.status).to.eq(303); // expect to be redirected after successful POST
                });
        });
    });
    describe('Create', () => {
        it('should navigate to New Case Page from Case Overview Page and submit a new case', () => {
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .openNewCase()
                .newCaseForm()
                .selectCaseType()
                .selectParty()
                .typeForm(mockText)
                .submitForm()
                .getAlertMessage()
                .should('contain.text', 'Case successfully created');
        });
    });

    describe('Update', () => {
        it('should add Referral to existing case', () => {
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .visitPage()
                .selectCase()
                .selectReferral()
                .submitReferral()
                .getOpenReferralItem()
                .should('exist');
        });

        it('should update form text to existing case and save', () => {
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .visitPage()
                .selectCase()
                .clearForm()
                .typeForm(mockUpdatedText)
                .submitUpdate()
                .getAlertMessage()
                .should('contain.text', 'Case successfully created');
        });
    });
});
