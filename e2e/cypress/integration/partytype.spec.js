import NewPartyTypePage from '../pages/newPartyType.page';
import PartytypesOverviewPage from '../pages/partytypeOverview.page';

const TYPE_NAME = 'New Test';
const EDITED_TYPE_NAME = 'Test Edited';

describe('PartyType Page', function () {
    describe('Create', () => {
        it('should navigate to New PartyType Form Page and save new attribute', () => {
            const newPartyTypePage = new NewPartyTypePage();
            newPartyTypePage
                .visitPage()
                .typeName(TYPE_NAME)
                .checkIsBuiltIn()
                .save();

            const partytypesOverviewPage = new PartytypesOverviewPage();
            partytypesOverviewPage
                .visitPage()
                .selectPartytype()
                .should('contain.text', TYPE_NAME);
        });
    });

    describe('Update', () => {
        it('should update name on existing PartyType', () => {
            var partytypesOverviewPage = new PartytypesOverviewPage();
            partytypesOverviewPage
                .visitPartytype()
                .clearName()
                .typeName(EDITED_TYPE_NAME)
                .save();

            partytypesOverviewPage = new PartytypesOverviewPage();
            partytypesOverviewPage
                .visitPage()
                .selectPartytype()
                .should('contain.text', EDITED_TYPE_NAME);
        });
    });
});
