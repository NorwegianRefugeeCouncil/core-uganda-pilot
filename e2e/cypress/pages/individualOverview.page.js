import { testId, URL } from '../helpers';
import IndividualPage from './individualPage';

const ID = {
    INDIVIDUAL_ROWS: testId('individual'),
    INDIVIDUAL: testId('individual'),
    SEARCH: testId('search'),
    SEARCH_BTN: testId('search-btn'),
};

export default class IndividualOverviewPage {
    visitPage = () => {
        cy.visit(URL.INDIVIDUALS);
        return this;
    };

    selectIndividual = () => {
        return cy.get(ID.INDIVIDUAL_ROWS).last();
    };

    visitIndividual = () => {
        cy.get(ID.INDIVIDUAL_ROWS).last().click();
        return new IndividualPage();
    };

    searchFor = value => {
        return cy
            .get(ID.SEARCH)
            .type(value)
            .get(ID.SEARCH_BTN)
            .click()
            .get(ID.INDIVIDUAL)
            .last()
            .invoke('attr', 'href');
    };
}
