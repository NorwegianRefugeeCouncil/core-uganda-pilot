import {nameAttr, testId, URL} from '../helpers';

const selector = {
    fullName: nameAttr('fullName'),
    displayName: nameAttr('displayName'),
    email: nameAttr('email'),
    birthDate: nameAttr('birthDate'),
    displacementStatus: nameAttr('displacementStatus'),
    gender: nameAttr('gender'),
    identificationDate: nameAttr('ugIdentificationDate'),
    identificationLocation: nameAttr('ugIdentificationLocation'),
    identificationSource: nameAttr('ugIdentificationSource'),
    admin2: nameAttr('ugAdmin2'),
    admin3: nameAttr('ugAdmin3'),
    admin4: nameAttr('ugAdmin4'),
    admin5: nameAttr('ugAdmin5'),
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
    taxonomyInput: testId('test-taxonomy'),
    addTaxonomyBtn: testId('add-taxonomy-btn'),
    taxonomyBadges: testId('badge-container'),
    perceivedPriority: nameAttr('perceivedPriority'),
    commentStartingPoint: nameAttr('commentStartingPoint'),
    commentOtherServices: nameAttr('commentOtherServices'),
    responseTab: '#response-tab',
    response: '#response',
    doneSituationAnalysis: testId('done-check-situation'),
    doneResponse: testId('done-check-response'),
};

const ugandaAttributeNames = [
    nameAttr('ugIdentificationDate'),
    nameAttr('ugIdentificationLocation'),
    nameAttr('ugIdentificationSource'),
    nameAttr('ugAdmin2'),
    nameAttr('ugAdmin3'),
    nameAttr('ugAdmin4'),
    nameAttr('ugAdmin5'),
    nameAttr('ugNationality'),
    nameAttr('ugSpokenLanguages'),
    nameAttr('ugPreferredLanguage'),
    nameAttr('ugPhysicalAddress'),
    nameAttr('ugInstructionOnMakingContact'),
    nameAttr('ugCanInitiateContact'),
    nameAttr('ugPreferredMeansOfContact'),
    nameAttr('ugRequireAnInterpreter'),
];

const colombiaAttributeNames = [
    nameAttr('coPrimaryNationality'),
    nameAttr('coSecondaryNationality'),
    nameAttr('coMaritalStatus'),
    nameAttr('coBeneficiaryType'),
    nameAttr('coEthnicity'),
    nameAttr('coRegistrationDate'),
    nameAttr('coRegistrationLocation'),
    nameAttr('coSourceOfIdentification'),
    nameAttr('coTypeOfSettlement'),
    nameAttr('coEmergencyCare'),
    nameAttr('coDurableSolutions'),
    nameAttr('coHardToReach'),
    nameAttr('coAttendedCovid19'),
    nameAttr('coIntroSource'),
    nameAttr('coAdmin1'),
    nameAttr('coAdmin2'),
    nameAttr('coAdmin3'),
    nameAttr('coAdmin4'),
    nameAttr('coAdmin5'),
    nameAttr('coJobOrEnterprise'),
    nameAttr('coTypeOfEnterprise'),
    nameAttr('coTimeInBusiness'),
    nameAttr('coTypeOfEmployment'),
    nameAttr('coFormsOfIncomeGeneration'),
    nameAttr('coLegalRepresentativeName'),
    nameAttr('coLegalRepresentativeAdditionalInfo'),
    nameAttr('coReasonsForRepresentation'),
    nameAttr('coGuardianshipIsLegal'),
    nameAttr('coAbleToGiveLegalConsent'),
];

const countryCodeAttributeMap = new Map([
    ['CO', colombiaAttributeNames],
    ['UG', ugandaAttributeNames],
]);

export default class IndividualPage {
    constructor(href) {
        if (href) {
            href.then(h => cy.visit(h));
        } else {
            cy.visit(URL.newIndividual);
        }
    }

