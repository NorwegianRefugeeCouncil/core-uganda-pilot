import { Urls } from '../helpers';

const NAME = '[data-testid=name]';
const IS_BULITIN = '[data-testid=isBuiltIn]';
const SAVE_BUTTON = '[data-testid=save-btn]';

export default class NewPartyTypePage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.NEW_PARTYTYPE_URL);
        cy.visit(Urls.NEW_PARTYTYPE_URL);
        return this;
    };

    typeName = (value) => {
        cy.get(NAME).type(value);
        return this;
    };

    checkIsBuiltIn = () => {
        cy.get(IS_BULITIN).check();
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
