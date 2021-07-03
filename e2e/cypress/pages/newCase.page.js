import { Urls } from '../helpers';
import CasePage from './case.page';

const NEW_CASE_FORM = '[data-testid=newCaseForm]';
const CASE_TYPE_PICKER = '[data-testid=caseTypePicker]';
const CASE_TYPE_OPTIONS = '[data-testid=caseTypeOption]';
const PARTY_PICKER = '[data-testid=partyPicker]';
const PARTY_OPTIONS = '[data-testid=partyOption]';
const DESCRIPTION = '[data-testid=description]';
const SUBMIT_BUTTON = '[data-testid=submitBtn]';

export default class NewCasePage {
    visitPage = () => {
        cy.visit(Urls.NEW_CASE_URL);
        return this;
    };

    selectCaseType = () => {
        cy.get(CASE_TYPE_OPTIONS)
            .first()
            .invoke('attr', 'value')
            .then((value) => cy.get(CASE_TYPE_PICKER).select(value));
        return this;
    };

    selectParty = () => {
        cy.get(PARTY_OPTIONS)
            .first()
            .invoke('attr', 'value')
            .then((value) => cy.get(PARTY_PICKER).select(value));
        return this;
    };

    typeDescription = (value) => {
        cy.get(DESCRIPTION).type(value);
        return this;
    };

    submitForm = () => {
        cy.get(SUBMIT_BUTTON).click();
        return this;
    };

    getNewCase = () => {
        const casePage = new CasePage();
        return casePage.getDescriptionValue();
    };

    verifyForm = () => {
        return cy.get(NEW_CASE_FORM);
    };
}