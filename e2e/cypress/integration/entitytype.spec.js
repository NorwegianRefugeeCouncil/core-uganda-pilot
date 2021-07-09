import NewEntityTypePage from '../pages/newEntityType.page';

describe('EntityType Page', function () {
    describe('Create', () => {
        it('should navigate to New EntityType Form Page and save new attribute', () => {
            const newEntityTypePage = new NewEntityTypePage();
            newEntityTypePage
                .visitPage()
                .typeName('Test')
                .checkIsBuiltIn()
                .save();
        });
    });
});
