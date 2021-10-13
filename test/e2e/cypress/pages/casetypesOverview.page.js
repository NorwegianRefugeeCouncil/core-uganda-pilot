import {testId, URL} from '../helpers';
import CaseTypePage from './caseTypePage';

const CASETYPE_ROWS = testId('casetype');
const NEW_CASETYPE_BTN = testId('new-casetype-btn');

export default class CasetypesOverviewPage {
    constructor() {
        this.visitPage();
    }

    visitPage = () => {
        cy.visit(URL.casetypes);
        return this;
    };

    clickNewCaseTypeBtn = () => {
        return cy.get(NEW_CASETYPE_BTN).click();
    };

    selectNewestCasetype = () => {
        return cy.get(CASETYPE_ROWS).last();
    };

    visitCasetype = () => {
        cy.get(CASETYPE_ROWS).last().click();
        return new CaseTypePage();
    };

    caseTypePageForNewest = () => {
        cy.get(CASETYPE_ROWS).last().invoke('attr', 'href').as('href');
        return new CaseTypePage(cy.get('@href'));
    };
}
