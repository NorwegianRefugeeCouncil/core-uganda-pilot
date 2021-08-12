import { URL, testId } from '../helpers';

const ID = {
    TEXT_ATTRIBUTES: testId('text-attribute'),
    RELATIONSHIP: testId('relationship'),
    RELATIONSHIP_TYPE: testId('relationshipType'),
    RELATED_PARTY: testId('relatedParty'),
    SEARCH_RESULT: testId('party-search-result'),
    ADD_BTN: testId('add-btn'),
    SAVE_BUTTON: testId('save-btn'),
};

export default class IndividualPage {
    constructor(href) {
        if (href) {
            href.then(h => cy.visit(h));
        }
    }

    visitPage = () => {
        cy.visit(URL.NEW_INDIVIDUAL);
        return this;
    };

    typeTextAttributes = value => {
        cy.get(ID.TEXT_ATTRIBUTES).each($el => {
            cy.wrap($el).clear().type(value);
        });
        return this;
    };

    verifyTextAttributes = value => {
        cy.get(ID.TEXT_ATTRIBUTES).each($el => {
            cy.wrap($el).should('have.value', value);
        });
    };

    getRelationship = () => cy.get(ID.RELATIONSHIP).first().find('a').last();

    selectRelationshipType = value => {
        cy.get(ID.RELATIONSHIP_TYPE).select(value);
        return this;
    };

    typeRelatedParty = value => {
        cy.get(ID.RELATED_PARTY).type(value).wait(1000);
        return this;
    };

    addParty = () => {
        cy.get(ID.SEARCH_RESULT).click().get(ID.ADD_BTN).click();
        return this;
    };

    save = () => {
        cy.get(ID.SAVE_BUTTON).click();
        return this;
    };

    clearTextAttributes = () => {
        cy.get(ID.TEXT_ATTRIBUTES).each($el => {
            cy.wrap($el).clear();
        });
        return this;
    };
}
