import { testId, URL } from '../helpers';
import NewCasePage from './newCase.page';
import CasePage from './case.page';

const NEW_CASE_BTN = testId('new-case-btn');
const CASE_ROW = testId('case-row');

export default class CasesOverviewPage {
    constructor() {
        this.visitPage();
    }

    visitPage = () => {
        cy.visit(URL.cases);
        return this;
    };

    selectLastCase = () => cy.get(CASE_ROW).last();
    newestCasePage = () => {
        this.selectLastCase().invoke('attr', 'data-testhref').as('href');
        return new CasePage(cy.get('@href'));
    };

    openNewCase = () => {
        cy.get(NEW_CASE_BTN).click();
        return this;
    };

    selectCase = () => {
        cy.get(CASE_ROW).last().click();
        return new NewCasePage();
    };

    newCaseForm = () => {
        return new NewCasePage();
    };

    clickNewCaseBtn = () => {
        return cy.get(NEW_CASE_BTN).click();
    };
}
