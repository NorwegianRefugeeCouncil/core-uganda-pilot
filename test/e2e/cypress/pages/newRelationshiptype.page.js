import { URL } from '../helpers';

const NAME = '[data-testid=name]';
const IS_DIRECTIONAL = '[data-testid=isDirectional]';
const FIRST_PARTY_ROLE = '[data-testid=firstPartyRole]';
const SECOND_PARTY_ROLE = '[data-testid=secondPartyRole]';
const FIRST_PARTY_TYPE = '[data-testid=firstPartyType]';
const SECOND_PARTY_TYPE = '[data-testid=secondPartyType]';
const SAVE_BUTTON = '[data-testid=save-btn]';

export default class NewRelationshiptypePage {
    visitPage = () => {
        cy.visit(URL.newRelationshipType);
        return this;
    };

    typeName = value => {
        cy.get(NAME).type(value);
        return this;
    };

    checkIsDirectional = () => {
        cy.get(IS_DIRECTIONAL).check();
        return this;
    };

    typeFirstPartyRole = value => {
        cy.get(FIRST_PARTY_ROLE).type(value);
        return this;
    };

    typeSecondPartyRole = value => {
        cy.get(SECOND_PARTY_ROLE).type(value);
        return this;
    };

    selectFristPartyType = value => {
        cy.get(FIRST_PARTY_TYPE).select(value);
        return this;
    };

    selectSecondPartyType = value => {
        cy.get(SECOND_PARTY_TYPE).select(value);
        return this;
    };

    save = () => {
        cy.get(SAVE_BUTTON).click();
        return this;
    };

    clearName = () => {
        cy.get(NAME).clear();
        return this;
    };
}
