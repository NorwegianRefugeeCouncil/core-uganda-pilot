import { URL } from '../helpers';
import CasePage from './case.page';
import NewCasePage from './newCase.page';

const OPEN_NEW_CASE_Btn = '[data-testid=openNewCase]';
const CASE_ROWS = '[data-testid=caseRow]';

export default class CasesOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', URL.CASE);
        cy.visit(URL.CASES);
        return this;
    };

    openNewCase = () => {
        cy.get(OPEN_NEW_CASE_Btn).click();
        return this;
    };

    selectCase = () => {
        cy.get(CASE_ROWS).last().click();
        return new CasePage();
    };

    newCaseForm = () => {
        return new NewCasePage();
    };
}
