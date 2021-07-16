import { Urls } from '../helpers';

const TEXT_ATTRIBUTES = '[data-cy=text-attribute]';
const RELATIONSHIP_TYPE = '[data-cy=relationshipType]';
const RELATED_PARTY = '[data-cy=relatedParty]';
const ADD_BTN = '[data-cy=add-btn]';
const SAVE_BUTTON = '[data-cy=save-btn]';

export default class NewIndividualPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.NEW_INDIVIDUAL_URL);
        cy.visit(Urls.NEW_INDIVIDUAL_URL);
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
