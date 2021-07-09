import NewPartyTypePage from '../pages/newPartyType.page';

describe('PartyType Page', function () {
    describe('Create', () => {
        it('should navigate to New PartyType Form Page and save new attribute', () => {
            const newPartyTypePage = new NewPartyTypePage();
            newPartyTypePage
                .visitPage()
                .typeName('Test')
                .checkIsBuiltIn()
                .save();
        });
    });
});
