import IndividualPage from '../pages/individualPage';
import '../support/commands';

describe('Country-specific Individual Attributes (Colombia)', function () {
    beforeEach(() => {
        cy.login('claudia.garcia@email.com');
    });
    describe('Check attributes for country', () => {
        it('CO (Colombia)', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage.verifyCountrySpecificAttributes('CO');
        });
    });
});

describe('Country-specific Individual Attributes (Uganda)', function () {
    beforeEach(() => {
        cy.login('stephen.kabagambe@email.com');
    });
    describe('Check attributes for country', () => {
        it('UG (Uganda)', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage.verifyCountrySpecificAttributes('UG');
        });
    });
});
