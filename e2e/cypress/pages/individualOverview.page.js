import { testId, URL } from '../helpers';

const selector = {
    individualRows: testId('individual'),
    individual: testId('individual'),
    newIndividualBtn: testId('new-individual-btn'),
    search: testId('search'),
    searchBtn: testId('search-btn'),
};

export default class IndividualOverviewPage {
    visitPage = () => {
        cy.visit(URL.INDIVIDUALS);
        return this;
    };

    newIndividual = () => {
        return cy.get(selector.newIndividualBtn).click();
    };

    searchFor = value => {
        return cy
            .get(selector.search)
            .type(value)
            .get(selector.searchBtn)
            .wait(2000) // server is slow
            .click()
            .get(selector.individual)
            .last()
            .invoke('attr', 'href');
    };
}
