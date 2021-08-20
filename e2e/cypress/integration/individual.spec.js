import IndividualPage from '../pages/individualPage';
import IndividualOverviewPage from '../pages/individualOverview.page';

const data = {
    attributes: {
        firstName: 'Test',
        lastName: 'Person',
        birthDate: '1978-10-30',
        email: 'test@email.com',
        status: '0',
        gender: '0',
        relationshipType: 'Is spouse of',
        relatedParty: 'POPPINS',
    },
    attributes_u: {
        firstName: 'Test - updated',
        lastName: 'Person - updated',
        birthDate: '1979-11-12',
        email: 'test.updated@email.com',
        status: '1',
        gender: '1',
        relationshipType: 'Is sibling of',
        relatedParty: 'DOE',
    },
    situationAnalysis: 'text content',
    response: {
        optionIdx: 0,
        priorityTxt: 'text content',
    },
};

function getTestIndividual(lastName) {
    const individualOverviewPage = new IndividualOverviewPage().visitPage();
    return new IndividualPage(individualOverviewPage.searchFor(lastName));
}

describe.skip('Individual Page', function () {
    describe('Navigate', () => {
        it('should navigate to new Individual page from Individuals overview', () => {
            const individualOverviewPage = new IndividualOverviewPage();
            individualOverviewPage.visitPage().newIndividual();
        });
    });
    describe('Attributes', () => {
        it('should create a new Individual', () => {
            const newIndividualPage = new IndividualPage();
            newIndividualPage.visitPage().inputAttributes(data.attributes).save();
        });
        it('should verify that the Individual was properly created', () => {
            const individualPage = getTestIndividual(data.attributes.lastName);
            individualPage.verifyAttributes(data.attributes);
        });
        it('should update attribute name on existing Individual', () => {
            const individualPage = getTestIndividual(data.attributes.lastName);
            individualPage.removeRelationship().inputAttributes(data.attributes_u).save();
        });
        it('should verify that the Individual was properly updated', () => {
            const individualPage = getTestIndividual(data.attributes_u.lastName);
            individualPage.verifyAttributes(data.attributes_u);
        });
    });

    describe('Situation Analysis', () => {
        it('should submit a situation analysis and verify it', () => {
            const individualPage = getTestIndividual(data.attributes_u.lastName);
            individualPage.visitSituationAnalysisTab().fillOutSituationAnalysis(data.situationAnalysis).save();
            individualPage.visitSituationAnalysisTab().verifySituationAnalysis(data.situationAnalysis);
        });
    });

    describe('Response', () => {
        it('should submit a response form and verify it', () => {
            const individualPage = getTestIndividual(data.attributes_u.lastName);
            individualPage.visitResponseTab().fillOutResponse(data.response).save();
            individualPage.visitResponseTab().verifyResponse(data.response);
        });
    });
});
