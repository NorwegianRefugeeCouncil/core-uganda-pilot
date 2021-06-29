const mockText = 'mock';
const mockUpdatedText = 'update';

describe('Create & edit casetype', function () {
    // this runs before each test and stores the option element we intend to
    // select in a cypress alias ie. as("A") to be called later with cy.get("@A")
    // this then gives us a reference to verify that the correct value is being assigned
    beforeEach(() => {
        cy.visit('settings/casetypes/new');
        cy.get('option[data-cy=partytype-option]')
            // Choosing the second index, because 1st would be selected by
            // default in the form and we wouldn't know if something went wrong
            .eq(1)
            .as('testPartyTypeOption');
        cy.get('option[data-cy=team-option]').eq(1).as('testTeamOption');

        // and a second set to use as update values
        cy.get('option[data-cy=partytype-option]')
            .eq(2)
            .as('testPartyTypeOption2');
        cy.get('option[data-cy=team-option]').eq(2).as('testTeamOption2');
    });

    describe('Create', () => {
        it('registers a new casetype', () => {
            cy.visit('/settings/casetypes/new');
            cy.get('input[name=name]').type(mockText);
            cy.get('@testPartyTypeOption')
                .invoke('prop', 'value')
                .then((value) => {
                    cy.get('select[name=partyTypeId]').select(value);
                });
            cy.get('@testTeamOption')
                .invoke('prop', 'value')
                .then((value) => {
                    cy.get('select[name=teamId]').select(value);
                });

            cy.get('button').contains('Save').click();
        });
        it('saved the casetype', () => {
            cy.visit('/settings/casetypes');
            cy.get('a[data-cy=casetype]').last().click({ force: true });
            // redirected to casetype view page
            cy.get('input[name=name]').should('have.value', mockText);
            cy.get('@testPartyTypeOption')
                .invoke('prop', 'value')
                .then((value) => {
                    cy.get('select[name=partyTypeId]').should(
                        'have.value',
                        value
                    );
                });

            cy.get('@testTeamOption')
                .invoke('prop', 'value')
                .then((value) => {
                    cy.get('select[name=teamId]').should('have.value', value);
                });
        });
    });
    describe('Update', () => {
        it('updates the casetype', () => {
            cy.visit('/settings/casetypes');
            cy.get('a[data-cy=casetype]').last().click({ force: true });

            cy.get('input[name=name]').clear().type(mockUpdatedText);
            cy.get('@testPartyTypeOption2')
                .invoke('prop', 'value')
                .then((value) => {
                    cy.get('select[name=partyTypeId]').select(value);
                });
            // disabled
            // cy.get('@testTeamOption2')
            //     .invoke('prop', 'value')
            //     .then(value => {
            //         cy.get('select[name=teamId]').select(value)
            //     })

            cy.get('button').contains('Save').click();
        });
        it('saved the updated casetype', () => {
            cy.visit('/settings/casetypes');
            cy.get('a[data-cy=casetype]').last().click({ force: true });
            // redirected to casetype view page
            cy.get('input[name=name]').should('have.value', mockUpdatedText);
            cy.get('@testPartyTypeOption2')
                .invoke('prop', 'value')
                .then((value) => {
                    cy.get('select[name=partyTypeId]').should(
                        'have.value',
                        value
                    );
                });
            // disabled
            // cy.get('@testTeamOption2')
            //     .invoke('prop', 'value')
            //     .then(value => {
            //         cy.get('select[name=teamId]').should('have.value', value)
            //     })
        });
    });
});
