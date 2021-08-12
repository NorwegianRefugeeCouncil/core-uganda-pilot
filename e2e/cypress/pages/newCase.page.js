import { URL, testId, TEST_CASE_TEMPLATE_FIELD } from '../helpers';

const ID = {
    CASETYPE_SELECT: testId('casetype-select'),
    CASETYPE_OPTION: testId('caseTypeOption'),
    PARTY_SELECT: testId('party-select'),
    PARTY_OPTION: testId('partyOption'),
    FORM: testId('form'),
    SUBMIT_BTN: testId('submitBtn')
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

    fillOutForm = value => {
        for (const id of Object.values(TEST_CASE_TEMPLATE_FIELD)) {
            cy.get(id).then($el => {
                const tag = $el[0].tagName;
                if (tag === 'INPUT') {
                    switch ($el[0].getAttribute('type')) {
                        case 'text':
                            cy.wrap($el).clear().type(value.text);
                            break;
                        case 'checkbox':
                            cy.wrap($el).each(ck => {
                                return value.checkbox ? cy.wrap(ck).check() : cy.wrap(ck).uncheck();
                            });
                    }
                } else if (tag === 'SELECT') {
                    cy.wrap($el).select(value.dropdown);
                } else if (tag === 'TEXTAREA') {
                    cy.wrap($el).clear().type(value.textarea);
                }
            });
        }
        return this;
    };

    submitForm = () => {
        return cy.get(ID.SUBMIT_BTN).click();
    };
}
