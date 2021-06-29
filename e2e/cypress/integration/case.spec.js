import CasesOverviewPage from '../pages/casesOverview.page';
import NewCasePage from '../pages/newCase.page';

const mockUpdatedText = 'update';

describe('Case Page', function () {
    describe('Create', () => {
        it('should navigate to New Case Form Page when NewCaseBtn is selected on Case Overview Page', () => {
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .visitPage()
                .openNewCase()
                .newCaseForm()
                .should('exist');
        });

        it('should create new case with Type, Party and Description', () => {
            const descriptionValue = 'Some Description';
            const newCasePage = new NewCasePage();
            newCasePage
                .visitPage()
                .selectCaseType()
                .selectParty()
                .typeDescription(descriptionValue)
                .submitForm()
                .getNewCase()
                .should('contain.text', descriptionValue);
        });
    });

    describe('Update', () => {
        it('should update description on existing case', () => {
            const newDescription = 'My new Description';
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .visitPage()
                .selectCase()
                .clearDescriptionValue()
                .typeDescription(newDescription)
                .submitUpdate()
                .getDescriptionValue()
                .should('contain.text', newDescription);
        });

        it('should add Referral to existing case', () => {
            const descriptionText = 'Referral Description';
            const caseOverviewPage = new CasesOverviewPage();
            caseOverviewPage
                .visitPage()
                .selectCase()
                .selectReferral()
                .typeReferralDescription(descriptionText)
                .submitReferral()
                .getOpenReferralItem()
                .should('exist');
        });

        xit('saved the updated case', () => {
            cy.visit('/cases');
            cy.get('tr')
                .last()
                .within(($row) => {
                    cy.wrap($row).should('contain.text');
                });
            cy.get('tr').last().click({ force: true });
            cy.get('textarea[data-cy=description]').should(
                'contain.text',
                mockUpdatedText
            );
            cy.get('input[data-cy=done-check]')
                .invoke('prop', 'checked')
                .then((checked) => expect(checked).to.be.true);
        });
    });
});
