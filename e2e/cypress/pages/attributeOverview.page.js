import { testId, URL } from '../helpers';
import AttributePage from './attributePage';

const ATTRIBUTE_ROWS = '[data-testid=attribute]';
const NEW_ATTR_BTN = testId('new-attribute-btn');

export default class AttributeOverviewPage {
    constructor() {
        this.visitPage();
    }

    visitPage = () => {
        cy.visit(URL.ATTRIBUTES);
        return this;
    };

    clickNewAttributeBtn = () => {
        return cy.get(NEW_ATTR_BTN).click();
    };

    selectLastAttribute = () => {
        return cy.get(ATTRIBUTE_ROWS).last();
    };

    visitAttribute = () => {
        this.selectLastAttribute().click();
    };

    // Expect newest attribute to be at the end of the list
    attributePageForNewest = () => {
        cy.get(ATTRIBUTE_ROWS).last().invoke('attr', 'href').as('href');
        return new AttributePage(cy.get('@href'));
    };
}
