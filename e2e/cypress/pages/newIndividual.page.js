import { URL } from '../helpers';

const TEXT_ATTRIBUTES = '[data-testid=text-attribute]';
const RELATIONSHIP_TYPE = '[data-testid=relationshipType]';
const RELATED_PARTY = '[data-testid=relatedParty]';
const ADD_BTN = '[data-testid=add-btn]';
const SAVE_BUTTON = '[data-testid=save-btn]';

export default class NewIndividualPage {
    visitPage = () => {
        cy.log('navigating to %s', URL.NEW_INDIVIDUAL);
        cy.visit(URL.NEW_INDIVIDUAL);
        return this;
    };

    typeTextAttributes = (value) => {
        cy.get(TEXT_ATTRIBUTES).each(($el) => {
            cy.wrap($el).type(value);
        });
        return this;
    };

    selectRelationshipType = (value) => {
        cy.get(RELATIONSHIP_TYPE).select(value);
        return this;
    };

    typeRelatedParty = (value) => {
        cy.get(RELATED_PARTY).type(value);
        return this;
    };

    addParty = () => {
        cy.get(ADD_BTN).click();
        return this;
    };

    save = () => {
        cy.get(SAVE_BUTTON).click();
        return this;
    };

    clearTextAttributes = () => {
        cy.get(TEXT_ATTRIBUTES).each(($el) => {
            cy.wrap($el).clear();
        });
        return this;
    };
}
