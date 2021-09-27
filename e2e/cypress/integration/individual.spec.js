import IndividualPage from '../pages/individualPage';
import IndividualOverviewPage from '../pages/individualOverview.page';
import '../support/commands';

const data = {
    attributes: {
        fullName: 'Test',
        displayName: 'Person',
        birthDate: '1978-10-30',
        email: 'test@email.com',
        displacementStatus: 'Refugee',
        gender: 'Male',
        identificationLocation: 'Kabusu Access Center',
        identificationSource: 'Walk-in Center',
        admin2: 'ABIM',
        admin3to5: 'test',
        relationshipType: 'Is spouse of',
        relatedParty: 'POPPINS',
    },
    attributesUpd: {
        fullName: 'Test - updated',
        displayName: 'Person - updated',
        birthDate: '1979-11-12',
        email: 'test.updated@email.com',
        displacementStatus: 'Internally displaced person',
        gender: 'Female',
        identificationLocation: 'Nsambya Access Center',
        identificationSource: 'FFRM Referral',
        admin2: 'ADJUMANI',
        admin3to5: 'test - updated',
        relationshipType: 'Is sibling of',
        relatedParty: 'DOE',
    },
    situationAnalysis: 'text content',
    response: {
        optionIdx: 0,
        priorityTxt: 'text content',
        commentStartingPoint: 'comment on service starting point',
        commentOtherServices: 'comment on other services',
    },
};

function getTestIndividual(displayName) {
    const individualOverviewPage = new IndividualOverviewPage().visitPage();
    return new IndividualPage(individualOverviewPage.searchFor(displayName));
}

describe('Individual Page', function () {
    beforeEach('Login', () => {
        cy.login('courtney.lare@email.com');
    });
    describe('Navigate', () => {
        it('should navigate to new Individual page from Individuals overview', () => {
            const individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage.visitPage().newIndividual();
        });
    });
    describe('Attributes', () => {
        it('should create a new Individual', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage.inputAttributes(data.attributes).save();
        });
        it('should verify that the Individual was properly created', () => {
            const individualPage = getTestIndividual(data.attributes.displayName);
            individualPage.verifyAttributes(data.attributes);
        });
        it.skip('should update attribute name on existing Individual', () => {
            const individualPage = getTestIndividual(data.attributes.displayName);
            individualPage.removeRelationship().inputAttributes(data.attributesUpd).save();
        });
        it.skip('should verify that the Individual was properly updated', () => {
            const individualPage = getTestIndividual(data.attributesUpd.displayName);
            individualPage.verifyAttributes(data.attributesUpd);
        });
    });

    describe('Situation Analysis', () => {
        it('should submit a situation analysis and verify it', () => {
            const individualPage = getTestIndividual(data.attributes.displayName);
            individualPage.visitSituationAnalysisTab().fillOutSituationAnalysis(data.situationAnalysis).save();
            individualPage.visitSituationAnalysisTab().verifySituationAnalysis(data.situationAnalysis);
        });
    });

    describe('Response', () => {
        it('should submit a response form and verify it', () => {
            const individualPage = getTestIndividual(data.attributes.displayName);
            individualPage.visitResponseTab().fillOutResponse(data.response).save();
            individualPage.visitResponseTab().verifyResponse(data.response);
        });
    });
});
