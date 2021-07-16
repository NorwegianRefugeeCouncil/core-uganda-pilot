import { Urls } from '../helpers';
import NewAttributePage from './newAttribute.page';

const ATTRIBUTE_ROWS = '[data-cy=attribute]';

export default class AttributeOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.ATTRIBUTE_URL);
        cy.visit(Urls.ATTRIBUTE_URL);
        return this;
    };

    selectAttribute = () => {
        return cy.get(ATTRIBUTE_ROWS).last();
    };

    visitAttribute = () => {
        cy.get(ATTRIBUTE_ROWS).last().click();
        return new NewAttributePage();
    };
}
