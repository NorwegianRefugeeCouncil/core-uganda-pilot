import AttributePage from '../pages/attributePage';
import AttributeOverviewPage from '../pages/attributeOverview.page';
import ids from '../fixtures/ids.json';
import { URL } from '../helpers';

const DATA = {
    name: 'Test attribute',
    nameUpd: 'Test attribute - updated',
    label: 'label',
    labelUpd: 'label updated',
    controlType: 'text',
    controlTypeUpd: 'dropdown',
    subjectType: ids.IndividualPartyTypeID,
};

describe('Attribute Page', function () {
    describe('Navigate', () => {
        it('should navigate to the New Attribute page from the attributes page', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.clickNewAttributeBtn().url().should('include', URL.newAttribute);
        });
    });
    describe('Create', () => {
        it('should create a new attribute', () => {
            const newAttributePage = new AttributePage();
            newAttributePage
                .setName(DATA.name)
                .setLabel(DATA.label)
                .selectControlType(DATA.controlType)
                .selectSubjectType(DATA.subjectType)
                .save();
        });
    });

    describe('Verify creation', () => {
        it('should verify that the attribute was created properly', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.selectLastAttribute().should('contain.text', DATA.name);
            const attrPage = attributeOverviewPage.newestAttributePage();

            // Verify values
            attrPage.getName().should('have.value', DATA.name);
            attrPage.getControlType().should('have.value', DATA.controlType);
            attrPage.getSubjectType().should('contain.value', DATA.subjectType);
            attrPage.getRequired().should('not.be.checked');
            attrPage.getPersonalInfo().should('not.be.checked');
        });
    });

    describe('Update', () => {
        it('should update name on existing Attribute', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            const attrPage = attributeOverviewPage.newestAttributePage();

            // Update values
            attrPage
                .setName(DATA.nameUpd)
                .setLabel(DATA.labelUpd)
                .selectControlType(DATA.controlTypeUpd)
                .checkRequired()
                .checkPersonalInfo()
                .save();
        });
    });

    describe('Verify update', () => {
        it('should verify that the attribute was updated properly', () => {
            const attributeOverviewPage = new AttributeOverviewPage();
            attributeOverviewPage.selectLastAttribute().should('contain.text', DATA.name);
            const attrPage = attributeOverviewPage.newestAttributePage();

            // Verify values
            attrPage.getName().should('have.value', DATA.nameUpd);
            attrPage.getLabel().should('have.value', DATA.labelUpd);
            attrPage.getControlType().should('have.value', DATA.controlTypeUpd);
            attrPage.getSubjectType().should('contain.value', DATA.subjectType);
            attrPage.getRequired().should('be.checked');
            attrPage.getPersonalInfo().should('be.checked');
        });
    });
});
