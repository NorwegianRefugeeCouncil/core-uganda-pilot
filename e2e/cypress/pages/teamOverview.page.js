import { testId, URL } from '../helpers';
import TeamPage from './team.page';

const TEAM_ROWS = testId('team');

export default class TeamsOverviewPage {
    visitPage = () => {
        cy.visit(URL.TEAMS);
        return this;
    };

    visitTeam = () => {
        cy.get(TEAM_ROWS).last().click();
        return new TeamPage();
    };
}
