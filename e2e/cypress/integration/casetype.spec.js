import NewCaseTypePage from '../pages/newCasetype.page';
import { caseTypeTemplate } from '../helpers';
import CasetypesOverviewPage from '../pages/casetypesOverview.page';

const TYPE_NAME = 'New Test';
const EDITED_TYPE_NAME = 'Test Edited';

describe('CaseType Page', function () {
    describe('Create', () => {
        it('should navigate to New CaseType Form Page and save new attribute', () => {
            const newCaseTypePage = new NewCaseTypePage();
            newCaseTypePage
                .visitPage()
                .typeName(TYPE_NAME)
                .selectPartyType('Individual')
                .selectTeam('Kampala Response Team')
                .typeTemplate(caseTypeTemplate)
                .save();

            const casetypesOverviewPage = new CasetypesOverviewPage();
            casetypesOverviewPage
                .visitPage()
                .selectCasetype()
                .should('contain.text', TYPE_NAME);
        });
    });

    describe('Update', () => {
        it('should update name on existing CaseType', () => {
            var casetypesOverviewPage = new CasetypesOverviewPage();
            casetypesOverviewPage
                .visitCasetype()
                .clearName()
                .typeName(EDITED_TYPE_NAME)
                .save();

            casetypesOverviewPage = new CasetypesOverviewPage();
            casetypesOverviewPage
                .visitPage()
                .selectCasetype()
                .should('contain.text', EDITED_TYPE_NAME);
        });
    });
});
