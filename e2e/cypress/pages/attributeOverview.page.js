import { URL } from '../helpers';
import AttributePage from './attributePage';

const ATTRIBUTE_ROWS = '[data-testid=attribute]';

export default class AttributeOverviewPage {
    constructor() {
        this.visitPage();
    }

    visitPage = () => {
        cy.log('navigating to %s', URL.ATTRIBUTE);
        cy.visit(URL.ATTRIBUTE);
        return this;
    };

    selectAttribute = () => {
        return cy.get(ATTRIBUTE_ROWS).last();
    };

    visitAttribute = () => {
        cy.get(ATTRIBUTE_ROWS).last().click();
    };
    attributePageForLatest = () => {
        cy.get(ATTRIBUTE_ROWS).last().invoke('attr', 'href').as('href');
        return new AttributePage(cy.get('@href'));
    };
}
