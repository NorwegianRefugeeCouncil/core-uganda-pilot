import CaseTypePage from '../pages/caseTypePage';
import { URL } from '../helpers';
import ids from '../fixtures/ids.json';
import testTemplate from '../fixtures/test_casetemplate.json';
import CasetypesOverviewPage from '../pages/casetypesOverview.page';

const DATA = {
    name: 'Test casetype',
    nameUpd: 'Test casetype - updated',
    partyType: ids.IndividualPartyTypeID,
    team: ids.DTeamID,
    form: testTemplate,
};

describe('CaseType Page', () => {
    describe('Navigate', () => {
        it('should navigate to new CaseType page from CaseTypes overview page', () => {
            const casetypesOverviewPage = new CasetypesOverviewPage();
            casetypesOverviewPage.clickNewCaseTypeBtn().url().should('include', URL.newCasetype);
        });
    });
    describe('Create', () => {
        it('should create a new CaseType', () => {
            const caseTypePage = new CaseTypePage();
            caseTypePage
                .setName(DATA.name)
                .setPartyType(DATA.partyType)
                .setTeam(DATA.team)
                .setTemplate(JSON.stringify(DATA.form, null, 2))
                .save();

            const casetypesOverviewPage = new CasetypesOverviewPage();
            casetypesOverviewPage.visitPage().selectNewestCasetype().should('contain.text', DATA.name);
        });
    });

    describe('Verify creation', () => {
        it('should verify that the casetype was properly created', () => {
            const casetypesOverviewPage = new CasetypesOverviewPage();
            const caseTypePage = casetypesOverviewPage.caseTypePageForNewest();

            // Verify values
            caseTypePage.getName().should('have.value', DATA.name);
            caseTypePage.getPartyTypeSelect().should('have.value', DATA.partyType);
            caseTypePage.getPartyTypeSelect().should('be.disabled');
            caseTypePage.getTeamSelect().should('have.value', DATA.team);
            caseTypePage.getTeamSelect().should('be.disabled');
            /* FIXME verifying that the saved template matches the provided one via UI is tricky because the text content of the textarea element will not match the initial input (because of formatting and null values) */
            caseTypePage
                .getTemplate()
                .invoke('val')
                .should(v => expect(v).not.to.be.empty);
            caseTypePage.getTemplate().should('have.attr', 'readonly');
        });
    });

    describe('Update', () => {
        it('should update name on existing CaseType', () => {
            const casetypesOverviewPage = new CasetypesOverviewPage();
            const caseTypePage = casetypesOverviewPage.caseTypePageForNewest();
            caseTypePage.setName(DATA.nameUpd).save();
        });
    });
    describe('Verify update', () => {
        it('should verify that the casetype was properly updated', () => {
            const casetypesOverviewPage = new CasetypesOverviewPage();
            const caseTypePage = casetypesOverviewPage.caseTypePageForNewest();

            // Verify values
            caseTypePage.getName().should('have.value', DATA.nameUpd);
            caseTypePage.getPartyTypeSelect().should('have.value', DATA.partyType);
            caseTypePage.getPartyTypeSelect().should('be.disabled');
            caseTypePage.getTeamSelect().should('have.value', DATA.team);
            caseTypePage.getTeamSelect().should('be.disabled');
            caseTypePage
                .getTemplate()
                .invoke('val')
                .should(v => expect(v).not.to.be.empty);
            caseTypePage.getTemplate().should('have.attr', 'readonly');
        });
    });
});
