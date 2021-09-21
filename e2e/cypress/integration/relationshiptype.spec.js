import NewRelationshiptypePage from '../pages/newRelationshiptype.page';
import RelationshiptypeOverviewPage from '../pages/relationshiptypeOverview.page';

const TYPE_NAME = 'New Test';
const EDITED_TYPE_NAME = 'Test Edited';

describe('Relationshiptype Page', function () {
    before('Login', () => {
        cy.login('courtney.lare@email.com');
    });
    describe('Create', () => {
        it('should navigate to New Relationshiptype Form Page and save new attribute', () => {
            const newRelationshiptypePage = new NewRelationshiptypePage();
            newRelationshiptypePage
                .visitPage()
                .typeName(TYPE_NAME)
                .checkIsDirectional()
                .typeFirstPartyRole('Is head of household of')
                .typeSecondPartyRole('Has for head of household')
                .selectFristPartyType('Individual')
                .selectSecondPartyType('Household')
                .save();

            const relationshiptypeOverviewPage = new RelationshiptypeOverviewPage();
            relationshiptypeOverviewPage.visitPage().selectRelationshiptype().should('contain.text', TYPE_NAME);
        });
    });

    describe('Update', () => {
        it('should update name on existing Relationshiptype', () => {
            var relationshiptypeOverviewPage = new RelationshiptypeOverviewPage();
            relationshiptypeOverviewPage.visitRelationshiptype().clearName().typeName(EDITED_TYPE_NAME).save();

            relationshiptypeOverviewPage = new RelationshiptypeOverviewPage();
            relationshiptypeOverviewPage.visitPage().selectRelationshiptype().should('contain.text', EDITED_TYPE_NAME);
        });
    });
});
