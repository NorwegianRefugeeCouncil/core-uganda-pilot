const PARTYSELECTOR = '[data-cy=partySelector]';
const ADD_BUTTON = '[data-cy=btn-add]';
const PARTYITEM = '[data-cy=partySelectorItem]';
const MEMBER = '[data-cy=member]';

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
