import { testId, URL } from '../helpers';
import AttributePage from './attributePage';

const ATTRIBUTE_ROW = '[data-testid=attribute]';
const NEW_ATTR_BTN = testId('new-attribute-btn');

export default class AttributeOverviewPage {
    constructor() {
        this.visitPage();
    }

    visitPage = () => {
        cy.visit(URL.attributes);
        return this;
    };

    clickNewAttributeBtn = () => {
        return cy.get(NEW_ATTR_BTN).click();
    };

    selectLastAttribute = () => {
        return cy.get(ATTRIBUTE_ROW).last();
    };

    // Expect newest attribute to be at the end of the list
    newestAttributePage = () => {
        cy.get(ATTRIBUTE_ROW).last().invoke('attr', 'href').as('href');
        return new AttributePage(cy.get('@href'));
    };
}
