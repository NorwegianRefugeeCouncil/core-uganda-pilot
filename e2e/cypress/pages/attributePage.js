import { URL, testId } from '../helpers';

const NAME = testId('name');
const VALUE_TYPE = testId('type');
const SUBJECT_TYPE = testId('subject');
const PERSONAL_INFO = testId('personal-info-chkbx');
const LANGUAGE_SELECT = testId('language-select');
const LANGUAGE_DSP = '-display';
const TRANSLATION_LONG = '-translation-long';
const TRANSLATION_SHORT = '-translation-short';
const SAVE_BUTTON = testId('save-btn');
const REMOVE_BUTTON = '-remove-btn';

export default class AttributePage {
    constructor(href) {
        if (href != null) {
            href.then((h) => cy.visit(h));
        } else {
            this.visitNewAttributePage();
        }
    }

    visitNewAttributePage = () => {
        cy.visit(URL.NEW_ATTRIBUTE);
        return this;
    };

    getName = () => cy.get(NAME);
    setName = (value) => {
        this.getName().clear().type(value);
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

    getLanguageDsp = (lang) => cy.get(testId(lang + LANGUAGE_DSP));

    getLanguage = () => cy.get(LANGUAGE_SELECT);
    selectLanguage = (value) => {
        this.getLanguage().select(value);
        return this;
    };

    getTranslationLong = (lang) => cy.get(testId(lang + TRANSLATION_LONG));
    setTranslationLong = (lang, value) => {
        this.getTranslationLong(lang).clear().type(value);
        return this;
    };

    getTranslationShort = (lang) => cy.get(testId(lang + TRANSLATION_SHORT));
    setTranslationShort = (lang, value) => {
        this.getTranslationShort(lang).clear().type(value);
        return this;
    };

    removeTranslation = (lang) => cy.get(testId(lang + REMOVE_BUTTON)).click();

    save = () => {
        cy.get(SAVE_BUTTON).click();
        return this;
    };

    clearName = () => {
        this.getName().clear();
        return this;
    };
}
