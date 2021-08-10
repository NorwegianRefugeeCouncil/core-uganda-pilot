import AttributePage from '../pages/attributePage';
import AttributeOverviewPage from '../pages/attributeOverview.page';
import ids from '../fixtures/ids.json';

const OPTIONS = {
    NAME: 'Test attribute',
    NAME_U: 'Test attribute - updated',
    VALUE_TYPE: 'String',
    SUBJECT_TYPE: ids.IndividualPartyTypeID,
    LANGUAGE_1: 'en',
    TRANSLATION_1_L: 'Long translation',
    TRANSLATION_1_L_U: 'Long translation - updated',
    TRANSLATION_1_S: 'Short translation',
    TRANSLATION_1_S_U: 'Short translation - updated',
    LANGUAGE_2: 'fr',
    TRANSLATION_2_L: 'Traduction longue',
    TRANSLATION_2_S: 'Traduction courte',
};

describe('Attribute Page', function () {
    describe('Navigate', () => {
        it('should navigate to the New Attribute page from the attributes page', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.clickNewAttributeBtn().url().should('include', 'attributes/new');
        });
    });
    describe('Create', () => {
        it('should create a new attribute', () => {
            const newAttributePage = new AttributePage();
            newAttributePage
                .setName(OPTIONS.NAME)
                .selectValueType(OPTIONS.VALUE_TYPE)
                .selectSubjectType(OPTIONS.SUBJECT_TYPE)
                .selectLanguage(OPTIONS.LANGUAGE_1)
                .setTranslationLong(OPTIONS.LANGUAGE_1, OPTIONS.TRANSLATION_1_L)
                .setTranslationShort(OPTIONS.LANGUAGE_1, OPTIONS.TRANSLATION_1_S)
                .selectLanguage(OPTIONS.LANGUAGE_2)
                .setTranslationLong(OPTIONS.LANGUAGE_2, OPTIONS.TRANSLATION_2_L)
                .setTranslationShort(OPTIONS.LANGUAGE_2, OPTIONS.TRANSLATION_2_S)
                .save();
        });
    });

    describe('Verify creation', () => {
        it('should verify that the attribute was created properly', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.selectLastAttribute().should('contain.text', OPTIONS.NAME);
            const attrPage = attributeOverviewPage.attributePageForNewest();

            // Verify values
            attrPage.getName().should('have.value', OPTIONS.NAME);
            attrPage.getValueType().should('have.value', OPTIONS.VALUE_TYPE);
            attrPage.getSubjectType().should('have.value', OPTIONS.SUBJECT_TYPE);
            attrPage.getPersonalInfo().should('not.be.checked');
            attrPage.getLanguageDsp(OPTIONS.LANGUAGE_1).should('exist');
            attrPage.getTranslationLong(OPTIONS.LANGUAGE_1).should('have.value', OPTIONS.TRANSLATION_1_L);
            attrPage.getTranslationShort(OPTIONS.LANGUAGE_1).should('have.value', OPTIONS.TRANSLATION_1_S);
            attrPage.getLanguageDsp(OPTIONS.LANGUAGE_2).should('exist');
            attrPage.getTranslationLong(OPTIONS.LANGUAGE_2).should('have.value', OPTIONS.TRANSLATION_2_L);
            attrPage.getTranslationShort(OPTIONS.LANGUAGE_2).should('have.value', OPTIONS.TRANSLATION_2_S);
        });
    });

    describe('Update', () => {
        it('should update name on existing Attribute', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            const attrPage = attributeOverviewPage.attributePageForNewest();

            // Update values

            attrPage.setName(OPTIONS.NAME_U);
            // TODO attrPage.selectValueType();
            attrPage.getPersonalInfo().check();
            attrPage.setTranslationLong(OPTIONS.LANGUAGE_1, OPTIONS.TRANSLATION_1_L_U);
            attrPage.setTranslationShort(OPTIONS.LANGUAGE_1, OPTIONS.TRANSLATION_1_S_U);
            attrPage.removeTranslation(OPTIONS.LANGUAGE_2);
            attrPage.save();
        });
    });

    describe('Verify update', () => {
        it('should verify that the attribute was updated properly', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.selectLastAttribute().should('contain.text', OPTIONS.NAME);
            const attrPage = attributeOverviewPage.attributePageForNewest();

            // Verify values
            attrPage.getName().should('have.value', OPTIONS.NAME_U);
            attrPage.getValueType().should('have.value', OPTIONS.VALUE_TYPE);
            attrPage.getSubjectType().should('have.value', OPTIONS.SUBJECT_TYPE);
            attrPage.getPersonalInfo().should('be.checked');
            attrPage.getLanguageDsp(OPTIONS.LANGUAGE_1).should('exist');
            attrPage.getTranslationLong(OPTIONS.LANGUAGE_1).should('have.value', OPTIONS.TRANSLATION_1_L_U);
            attrPage.getTranslationShort(OPTIONS.LANGUAGE_1).should('have.value', OPTIONS.TRANSLATION_1_S_U);
            attrPage.getLanguageDsp(OPTIONS.LANGUAGE_2).should('not.exist');
        });
    });
});
