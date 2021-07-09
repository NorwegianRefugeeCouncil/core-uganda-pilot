import NewIndividualPage from '../pages/newIndividual.page';

describe('Individual Page', function () {
    describe('Create', () => {
        it('should navigate to New Individual Form Page and save new attribute', () => {
            const newIndividualPage = new NewIndividualPage();
            newIndividualPage
                .visitPage()
                .typeTextAttributes('Test')
                .selectRelationshipType('Is spouse of')
                .typeRelatedParty('Doe, John')
                .addParty()
                .save();
        });
    });
});
