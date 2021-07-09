import NewAttributePage from '../pages/newAttribute.page';

describe('Attribute Page', function () {
    describe('Create', () => {
        it('should navigate to New Attribute Form Page and save new attribute', () => {
            const newAttributePage = new NewAttributePage();
            newAttributePage
                .visitPage()
                .typeName('Test')
                .selectValueType('String')
                .selectSubjectType('Individual')
                .selectLanguage('English')
                .save();
        });
    });
});
