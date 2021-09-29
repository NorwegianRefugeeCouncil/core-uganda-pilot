import { URL, testId } from '../helpers';
import testTemplate from '../fixtures/test_casetemplate.json';

const ID = {
    casetypeSelect: testId('casetype-select'),
    casetypeOption: testId('caseTypeOption'),
    partySelect: testId('party-select'),
    partyOption: testId('partyOption'),
    form: testId('form'),
    submitBtn: testId('submitBtn'),
};

export default class NewCasePage {
    constructor(noredirect) {
        if (!noredirect) {
            this.visitPage();
        }
    }

    visitPage = () => {
        cy.visit(URL.newCase);
        return this;
    };

    getCaseTypeSelect = () => cy.get(ID.casetypeSelect);
    setCaseType = value => {
        this.getCaseTypeSelect().select(value);
        return this;
    };

    getPartySelect = () => cy.get(ID.partySelect);
    setParty = value => {
        this.getPartySelect().select(value);
        return this;
    };

    fillOutForm = data => {
        for (const { type } of testTemplate.controls) {
            cy.get(testId('test-' + type)).then($el => {
                switch (type) {
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
                    case 'time':
                        cy.wrap($el).clear().type(data.time);
                        break;
                    case 'textarea':
                        cy.wrap($el).clear().type(data.textarea);
                        break;
                    case 'dropdown':
                        cy.wrap($el).select(data.dropdown);
                        break;
                    case 'boolean':
                        data.boolean ? cy.wrap($el).check() : cy.wrap($el).uncheck();
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
                        break;
                    case 'file':
                    case 'taxonomy':
                    default:
                        break;
                }
            });
        }
        return this;
    };

    submitForm = () => {
        return cy.get(ID.submitBtn).click().wait(200);
    };
}
