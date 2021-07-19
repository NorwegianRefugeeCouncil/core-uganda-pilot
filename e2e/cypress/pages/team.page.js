const PARTYSELECTOR = '[data-testid=partySelector]';
const ADD_BUTTON = '[data-testid=btn-add]';
const PARTYITEM = '[data-testid=partySelectorItem]';
const MEMBER = '[data-testid=member]';

export default class TeamPage {
    typeParty = (value) => {
        cy.get(PARTYSELECTOR).type(value);
        return this;
    };

    add = () => {
        cy.get(ADD_BUTTON).click();
        return this;
    };

    selectParty = () => {
        cy.get(PARTYITEM).click();
        return this;
    };

    selectTeamMembers = () => {
        return cy.get(MEMBER);
    };
}
