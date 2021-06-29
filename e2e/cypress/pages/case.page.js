const DESCRIPTION = '[data-testid=description]'
const UPDATE_FORM_Btn = '[data-testid=updateFormBtn]'
const REFERRAL_PICKER = '[data-testid=referralPicker]'
const REFERRAL_CASES = '[data-testid=referralCaseItem]'
const REFERRAL_DESCRIPTION = '[data-testid=referralDescription]'
const SUBMIT_REFERRAL_BUTTON = '[data-testid=submitReferralBtn]'
const REFERRAL_CASE_OPEN = '[data-testid=referralCaseOpen]'

export default class CasePage {

    clearDescriptionValue = () => {
        cy.get(DESCRIPTION).clear();
        return this;
    }

    typeDescription = (value) => {
        cy.get(DESCRIPTION).type(value);
        return this;
    }

    getDescriptionValue = () => {
        return cy.get(DESCRIPTION);
    }

    submitUpdate = () => {
        cy.get(UPDATE_FORM_Btn).click();
        return this;
    }

    selectReferral = () => {
        cy.get(REFERRAL_PICKER).click();
        cy.get(REFERRAL_CASES).first().click();
        return this;
    }

    typeReferralDescription = (value) => {
        cy.get(REFERRAL_DESCRIPTION).type(value);
        return this;
    }

    submitReferral = () => {
        cy.get(SUBMIT_REFERRAL_BUTTON).click();
        return this;
    }

    getOpenReferralItem = () => {
        return cy.get(REFERRAL_CASE_OPEN);
    }
}
