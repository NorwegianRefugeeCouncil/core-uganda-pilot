import { URL } from '../helpers';
import NewIndividualPage from './newIndividual.page';

const INDIVIDUAL_ROWS = '[data-testid=individual]';

export default class IndividualOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', URL.INDIVIDUALS);
        cy.visit(URL.INDIVIDUALS);
        return this;
    };

    selectIndividual = () => {
        return cy.get(INDIVIDUAL_ROWS).last();
    };

    visitIndividual = () => {
        cy.get(INDIVIDUAL_ROWS).last().click();
        return new NewIndividualPage();
    };
}
