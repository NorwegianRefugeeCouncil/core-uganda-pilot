import { Urls } from '../helpers';

const NAME = '[data-testid=name]';
const VALUE_TYPE = '[data-testid=type]';
const SUBJECT_TYPE = '[data-testid=subject]';
// const PERSONAL_INFO = '[data-testid=personal-info-chkbx]';
const LANGUAGE = '[data-testid=language]';
const SAVE_BUTTON = '[data-testid=save-btn]';

export default class NewAttributePage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.NEW_ATTRIBUTE_URL);
        cy.visit(Urls.NEW_ATTRIBUTE_URL);
        return this;
    };

    typeName = (value) => {
        cy.get(NAME).type(value);
        return this;
    };

    selectValueType = (value) => {
        cy.get(VALUE_TYPE).select(value);
        return this;
    };

    selectSubjectType = (value) => {
        cy.get(SUBJECT_TYPE).select(value);
        return this;
    };

    selectLanguage = (value) => {
        cy.get(LANGUAGE).select(value);
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