    inputAttributes = data => {
        const {
            fullName,
            displayName,
            birthDate,
            email,
            displacementStatus,
            gender,
            identificationDate,
            identificationLocation,
            identificationSource,
            admin2,
            admin3to5,
            relationshipType,
            relatedParty,
        } = data;
        this.typeFirstName(fullName);
        this.typeLastName(displayName);
        this.enterBirthDate(birthDate);
        this.typeEmail(email);
        this.selectDisplacementStatus(displacementStatus);
        this.selectGender(gender);
        this.enterIdentificationDate(identificationDate);
        this.selectIdentificationLocation(identificationLocation);
        this.selectIdentificationSource(identificationSource);
        this.selectAdmin2(admin2);
        this.typeAdmin3to5(admin3to5);
        this.addRelationship({relationshipType, relatedParty});
        return this;
    };
    verifyAttributes = data => {
        const {
            fullName,
            displayName,
            birthDate,
            email,
            displacementStatus,
            gender,
            identificationDate,
            identificationLocation,
            identificationSource,
            admin2,
            admin3to5,
            relationshipType,
            relatedParty,
        } = data;
        this.getFirstName().should('have.value', fullName);
        this.getLastName().should('have.value', displayName);
        this.getBirthDate().should('have.value', birthDate);
        this.getEmail().should('have.value', email);
        this.getDisplacementStatus().should('have.value', displacementStatus);
        this.getGender().should('have.value', gender);
        this.getIdentificationDate().should('have.value', identificationDate);
        this.getIdentificationLocation().should('have.value', identificationLocation);
        this.getIdentificationSource().should('have.value', identificationSource);
        this.getAdmin2().should('have.value', admin2);
        this.getAdmin3().should('have.value', admin3to5);
        this.getAdmin4().should('have.value', admin3to5);
        this.getAdmin5().should('have.value', admin3to5);
        const r = this.getRelationShip();
        r.should('contain.text', relationshipType.toLowerCase());
        r.should('contain.text', relatedParty);
        return this;
    };

    verifyCountrySpecificAttributes = countryCode => {
        const expectedAttributes = countryCodeAttributeMap.get(countryCode);
        expectedAttributes.forEach(attrName => {
            cy.get(attrName).should('exist');
        });
    };

    getFirstName = () => cy.get(selector.fullName);
    getLastName = () => cy.get(selector.displayName);
    getEmail = () => cy.get(selector.email);
    getBirthDate = () => cy.get(selector.birthDate);
    getDisplacementStatus = () => cy.get(selector.displacementStatus);
    getGender = () => cy.get(selector.gender);
    getIdentificationDate = () => cy.get(selector.identificationDate);
    getIdentificationLocation = () => cy.get(selector.identificationLocation);
    getIdentificationSource = () => cy.get(selector.identificationSource);
    getAdmin2 = () => cy.get(selector.admin2);
    getAdmin3 = () => cy.get(selector.admin3);
    getAdmin4 = () => cy.get(selector.admin4);
    getAdmin5 = () => cy.get(selector.admin5);
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
    enterIdentificationDate = value => this.getIdentificationDate().invoke('val', value);
    selectIdentificationLocation = value => this.getIdentificationLocation().select(value);
    selectIdentificationSource = value => this.getIdentificationSource().select(value);
    selectAdmin2 = value => this.getAdmin2().select(value);
    typeAdmin3to5 = value => [this.getAdmin3, this.getAdmin4, this.getAdmin5].forEach(fn => fn().clear().type(value));
    selectRelationshipType = value => this.getRelationshipType().select(value);
    typeRelatedParty = value => this.getRelatedParty().clear().type(value).wait(500);
    addRelatedParty = () => this.getSearchResult().click();
    addRelationship = ({relationshipType, relatedParty}) => {
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
        this.getSituationAnalysis().find(selector.doneSituationAnalysis).check();
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
            cy.get(selector.doneResponse).check();
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
        });
    };

    verifyTaxonomyInput = $t => {
        cy.wrap($t).find(selector.taxonomyBadges).should('have.length', 1);
    };

    save = () => {
        cy.get(selector.saveBtn).click().wait(500);
    };
}