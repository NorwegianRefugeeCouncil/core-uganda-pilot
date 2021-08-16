import { URL, testId } from '../helpers';
import testTemplate from '../fixtures/test_casetemplate.json';

const ID = {
    CASETYPE_SELECT: testId('casetype-select'),
    CASETYPE_OPTION: testId('caseTypeOption'),
    PARTY_SELECT: testId('party-select'),
    PARTY_OPTION: testId('partyOption'),
    FORM: testId('form'),
    SUBMIT_BTN: testId('submitBtn'),
};

export default class NewCasePage {
    constructor(noredirect) {
        if (!noredirect) {
            this.visitPage();
        }
    }

    visitPage = () => {
        cy.visit(URL.NEW_CASE);
        return this;
    };

    getCaseTypeSelect = () => cy.get(ID.CASETYPE_SELECT);
    setCaseType = value => {
        this.getCaseTypeSelect().select(value);
        return this;
    };

    getPartySelect = () => cy.get(ID.PARTY_SELECT);
    setParty = value => {
        this.getPartySelect().select(value);
        return this;
    };

    fillOutForm = data => {
        for (const { type } of testTemplate.formElements) {
            cy.get(testId('test-' + type)).then($el => {
                const tag = $el[0].tagName;
                if (tag === 'INPUT') {
                    switch ($el[0].getAttribute('type')) {
                        case 'text':
                            cy.wrap($el).clear().type(data.text);
                            break;
                        case 'email':
                            cy.wrap($el).clear().type(data.email);
                            break;
                        case 'tel':
                            cy.wrap($el).clear().type(data.phone);
                            break;
                        case 'url':
                            cy.wrap($el).clear().type(data.url);
                            break;
                        case 'date':
                            cy.wrap($el).clear().type(data.date);
                            break;
                        case 'checkbox':
                            cy.wrap($el).each(ck => {
                                if (ck.val() === data.checkbox) cy.wrap(ck).check();
                                else cy.wrap(ck).uncheck();
                            });
                            break;
                        case 'radio':
                            cy.wrap($el).each(ck => {
                                if (ck.val() === data.radio) cy.wrap(ck).check();
                            });
                    }
                } else if (tag === 'SELECT') {
                    cy.wrap($el).select(data.dropdown);
                } else if (tag === 'TEXTAREA') {
                    cy.wrap($el).clear().type(data.textarea);
                }
            });
        }
        return this;
    };

    submitForm = () => {
        return cy.get(ID.SUBMIT_BTN).click();
    };
}
