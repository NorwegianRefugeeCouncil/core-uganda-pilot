import { URL } from '../helpers';
import ids from '../fixtures/ids.json';
import CasesOverviewPage from '../pages/casesOverview.page';
import NewCasePage from '../pages/newCase.page';
import CaseTypePage from '../pages/caseTypePage';
import testTemplate from '../fixtures/test_casetemplate.json';
import CasetypesOverviewPage from '../pages/casetypesOverview.page';
import '../support/commands';

const DATA = {
    name: 'test casetype',
    partyTypeId: ids.IndividualPartyTypeID,
    partyId: ids.BoDiddleyID,
    teamId: ids.DTeamID,
    formText: testTemplate,
    form: {
        text: 'test',
        email: 'test@whatever.net',
        phone: '0897-938900',
        url: 'https://www.hello-world.com',
        date: '1933-02-24',
        textarea: 'test',
        dropdown: '0',
        checkbox: '0',
        radio: '0',
    },
    formUpd: {
        text: 'test-updated',
        email: 'test_updated@whatever.net',
        phone: '0897-938901',
        url: 'https://www.hello-world-updated.com',
        date: '1933-02-25',
        textarea: 'test-updated',
        dropdown: '1',
        checkbox: '1',
        radio: '1',
    },
};

describe('Case Page', function () {
    beforeEach('Login', () => {
        cy.login('courtney.lare@email.com');
    });
    before('Seed DB with test casetype', () => {
        cy.login('courtney.lare@email.com');
        const caseTypePage = new CaseTypePage();
        caseTypePage
            .setName(DATA.name)
            .setPartyType(DATA.partyTypeId)
            .setTeam(DATA.teamId)
            .setTemplate(JSON.stringify(DATA.formText, null, 2))
            .save();

        // store the new CaseTypeId in the DATA object
        const casetypesOverviewPage = new CasetypesOverviewPage();
        casetypesOverviewPage
            .visitPage()
            .selectNewestCasetype()
            .invoke('attr', 'href')
            .then(h => {
                // h should look like "/settings/casetypes/<id>"
                const sepIdx = h.lastIndexOf('/');
                DATA.caseTypeId = h.slice(sepIdx + 1);
            });
    });
    describe('Navigate', () => {
        it('should navigate to new Case page from the case overview page', () => {
            const casesOverviewPage = new CasesOverviewPage();
            casesOverviewPage.clickNewCaseBtn().url().should('include', URL.newCase);
        });
    });
    describe('Create', () => {
        it('should create a new Case', () => {
            const newCasePage = new NewCasePage();
            newCasePage.setCaseType(DATA.caseTypeId).setParty(DATA.partyId).fillOutForm(DATA.form).submitForm();
        });
    });

    describe('Verify creation', () => {
        it('should verify that the Case was properly created', () => {
            const casesOverviewPage = new CasesOverviewPage();
            casesOverviewPage.selectLastCase().should('contain.text', DATA.name);
            const casePage = casesOverviewPage.newestCasePage();

            // Verify values
            casePage.getDonePill().should('contain.text', 'Open');
            casePage.verifyForm(DATA.form);
            casePage.getDoneCheck().should('not.be.checked');
        });
    });

    describe('Update', () => {
        it('should update the case', () => {
            const casesOverviewPage = new CasesOverviewPage();
            const casePage = casesOverviewPage.newestCasePage();

            // Verify values
            new NewCasePage(true).fillOutForm.apply(casePage, [DATA.formUpd]);
            casePage.setDoneCheck().save();
        });
    });

    describe('Verify update', () => {
        it('should verify that the Case was properly updated', () => {
            const casesOverviewPage = new CasesOverviewPage();
            casesOverviewPage.selectLastCase().should('contain.text', DATA.name);
            const casePage = casesOverviewPage.newestCasePage();

            // Verify values
            casePage.getDonePill().should('contain.text', 'Closed');
            casePage.verifyForm(DATA.formUpd);
            casePage.getDoneCheck().should('not.exist');
        });
    });
});
