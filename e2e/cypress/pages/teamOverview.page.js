import { Urls } from '../helpers';
import TeamPage from './team.page';

const TEAM_ROWS = '[data-cy=team]';

export default class TeamsOverviewPage {
    visitPage = () => {
        cy.log('navigating to %s', Urls.TEAMS_URL);
        cy.visit(Urls.TEAMS_URL);
        return this;
    };

    visitTeam = () => {
        cy.get(TEAM_ROWS).first().click();
        return new TeamPage();
    };
}
