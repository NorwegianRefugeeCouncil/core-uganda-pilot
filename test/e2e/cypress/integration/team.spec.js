import TeamsOverviewPage from '../pages/teamOverview.page';
import '../support/commands';

const SEARCH_NAME = 'DIKUA, Robert';
const MEMBER_SHOWN_NAME = 'Robert Dikua';

describe('Teams Page', function () {
    beforeEach('Login', () => {
        cy.login('courtney.lare@email.com');
    });
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
