import NewIndividualPage from '../pages/newIndividual.page';
import IndividualOverviewPage from '../pages/individualOverview.page';

const ATTRIBUTE_NAME = 'Test';
const EDITED_ATTRIBUTE_NAME = 'Edited';

describe('Individual Page', function () {
    describe('Create', () => {
        it('should navigate to New Individual Form Page and save new attribute', () => {
            const newIndividualPage = new NewIndividualPage();
            newIndividualPage
                .visitPage()
                .typeTextAttributes(ATTRIBUTE_NAME)
                .selectRelationshipType('Is spouse of')
                .typeRelatedParty('Doe, John')
                .addParty()
                .save();

            const individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage
                .visitPage()
                .selectIndividual()
                .should('contain.text', ATTRIBUTE_NAME);
        });
    });

    describe('Update', () => {
        it('should update attribute name on existing Individual', () => {
            var individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage
                .visitIndividual()
                .clearTextAttributes()
                .typeTextAttributes(EDITED_ATTRIBUTE_NAME)
                .save();

            individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage
                .visitPage()
                .selectIndividual()
                .should('contain.text', EDITED_ATTRIBUTE_NAME);
        });
    });
});
