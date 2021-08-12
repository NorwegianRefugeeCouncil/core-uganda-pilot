import { TEST_CASE_TEMPLATE_FIELD, testId } from '../helpers';

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
            cy.visit(href);
        } else {
            this.visitPage();
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
        for (const id of Object.values(TEST_CASE_TEMPLATE_FIELD)) {
            cy.get(id).then($el => {
                const tag = $el[0].tagName;
                switch (tag) {
                    case 'INPUT':
                        switch ($el[0].getAttribute('type')) {
                            case 'text':
                                cy.wrap($el).should('have.value', value.textinput);
                                break;
                            case 'checkbox':
                                value.checkbox
                                    ? cy.wrap($el).should('be.checked')
                                    : cy.wrap($el).should('not.be.checked');
                                break;
                            default:
                                break;
                        }
                        break;
                    case 'SELECT':
                        cy.wrap($el).should('have.value', value.dropdown);
                        break;
                    case 'TEXTAREA':
                        cy.wrap($el).should('have.value', value.textarea);
                        break;
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

    save = () => cy.get(ID.SAVE_BTN).click();

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
