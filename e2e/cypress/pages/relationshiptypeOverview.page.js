import { URL } from '../helpers';
import NewRelationshiptypePage from './newRelationshiptype.page';

const RELATIONSHIPTYPE_ROWS = '[data-testid=relationshiptype]';

export default class RelationshiptypeOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', URL.RELATIONSHIPTYPE);
        cy.visit(URL.RELATIONSHIPTYPE);
        return this;
    };

    selectRelationshiptype = () => {
        return cy.get(RELATIONSHIPTYPE_ROWS).last();
    };

    visitRelationshiptype = () => {
        cy.get(RELATIONSHIPTYPE_ROWS).last().click();
        return new NewRelationshiptypePage();
    };
}
