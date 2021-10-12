import {testId, URL} from '../helpers';

const NAME = testId('name');
const PARTY_TYPE = testId('partytype');
const TEAM = testId('team');
const TEMPLATE = testId('template');
const SAVE_BUTTON = testId('save-btn');

export default class CaseTypePage {
    constructor(href) {
        if (href != null) {
            href.then(h => cy.visit(h));
        } else {
            this.visitPage();
        }
    }

    visitPage = () => {
        cy.visit(URL.newCasetype);
        return this;
    };

    getName = () => cy.get(NAME);
    setName = value => {
        this.getName().clear().type(value);
        return this;
    };

    getPartyTypeSelect = () => cy.get(PARTY_TYPE);
    setPartyType = value => {
        this.getPartyTypeSelect().select(value);
        return this;
    };

    getTeamSelect = () => cy.get(TEAM);
    setTeam = value => {
        this.getTeamSelect().select(value);
        return this;
    };

    getTemplate = () => cy.get(TEMPLATE);
    setTemplate = value => {
        this.getTemplate().invoke('val', value);
        return this;
    };

    save = () => {
        cy.get(SAVE_BUTTON).click();
        return this;
    };

    clearName = () => {
        cy.get(NAME).clear();
        return this;
    };
}
