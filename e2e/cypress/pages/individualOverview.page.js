import { testId, URL } from '../helpers';
import IndividualPage from './individualPage';

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
    selectIndividual = () => {
        return cy.get(selector.individualRows).last();
    };

    visitIndividual = () => {
        cy.get(selector.individualRows).last().click();
        return new IndividualPage();
    };

    searchFor = value => {
        return cy
            .get(selector.search)
            .type(value)
            .get(selector.searchBtn)
            .click()
            .get(selector.individual)
            .last()
            .invoke('attr', 'href');
    };
}
