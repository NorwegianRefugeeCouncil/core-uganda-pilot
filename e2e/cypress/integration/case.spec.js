import CasesOverviewPage from '../pages/casesOverview.page';

const mockText = 'Test';
const mockUpdatedText = 'Updated Test';

describe('Case Page', function () {
    describe('Create', () => {
        it('should navigate to New Case Form Page when NewCaseBtn is selected on Case Overview Page', () => {
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .visitPage()
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
