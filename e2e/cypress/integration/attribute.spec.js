import NewAttributePage from '../pages/newAttribute.page';
import AttributeOverviewPage from '../pages/attributeOverview.page';

const TYPE_NAME = 'New Test';
const EDITED_TYPE_NAME = 'Test Edited';

describe('Attribute Page', function () {
    describe('Create', () => {
        it('should navigate to New Attribute Form Page and save new attribute', () => {
            const newAttributePage = new NewAttributePage();
            newAttributePage
                .visitPage()
                .typeName(TYPE_NAME)
                .selectValueType('String')
                .selectSubjectType('Individual')
                .selectLanguage('English')
                .save();

            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage
                .visitPage()
                .selectAttribute()
                .should('contain.text', TYPE_NAME);
        });
    });

    describe('Update', () => {
        it('should update name on existing Attribute', () => {
            var attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage
                .visitAttribute()
                .clearName()
                .typeName(EDITED_TYPE_NAME)
                .save();

            attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage
                .visitPage()
                .selectAttribute()
                .should('contain.text', EDITED_TYPE_NAME);
        });
    });
});
