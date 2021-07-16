import { Urls } from '../helpers';
import NewCaseTypePage from './newCasetype.page';

const CASETYPE_ROWS = '[data-cy=casetype]';

export default class CasetypesOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.CASETYPE_URL);
        cy.visit(Urls.CASETYPE_URL);
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
