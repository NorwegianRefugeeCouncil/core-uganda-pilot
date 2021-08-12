import { URL } from '../helpers';
import TeamPage from './team.page';

const TEAM_ROWS = '[data-testid=team]';

export default class TeamsOverviewPage {
    visitPage = () => {
        cy.visit(URL.TEAMS);
        return this;
    };

    visitTeam = () => {
        cy.get(TEAM_ROWS).first().click();
        return new TeamPage();
    };
}
