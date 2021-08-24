import IndividualPage from '../pages/individualPage';
import IndividualOverviewPage from '../pages/individualOverview.page';

const DATA = {
    ATTRIBUTE_NAME: 'Test',
    EDITED_ATTRIBUTE_NAME: 'Edited',
    RELATIONSHIP_TYPE: 'Is spouse of',
    RELATED_PARTY: 'DOE',
};
const ATTRIBUTE_NAME = 'Test';
const EDITED_ATTRIBUTE_NAME = 'Edited';

// FIXME rewrite
describe.skip('Individual Page', function () {
    describe('Create', () => {
        it('should navigate to New Individual Form Page and save new Individual', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage
                .visitPage()
                .typeTextAttributes(ATTRIBUTE_NAME)
                .selectRelationshipType(DATA.RELATIONSHIP_TYPE)
                .typeRelatedParty(DATA.RELATED_PARTY)
                .addParty()
                .save();
        });
    });

    describe('Verify new Individual', () => {
        it('should verify that the Individual was properly created', () => {
            const individualOverviewPage = new IndividualOverviewPage().visitPage();
            const individualPage = new IndividualPage(individualOverviewPage.searchFor(ATTRIBUTE_NAME));
            // Verify values
            individualPage.verifyTextAttributes(ATTRIBUTE_NAME);
            individualPage.getRelationship().should('contain.text', DATA.RELATED_PARTY);
        });
    });

    // TODO update this after COR-187 is done
    describe.skip('Update', () => {
        it('should update attribute name on existing Individual', () => {
            const individualOverviewPage = new IndividualOverviewPage().visitPage();
            const individualPage = new IndividualPage(individualOverviewPage.searchFor(ATTRIBUTE_NAME));
            individualPage.typeTextAttributes(EDITED_ATTRIBUTE_NAME).save();
        });
    });
});
