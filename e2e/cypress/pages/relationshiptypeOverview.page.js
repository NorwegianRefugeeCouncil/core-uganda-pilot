import { URL } from '../helpers';
import NewRelationshiptypePage from './newRelationshiptype.page';

const RELATIONSHIPTYPE_ROWS = '[data-testid=relationshiptype]';

export default class RelationshiptypeOverviewPage {
    visitPage = () => {
        cy.visit(URL.RELATIONSHIPTYPES);
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
