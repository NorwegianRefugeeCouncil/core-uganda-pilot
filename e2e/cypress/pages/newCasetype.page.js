import { Urls } from '../helpers';

const NAME = '[data-testid=name]';
const PARTY_TYPE = '[data-testid=partytype]';
const TEAM = '[data-testid=team]';
const TEMPLATE = '[data-testid=template]';
const SAVE_BUTTON = '[data-testid=save-btn]';

export default class NewCaseTypePage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.NEW_CASETYPE_URL);
        cy.visit(Urls.NEW_CASETYPE_URL);
        return this;
    };

    typeName = (value) => {
        cy.get(NAME).type(value);
        return this;
    };

    selectPartyType = (value) => {
        cy.get(PARTY_TYPE).select(value);
        return this;
    };

    selectTeam = (value) => {
        cy.get(TEAM).select(value);
        return this;
    };

    typeTemplate = (value) => {
        cy.get(TEMPLATE).invoke('val', value);
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
