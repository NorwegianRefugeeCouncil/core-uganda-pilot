import { Urls } from '../helpers';
import NewPartyTypePage from './newPartyType.page';

const PARTYTYPE_ROWS = '[data-cy=partytype]';

export default class CasetypesOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.PARTYTYPE_URL);
        cy.visit(Urls.PARTYTYPE_URL);
        return this;
    };

    selectPartytype = () => {
        return cy.get(PARTYTYPE_ROWS).last();
    };

    visitPartytype = () => {
        cy.get(PARTYTYPE_ROWS).last().click();
        return new NewPartyTypePage();
    };
}
