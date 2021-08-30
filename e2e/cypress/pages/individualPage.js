import { URL, testId, nameAttr } from '../helpers';

const selector = {
    firstName: nameAttr('firstName'),
    lastName: nameAttr('lastName'),
    email: nameAttr('email'),
    birthDate: nameAttr('birthDate'),
    displacementStatus: nameAttr('displacementStatus'),
    gender: nameAttr('gender'),
    relationship: testId('relationship'),
    relationshipType: testId('relationship-type'),
    relatedParty: testId('related-party'),
    searchResult: testId('party-search-result'),
    addRelationshipBtn: testId('add-relationship-btn'),
    removeRelationshipBtn: testId('remove-relationship-btn'),
    saveBtn: testId('save-btn'),
    situationAnalysisTab: '#situation-analysis-tab',
    situationAnalysis: '#situation-analysis',
    textArea: testId('test-textarea'),
    taxonomyInput: testId('test-taxonomyinput'),
    addTaxonomyBtn: testId('add-taxonomy-btn'),
    taxonomyBadges: testId('badge-container'),
    perceivedPriority: nameAttr('perceivedPriority'),
    commentStartingPoint: nameAttr('commentStartingPoint'),
    commentOtherServices: nameAttr('commentOtherServices'),
    responseTab: '#response-tab',
    response: '#response',
};

export default class IndividualPage {
    constructor(href) {
        if (href) {
            href.then(h => cy.visit(h));
        }
    }

    visitPage = () => {
        cy.visit(URL.NEW_INDIVIDUAL);
        return this;
    };

    inputAttributes = data => {
        const { firstName, lastName, birthDate, email, status, gender, relationshipType, relatedParty } = data;
        this.typeFirstName(firstName);
        this.typeLastName(lastName);
        this.enterBirthDate(birthDate);
        this.typeEmail(email);
        this.selectDisplacementStatus(status);
        this.selectGender(gender);
        this.addRelationship({ relationshipType, relatedParty });
        return this;
    };
    verifyAttributes = data => {
        const { firstName, lastName, birthDate, email, status, gender, relationshipType, relatedParty } = data;
        this.getFirstName().should('have.value', firstName);
        this.getLastName().should('have.value', lastName);
        this.getBirthDate().should('have.value', birthDate);
        this.getEmail().should('have.value', email);
        this.getDisplacementStatus().should('have.value', status);
        this.getGender().should('have.value', gender);
        const r = this.getRelationShip();
        r.should('contain.text', relationshipType.toLowerCase());
        r.should('contain.text', relatedParty);
        return this;
    };

    getFirstName = () => cy.get(selector.firstName);
    getLastName = () => cy.get(selector.lastName);
    getEmail = () => cy.get(selector.email);
    getBirthDate = () => cy.get(selector.birthDate);
    getDisplacementStatus = () => cy.get(selector.displacementStatus);
    getGender = () => cy.get(selector.gender);
    getRelationshipType = () => cy.get(selector.relationshipType);
    getRelatedParty = () => cy.get(selector.relatedParty);
    getSearchResult = () => cy.get(selector.searchResult);
    getRelationShip = () => cy.get(selector.relationship).first();

    typeFirstName = value => this.getFirstName().clear().type(value);
    typeLastName = value => this.getLastName().clear().type(value);
    typeEmail = value => this.getEmail().clear().type(value);
    enterBirthDate = value => this.getBirthDate().invoke('val', value);
    selectDisplacementStatus = value => this.getDisplacementStatus().select(value);
    selectGender = value => this.getGender().select(value);
    selectRelationshipType = value => this.getRelationshipType().select(value);
    typeRelatedParty = value => this.getRelatedParty().clear().type(value).wait(500);
    addRelatedParty = () => this.getSearchResult().click();
    addRelationship = ({ relationshipType, relatedParty }) => {
        this.selectRelationshipType(relationshipType);
        this.typeRelatedParty(relatedParty);
        this.addRelatedParty();
        cy.get(selector.addRelationshipBtn).click();
    };
    removeRelationship = () => {
        this.getRelationShip().find(selector.removeRelationshipBtn).click();
        return this;
    };

    visitSituationAnalysisTab = () => {
        cy.get(selector.situationAnalysisTab).click();
        return this;
    };
    getSituationAnalysis = () => cy.get(selector.situationAnalysis);
    fillOutSituationAnalysis = data => {
        this.getSituationAnalysis()
            .find(selector.textArea)
            .each($t => cy.wrap($t).clear().type(data));
        return this;
    };
    verifySituationAnalysis = data => {
        this.getSituationAnalysis()
            .find(selector.textArea)
            .each($t => cy.wrap($t).should('have.value', data));
    };

    visitResponseTab = () => {
        cy.get(selector.responseTab).click();
        return this;
    };

    getResponse = () => cy.get(selector.response);

    fillOutResponse = data => {
        this.getResponse().within(() => {
            cy.get(selector.taxonomyInput).each($t => this.fillTaxonomyInput($t, data.optionIdx));
            cy.get(selector.perceivedPriority).clear().type(data.priorityTxt);
            cy.get(selector.commentStartingPoint).clear().type(data.commentStartingPoint);
            cy.get(selector.commentOtherServices).clear().type(data.commentOtherServices);
        });
        return this;
    };

    fillTaxonomyInput = ($t, i) => {
        cy.wrap($t).within(() => {
            cy.get('select').each($s => {
                cy.wrap($s).find('option').eq(i).as('opt');
                cy.get('@opt')
                    .invoke('attr', 'value')
                    .then(v => cy.wrap($s).select(v));
            });
            cy.get(selector.addTaxonomyBtn).click();
        });
    };

    verifyResponse = data => {
        this.getResponse().within(() => {
            cy.get(selector.taxonomyInput).each($t => this.verifyTaxonomyInput($t, data.optionIdx));
            cy.get(selector.perceivedPriority).should('have.value', data.priorityTxt);
            cy.get(selector.commentStartingPoint).should('have.value', data.commentStartingPoint);
            cy.get(selector.commentOtherServices).should('have.value', data.commentOtherServices);
        });
    };

    verifyTaxonomyInput = $t => {
        cy.wrap($t).find(selector.taxonomyBadges).should('have.length', 1);
    };

    save = () => {
        cy.get(selector.saveBtn).click();
    };
}
