import { Urls } from '../helpers';

const NAME = '[data-cy=name]';
const IS_DIRECTIONAL = '[data-cy=isDirectional]';
const FIRST_PARTY_ROLE = '[data-cy=firstPartyRole]';
const SECOND_PARTY_ROLE = '[data-cy=secondPartyRole]';
const FIRST_PARTY_TYPE = '[data-cy=firstPartyType]';
const SECOND_PARTY_TYPE = '[data-cy=secondPartyType]';
const SAVE_BUTTON = '[data-cy=save-btn]';

export default class NewRelationshiptypePage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.NEW_RELATIONSHIPTYPE_URL);
        cy.visit(Urls.NEW_RELATIONSHIPTYPE_URL);
        return this;
    };

    typeName = (value) => {
        cy.get(NAME).type(value);
        return this;
    };

    checkIsDirectional = () => {
        cy.get(IS_DIRECTIONAL).check();
        return this;
    };

    typeFirstPartyRole = (value) => {
        cy.get(FIRST_PARTY_ROLE).type(value);
        return this;
    };

    typeSecondPartyRole = (value) => {
        cy.get(SECOND_PARTY_ROLE).type(value);
        return this;
    };

    selectFristPartyType = (value) => {
        cy.get(FIRST_PARTY_TYPE).select(value);
        return this;
    };

    selectSecondPartyType = (value) => {
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
