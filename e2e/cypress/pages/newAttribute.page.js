import { Urls } from '../helpers';

const NAME = '[data-cy=name]';
const VALUE_TYPE = '[data-cy=type]';
const SUBJECT_TYPE = '[data-cy=subject]';
// const PERSONAL_INFO = '[data-cy=personal-info-chkbx]';
const LANGUAGE = '[data-cy=language]';
const SAVE_BUTTON = '[data-cy=save-btn]';

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
