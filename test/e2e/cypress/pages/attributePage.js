import {testId, URL} from '../helpers';

const NAME = testId('name');
const LABEL = testId('label');
const CONTROL_TYPE = testId('control-type');
const SUBJECT_TYPE = testId('subject');
const PERSONAL_INFO = testId('personal-info-chkbx');
const IS_REQUIRED = testId('required-checkbox');
const SAVE_BUTTON = testId('save-btn');

export default class AttributePage {
    constructor(href) {
        if (href != null) {
            href.then(h => cy.visit(h));
        } else {
            this.visitNewAttributePage();
        }
    }

    visitNewAttributePage = () => {
        cy.visit(URL.newAttribute);
        return this;
    };

    getName = () => cy.get(NAME);
    setName = value => {
        this.getName().clear().type(value);
        return this;
    };

    getLabel = () => cy.get(LABEL);
    setLabel = value => {
        this.getLabel().clear().type(value);
        return this;
    };

    getControlType = () => cy.get(CONTROL_TYPE);
    selectControlType = value => {
        this.getControlType().select(value);
        return this;
    };

    getSubjectType = () => cy.get(SUBJECT_TYPE);
    selectSubjectType = value => {
        this.getSubjectType().select(value);
        return this;
    };

    getRequired = () => cy.get(IS_REQUIRED);
    checkRequired = () => {
        this.getRequired().check();
        return this;
    };
    uncheckRequired = () => {
        this.getRequired().uncheck();
        return this;
    };

    getPersonalInfo = () => cy.get(PERSONAL_INFO);
    checkPersonalInfo = () => {
        this.getPersonalInfo().check();
        return this;
    };

    save = () => {
        cy.get(SAVE_BUTTON).click().wait(200);
        return this;
    };
}
