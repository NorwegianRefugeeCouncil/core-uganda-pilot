import { URL } from '../helpers';
import NewCaseTypePage from './newCasetype.page';

const CASETYPE_ROWS = '[data-testid=casetype]';

export default class CasetypesOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', URL.CASETYPE);
        cy.visit(URL.CASETYPE);
        return this;
    };

    selectCasetype = () => {
        return cy.get(CASETYPE_ROWS).last();
    };

    visitCasetype = () => {
        cy.get(CASETYPE_ROWS).last().click();
        return new NewCaseTypePage();
    };
}
