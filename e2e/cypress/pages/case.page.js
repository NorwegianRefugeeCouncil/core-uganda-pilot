const ALERT = '[data-testid=alert]';
const UPDATE_FORM_Btn = '[data-testid=updateFormBtn]';
const REFERRAL_PICKER = '[data-testid=referralPicker]';
const REFERRAL_CASES = '[data-testid=referralCaseItem]';
const SUBMIT_REFERRAL_BUTTON = '[data-testid=submitReferralBtn]';
const REFERRAL_CASE_OPEN = '[data-testid=referralCaseOpen]';
const FORM = '[data-testid=form]';

export default class CasePage {
    getAlertMessage = () => {
        return cy.get(ALERT);
    };

    clearForm = () => {
        cy.get(FORM).each(($el) => {
            cy.wrap($el).clear();
        });
        return this;
    };

    typeForm = (value) => {
        cy.get(FORM).each(($el) => {
            cy.wrap($el).type(value);
        });
        return this;
    };

    submitUpdate = () => {
        cy.get(UPDATE_FORM_Btn).click();
        return this;
    };

    selectReferral = () => {
        cy.get(REFERRAL_PICKER).click();
        cy.get(REFERRAL_CASES).first().click();
        return this;
    };

    submitReferral = () => {
        cy.get(SUBMIT_REFERRAL_BUTTON).click();
        return this;
    };

    getOpenReferralItem = () => {
        return cy.get(REFERRAL_CASE_OPEN);
    };
}
