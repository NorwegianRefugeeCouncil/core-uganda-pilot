import { testId } from '../helpers';
import testTemplate from '../fixtures/test_casetemplate.json';

const ID = {
    FLASH: testId('flash'),
    RECIPIENT: testId('recipient'),
    DONE_PILL: testId('done-pill'),
    DONE_CHECK: testId('done-check'),
    SAVE_BTN: testId('save-btn'),
    REFERRAL_PICKER: testId('referralPicker'),
    REFERRAL_CASES: testId('referralCaseItem'),
    SUBMIT_REFERRAL_BUTTON: testId('submitReferralBtn'),
    REFERRAL_CASE_OPEN: testId('referralCaseOpen'),
    FORM: testId('form'),
};

export default class CasePage {
    constructor(href) {
        if (href != null) {
            href.then(h => cy.visit(h));
        } else {
            cy.visit(URL.newCase);
        }
    }

    getFlash = () => cy.get(ID.FLASH);

    getRecipient = () => cy.get(ID.RECIPIENT);
    getDonePill = () => cy.get(ID.DONE_PILL);

    clearForm = () => {
        cy.get(ID.FORM).each($el => {
            cy.wrap($el).clear();
        });
        return this;
    };

    verifyForm = value => {
        for (const { type } of testTemplate.controls) {
            cy.get(testId('test-' + type)).then($el => {
                switch (type) {
                    case 'text':
                        cy.wrap($el).should('have.value', value.text);
                        break;
                    case 'email':
                        cy.wrap($el).should('have.value', value.email);
                        break;
                    case 'tel':
                        cy.wrap($el).should('have.value', value.phone);
                        break;
                    case 'url':
                        cy.wrap($el).should('have.value', value.url);
                        break;
                    case 'date':
                        cy.wrap($el).should('have.value', value.date);
                        break;
                    case 'time':
                        cy.wrap($el).should('have.value', value.time);
                        break;
                    case 'textarea':
                        cy.wrap($el).should('have.value', value.textarea);
                        break;
                    case 'dropdown':
                        cy.wrap($el).should('have.value', value.dropdown);
                        break;
                    case 'boolean':
                        value.boolean ? cy.wrap($el).should('be.checked') : cy.wrap($el).should('not.be.checked');
                        break;
                    case 'checkbox':
                        value.checkbox ? cy.wrap($el).should('be.checked') : cy.wrap($el).should('not.be.checked');
                        break;
                    case 'radio':
                        value.radio ? cy.wrap($el).should('be.checked') : cy.wrap($el).should('not.be.checked');
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

    getDoneCheck = () => cy.get(ID.DONE_CHECK);
    setDoneCheck = () => {
        cy.get(ID.DONE_CHECK).check();
        return this;
    };

    save = () => cy.get(ID.SAVE_BTN).click().wait(200);

    typeForm = value => {
        cy.get(ID.FORM).each($el => {
            cy.wrap($el).type(value);
        });
        return this;
    };

    submitUpdate = () => {
        cy.get(ID.SAVE_BTN).click();
        return this;
    };

    selectReferral = () => {
        cy.get(ID.REFERRAL_PICKER).click();
        cy.get(ID.REFERRAL_CASES).first().click();
        return this;
    };

    submitReferral = () => {
        cy.get(ID.SUBMIT_REFERRAL_BUTTON).click();
        return this;
    };

    getOpenReferralItem = () => {
        return cy.get(ID.REFERRAL_CASE_OPEN);
    };
}
