import IndividualPage from '../pages/individualPage';
import IndividualOverviewPage from '../pages/individualOverview.page';
import '../support/commands';

describe('Country-specific Individual Attributes (Colombia)', function () {
    beforeEach(() => {
        cy.login('claudia.garcia@email.com');
    });
    describe('Navigate', () => {
        it('should navigate to new Individual page from Individuals overview', () => {
            const individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage.visitPage().newIndividual();
        });
    });
    describe('Check attributes for country', () => {
        it('CO (Colombia)', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage.visitPage();
            newIndividualPage.verifyCountrySpecificAttributes('CO');
        });
    });
});

describe('Country-specific Individual Attributes (Uganda)', function () {
    beforeEach(() => {
        cy.login('stephen.kabagambe@email.com');
    });
    describe('Navigate', () => {
        it('should navigate to new Individual page from Individuals overview', () => {
            const individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage.visitPage().newIndividual();
        });
    });
    describe('Check attributes for country', () => {
        it('UG (Uganda)', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage.visitPage();
            newIndividualPage.verifyCountrySpecificAttributes('UG');
        });
    });
});
