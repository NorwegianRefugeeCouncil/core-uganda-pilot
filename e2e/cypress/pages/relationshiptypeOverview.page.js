import { Urls } from '../helpers';
import NewRelationshiptypePage from './newRelationshiptype.page';

const RELATIONSHIPTYPE_ROWS = '[data-cy=relationshiptype]';

export default class RelationshiptypeOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.RELATIONSHIPTYPE_URL);
        cy.visit(Urls.RELATIONSHIPTYPE_URL);
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
