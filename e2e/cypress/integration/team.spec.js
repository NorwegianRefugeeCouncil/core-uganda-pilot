import TeamsOverviewPage from '../pages/teamOverview.page';

const SEARCH_NAME = 'MALVINA, Devi';
const MEMBER_SHOWN_NAME = 'Devi Malvina';

describe('Teams Page', function () {
    describe('Create', () => {
        it('should navigate to Team Page, select the first team and add a new member', () => {
            const teamsOverviewPage = new TeamsOverviewPage();
            teamsOverviewPage
                .visitPage()
                .visitTeam()
                .typeParty(SEARCH_NAME)
                .selectParty()
                .add()
                .selectTeamMembers()
                .should('contain.text', MEMBER_SHOWN_NAME);
        });
    });
});
