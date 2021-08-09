import AttributePage from '../pages/attributePage';
import AttributeOverviewPage from '../pages/attributeOverview.page';

const OPTIONS = {
    NAME: 'New Test',
    EDITED_NAME: 'Test Edited',
    VALUE_TYPE: 'String',
    SUBJECT_TYPE: 'a842e7cb-3777-423a-9478-f1348be3b4a5', // Individual PartyTypeID
    LANGUAGE: 'English',
    TRANSLATION_LONG: 'Long',
    TRANSLATION_SHORT: 'Short',
};

describe('Attribute Page', function () {
    describe('Create', () => {
        it('should navigate to New Attribute Form Page and save new attribute', () => {
            const newAttributePage = new AttributePage();
            newAttributePage
                .typeName(OPTIONS.NAME)
                .selectValueType(OPTIONS.VALUE_TYPE)
                .selectSubjectType(OPTIONS.SUBJECT_TYPE)
                .selectLanguage(OPTIONS.LANGUAGE)
                .setTranslationLong(OPTIONS.TRANSLATION_LONG)
                .setTranslationShort(OPTIONS.TRANSLATION_SHORT)
                .save();

            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage
                .selectAttribute()
                .should('contain.text', OPTIONS.NAME);

            // Verify values
            const attrPage = attributeOverviewPage.attributePageForLatest();
            attrPage.getName().should('have.value', OPTIONS.NAME);
            attrPage.getValueType().should('have.value', OPTIONS.VALUE_TYPE);
            attrPage
                .getSubjectType()
                .should('have.value', OPTIONS.SUBJECT_TYPE);
            attrPage.getPersonalInfo().should('not.be.checked');
            attrPage.getLanguageDsp().should('contain.text', OPTIONS.LANGUAGE);
            attrPage
                .getTranslationLong()
                .should('have.value', OPTIONS.TRANSLATION_LONG);
            attrPage
                .getTranslationShort()
                .should('have.value', OPTIONS.TRANSLATION_SHORT);
        });
    });

    describe('Update', () => {
        it('should update name on existing Attribute', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            const attrPage = attributeOverviewPage
                .attributePageForLatest()
                .clearName()
                .typeName(OPTIONS.EDITED_NAME)
                .save();
            attrPage.getName().should('have.value', OPTIONS.EDITED_NAME);
        });
    });
});
