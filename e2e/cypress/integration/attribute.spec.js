import AttributePage from '../pages/attributePage';
import AttributeOverviewPage from '../pages/attributeOverview.page';
import ids from '../fixtures/ids.json';
import { URL } from '../helpers';
import './commands';

const DATA = {
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

// FIXME COR-204
describe.skip('Attribute Page', function () {
    before('Login', () => {
        cy.login('courtney.lare@email.com');
    });
    describe('Navigate', () => {
        it('should navigate to the New Attribute page from the attributes page', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.clickNewAttributeBtn().url().should('include', URL.NEW_ATTRIBUTE);
        });
    });
    describe('Create', () => {
        it('should create a new attribute', () => {
            const newAttributePage = new AttributePage();
            newAttributePage
                .setName(DATA.NAME)
                .selectValueType(DATA.VALUE_TYPE)
                .selectSubjectType(DATA.SUBJECT_TYPE)
                .selectLanguage(DATA.LANGUAGE_1)
                .setTranslationLong(DATA.LANGUAGE_1, DATA.TRANSLATION_1_L)
                .setTranslationShort(DATA.LANGUAGE_1, DATA.TRANSLATION_1_S)
                .selectLanguage(DATA.LANGUAGE_2)
                .setTranslationLong(DATA.LANGUAGE_2, DATA.TRANSLATION_2_L)
                .setTranslationShort(DATA.LANGUAGE_2, DATA.TRANSLATION_2_S)
                .save();
        });
    });

    describe('Verify creation', () => {
        it('should verify that the attribute was created properly', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.selectLastAttribute().should('contain.text', DATA.NAME);
            const attrPage = attributeOverviewPage.attributePageForNewest();

            // Verify values
            attrPage.getName().should('have.value', DATA.NAME);
            attrPage.getValueType().should('have.value', DATA.VALUE_TYPE);
            attrPage.getSubjectType().should('have.value', DATA.SUBJECT_TYPE);
            attrPage.getPersonalInfo().should('not.be.checked');
            attrPage.getLanguageDsp(DATA.LANGUAGE_1).should('exist');
            attrPage.getTranslationLong(DATA.LANGUAGE_1).should('have.value', DATA.TRANSLATION_1_L);
            attrPage.getTranslationShort(DATA.LANGUAGE_1).should('have.value', DATA.TRANSLATION_1_S);
            attrPage.getLanguageDsp(DATA.LANGUAGE_2).should('exist');
            attrPage.getTranslationLong(DATA.LANGUAGE_2).should('have.value', DATA.TRANSLATION_2_L);
            attrPage.getTranslationShort(DATA.LANGUAGE_2).should('have.value', DATA.TRANSLATION_2_S);
        });
    });

    describe('Update', () => {
        it('should update name on existing Attribute', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            const attrPage = attributeOverviewPage.attributePageForNewest();

            // Update values

            attrPage.setName(DATA.NAME_U);
            // TODO attrPage.selectValueType();
            attrPage.getPersonalInfo().check();
            attrPage.setTranslationLong(DATA.LANGUAGE_1, DATA.TRANSLATION_1_L_U);
            attrPage.setTranslationShort(DATA.LANGUAGE_1, DATA.TRANSLATION_1_S_U);
            attrPage.removeTranslation(DATA.LANGUAGE_2);
            attrPage.save();
        });
    });

    describe('Verify update', () => {
        it('should verify that the attribute was updated properly', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.selectLastAttribute().should('contain.text', DATA.NAME);
            const attrPage = attributeOverviewPage.attributePageForNewest();

            // Verify values
            attrPage.getName().should('have.value', DATA.NAME_U);
            attrPage.getValueType().should('have.value', DATA.VALUE_TYPE);
            attrPage.getSubjectType().should('have.value', DATA.SUBJECT_TYPE);
            attrPage.getPersonalInfo().should('be.checked');
            attrPage.getLanguageDsp(DATA.LANGUAGE_1).should('exist');
            attrPage.getTranslationLong(DATA.LANGUAGE_1).should('have.value', DATA.TRANSLATION_1_L_U);
            attrPage.getTranslationShort(DATA.LANGUAGE_1).should('have.value', DATA.TRANSLATION_1_S_U);
            attrPage.getLanguageDsp(DATA.LANGUAGE_2).should('not.exist');
        });
    });
});
