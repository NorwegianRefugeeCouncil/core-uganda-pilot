import { testId, URL } from '../helpers';
import ids from '../fixtures/ids.json';
import CasesOverviewPage from '../pages/casesOverview.page';
import NewCasePage from '../pages/newCase.page';
import CasePage from '../pages/case.page';
import CaseTypePage from '../pages/caseTypePage';
import testTemplate from '../fixtures/test_casetemplate.json';
import CasetypesOverviewPage from '../pages/casetypesOverview.page';

const DATA = {
    NAME: 'test casetype',
    PARTYTYPEID: ids.IndividualPartyTypeID,
    PARTYID: ids.BoDiddleyID,
    TEAMID: ids.DTeamID,
    TEMPLATE: testTemplate,
    FORM: {
        text: 'test',
        email: 'test@whatever.net',
        phone: '0897-938900',
        url: 'www.hello-world.com',
        date: '1933-02-24',
        textarea: 'test',
        dropdown: '0',
        checkbox: '0',
        radio: '0',
    },
    FORM_U: {
        text: 'test-updated',
        email: 'test_updated@whatever.net',
        phone: '0897-938901',
        url: 'www.hello-world-updated.com',
        date: '1933-02-25',
        textarea: 'test-updated',
        dropdown: '1',
        checkbox: '1',
        radio: '1',
    },
};

describe('Case Page', function () {
    before('Seed DB with test casetype', () => {
        cy.login('courtney.lare@email.com');
        const caseTypePage = new CaseTypePage();
        caseTypePage
            .setName(DATA.NAME)
            .setPartyType(DATA.PARTYTYPEID)
            .setTeam(DATA.TEAMID)
            .setTemplate(JSON.stringify(DATA.TEMPLATE, null, 2))
            .save();

        // Get the new CaseTypeId
        const casetypesOverviewPage = new CasetypesOverviewPage();
        casetypesOverviewPage
            .visitPage()
            .selectNewestCasetype()
            .invoke('attr', 'href')
            .then(h => {
                // h should look like "/settings/casetypes/<id>"
                const sepIdx = h.lastIndexOf('/');
                DATA.CASETYPEID = h.slice(sepIdx + 1);
            });
    });
    describe('Navigate', () => {
        it('should navigate to new Case page from the case overview page', () => {
            const casesOverviewPage = new CasesOverviewPage();
            casesOverviewPage.clickNewCaseBtn().url().should('include', URL.NEW_CASE);
        });
    });
    describe('Create', () => {
        it('should create a new Case', () => {
            const newCasePage = new NewCasePage();
            newCasePage.setCaseType(DATA.CASETYPEID).setParty(DATA.PARTYID).fillOutForm(DATA.FORM).submitForm();
            // TODO get and store caseID
            // store caseId
            // .get(testId('case-id'))
            // .then($c => (DATA.CASEID = $c.text()));
        });
    });

    // FIXME
    describe.skip('Verify creation', () => {
        it('should verify that the Case was properly created', () => {
            const casePage = new CasePage(URL.CASES + '/' + DATA.CASEID);

            // Verify values
            casePage.getRecipient().should('contain.text', ids.TestIndividualName);
            casePage.getDonePill().should('contain.text', 'Open');
            casePage.verifyForm(DATA.FORM);
            casePage.getDoneCheck().should('not.be.checked');
        });
    });

    // FIXME
    describe.skip('Update', () => {
        it('should update the case', () => {
            const casePage = new CasePage(URL.CASES + '/' + DATA.CASEID);

            // Verify values
            new NewCasePage(true).fillOutForm.apply(casePage, [DATA.FORM_U]);
            casePage.setDoneCheck().save();
        });
    });

    // FIXME
    describe.skip('Verify update', () => {
        it('should verify that the Case was properly updated', () => {
            const casePage = new CasePage(URL.CASES + '/' + DATA.CASEID);

            // Verify values
            casePage.getRecipient().should('contain.text', ids.TestIndividualName);
            casePage.getDonePill().should('contain.text', 'Closed');
            casePage.verifyForm(DATA.FORM_U);
            casePage.getDoneCheck().should('not.exist');
        });
    });
});
