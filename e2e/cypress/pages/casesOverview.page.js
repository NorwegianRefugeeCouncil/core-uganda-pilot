import { testId, URL } from '../helpers';
import NewCasePage from './newCase.page';

const NEW_CASE_BTN = testId('new-case-btn');
const CASE_ROWS = '[data-testid=caseRow]';

export default class CasesOverviewPage {
    constructor() {
        this.visitPage();
    }

    visitPage = () => {
        cy.visit(URL.CASES);
        return this;
    };

    openNewCase = () => {
        cy.get(NEW_CASE_BTN).click();
        return this;
    };

    selectCase = () => {
        cy.get(CASE_ROWS).last().click();
        return new NewCasePage();
    };

    newCaseForm = () => {
        return new NewCasePage();
    };

    clickNewCaseBtn = () => {
        return cy.get(NEW_CASE_BTN).click();
    };
}
