import NewCaseTypePage from '../pages/newCasetype.page';
import { caseTypeTemplate } from '../helpers';

describe('CaseType Page', function () {
    describe('Create', () => {
        it('should navigate to New CaseType Form Page and save new attribute', () => {
            const newCaseTypePage = new NewCaseTypePage();
            newCaseTypePage
                .visitPage()
                .typeName('Test')
                .selectPartyType('Individual')
                .selectTeam('Kampala Response Team')
                .typeTemplate(caseTypeTemplate)
                .save();
        });
    });
});
