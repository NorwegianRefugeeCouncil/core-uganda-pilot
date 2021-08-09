import { URL, testId } from '../helpers';

const NAME = testId('name');
const VALUE_TYPE = testId('type');
const SUBJECT_TYPE = testId('subject');
const PERSONAL_INFO = testId('personal-info-chkbx');
const LANGUAGE = testId('language');
const LANGUAGE_DSP = testId('language-display');
const TRANSLATION_LONG = testId('translation-long');
const TRANSLATION_SHORT = testId('translation-short');
const SAVE_BUTTON = testId('save-btn');

export default class AttributePage {
    constructor(href) {
        if (href != null) {
            href.then((h) => cy.visit(h));
        } else {
            this.visitPage();
        }
    }

    visitPage = () => {
        cy.log('navigating to %s', URL.NEW_ATTRIBUTE);
        cy.visit(URL.NEW_ATTRIBUTE);
        return this;
    };

    getName = () => cy.get(NAME);
    typeName = (value) => {
        this.getName().type(value);
        return this;
    };

    getValueType = () => cy.get(VALUE_TYPE);
    selectValueType = (value) => {
        this.getValueType().select(value);
        return this;
    };

    getSubjectType = () => cy.get(SUBJECT_TYPE);
    selectSubjectType = (value) => {
        this.getSubjectType().select(value);
        return this;
    };

    getPersonalInfo = () => cy.get(PERSONAL_INFO);
    checkPersonalInfo = () => {
        this.getPersonalInfo().check();
        return this;
    };

    getLanguage = () => cy.get(LANGUAGE);
    getLanguageDsp = () => cy.get(LANGUAGE_DSP);
    selectLanguage = (value) => {
        this.getLanguage().select(value);
        return this;
    };

    getTranslationLong = () => cy.get(TRANSLATION_LONG);
    setTranslationLong = (value) => {
        this.getTranslationLong().type(value);
        return this;
    };

    getTranslationShort = () => cy.get(TRANSLATION_SHORT);
    setTranslationShort = (value) => {
        this.getTranslationShort().type(value);
        return this;
    };

    save = () => {
        cy.get(SAVE_BUTTON).click();
        return this;
    };

    clearName = () => {
        this.getName().clear();
        return this;
    };
}
