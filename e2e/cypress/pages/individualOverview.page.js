import { Urls } from '../helpers';
import NewIndividualPage from './newIndividual.page';

const INDIVIDUAL_ROWS = '[data-testid=individual]';

export default class IndividualOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.INDIVIDUALS_URL);
        cy.visit(Urls.INDIVIDUALS_URL);
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
