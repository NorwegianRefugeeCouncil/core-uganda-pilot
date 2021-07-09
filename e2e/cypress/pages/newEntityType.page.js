import { Urls } from '../helpers';

const NAME = '[data-cy=name]';
const IS_BULITIN = '[data-cy=isBuiltIn]';;
const SAVE_BUTTON = '[data-cy=save-btn]';

export default class NewEntityTypePage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.NEW_ENTITYTYPE_URL);
        cy.visit(Urls.NEW_ENTITYTYPE_URL);
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
}
