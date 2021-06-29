const mockText = 'mock';
const mockUpdatedText = 'update';

describe('Create & edit relationship types', function () {
    // this runs before each test and stores the option element we intend to
    // select in a cypress alias ie. as("A") to be called later with cy.get("@A")
    // this then gives us a reference to verify that the correct value is being assigned
    beforeEach(() => {
        cy.visit('/settings/relationshiptypes/new');
        cy.get('option[data-cy=first-partytype-option]')
            // Choosing the second index, because 1st would be selected by
            // default in the form and we wouldn't know if something went wrong
            .eq(1)
            .as('testFirstPartyTypeOption');
        cy.get('option[data-cy=second-partytype-option]')
            .eq(2)
            .as('testSecondPartyTypeOption');

        // and a second set to use as update values
        cy.get('option[data-cy=first-partytype-option]')
            .eq(2)
            .as('testFirstPartyTypeOption2');
        cy.get('option[data-cy=second-partytype-option]')
            .eq(3)
            .as('testSecondPartyTypeOption2');
    });

    context('Non-directional relationship type', () => {
        describe('Create', () => {
            it('registers a new relationship type', () => {
                cy.visit('/settings/relationshiptypes/new');

                cy.get('input[name=name]').type(mockText);
                cy.get('input[name=firstPartyRole]').type(mockText);
                cy.get('@testFirstPartyTypeOption')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).select(value);
                    });

                cy.get('button').contains('Save').click();
            });
            it('saved the relationship type', () => {
                cy.visit('/settings/relationshiptypes');
                cy.get('a[data-cy=relationshiptype]')
                    .last()
                    .click({ force: true });
                // redirected to relationshiptype view page

                cy.get('input[name=name]').should('have.value', mockText);
                cy.get('input[name=isDirectional]')
                    .invoke('prop', 'checked')
                    .then((checked) => expect(checked).to.be.false);
                cy.get('input[name=firstPartyRole]').should(
                    'have.value',
                    mockText
                );
                cy.get('@testFirstPartyTypeOption')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).should('have.value', value);
                    });
            });
        });
        describe('Update', () => {
            it('updates the relationship type', () => {
                cy.visit('/settings/relationshiptypes');
                cy.get('a[data-cy=relationshiptype]')
                    .last()
                    .click({ force: true });
                // redirected to relationshiptype view page

                cy.get('input[name=name]').clear().type(mockUpdatedText);
                cy.get('input[name=firstPartyRole]')
                    .clear()
                    .type(mockUpdatedText);
                cy.get('@testFirstPartyTypeOption2')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).select(value);
                    });

                cy.get('button').contains('Save').click();
            });
            it('saved the updated relationship type', () => {
                cy.visit('/settings/relationshiptypes');
                cy.get('a[data-cy=relationshiptype]')
                    .last()
                    .click({ force: true });
                // redirected to relationshiptype view page

                cy.get('input[name=name]').should(
                    'have.value',
                    mockUpdatedText
                );
                cy.get('input[name=isDirectional]')
                    .invoke('prop', 'checked')
                    .then((checked) => expect(checked).to.be.false);
                cy.get('input[name=firstPartyRole]').should(
                    'have.value',
                    mockUpdatedText
                );
                cy.get('@testFirstPartyTypeOption2')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).should('have.value', value);
                    });
            });
        });
    });
    context('Directional relationship type', () => {
        describe('Create', () => {
            it('registers a new relationship type', () => {
                cy.visit('/settings/relationshiptypes/new');

                cy.get('input[name=name]').type(mockText);
                cy.get('input[name=isDirectional]').check();
                cy.get('input[name=firstPartyRole]').type(mockText);
                cy.get('input[name=secondPartyRole]').type(mockText);
                cy.get('@testFirstPartyTypeOption')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).select(value);
                    });
                cy.get('@testSecondPartyTypeOption')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].secondPartyTypeId"]'
                        ).select(value);
                    });

                cy.get('button').contains('Save').click();
            });
            it('saved the relationship type', () => {
                cy.visit('/settings/relationshiptypes');
                cy.get('a[data-cy=relationshiptype]')
                    .last()
                    .click({ force: true });
                // redirected to relationshiptype view page

                cy.get('input[name=name]').should('have.value', mockText);
                cy.get('input[name=isDirectional]')
                    .invoke('prop', 'checked')
                    .then((checked) => expect(checked).to.be.true);
                cy.get('input[name=firstPartyRole]').should(
                    'have.value',
                    mockText
                );
                cy.get('input[name=secondPartyRole]').should(
                    'have.value',
                    mockText
                );
                cy.get('@testFirstPartyTypeOption')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).should('have.value', value);
                    });
                cy.get('@testSecondPartyTypeOption')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].secondPartyTypeId"]'
                        ).should('have.value', value);
                    });
            });
        });
        describe('Update', () => {
            it('updates the relationship type', () => {
                cy.visit('/settings/relationshiptypes');
                cy.get('a[data-cy=relationshiptype]')
                    .last()
                    .click({ force: true });
                // redirected to relationshiptype view page

                cy.get('input[name=name]').clear().type(mockUpdatedText);
                cy.get('input[name=firstPartyRole]')
                    .clear()
                    .type(mockUpdatedText);
                cy.get('input[name=secondPartyRole]')
                    .clear()
                    .type(mockUpdatedText);
                cy.get('@testFirstPartyTypeOption2')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).select(value);
                    });
                cy.get('@testSecondPartyTypeOption2')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].secondPartyTypeId"]'
                        ).select(value);
                    });

                cy.get('button').contains('Save').click();
            });
            it('saved the updated relationship type', () => {
                cy.visit('/settings/relationshiptypes');
                cy.get('a[data-cy=relationshiptype]')
                    .last()
                    .click({ force: true });
                // redirected to relationshiptype view page

                cy.get('input[name=name]').should(
                    'have.value',
                    mockUpdatedText
                );
                cy.get('input[name=isDirectional]')
                    .invoke('prop', 'checked')
                    .then((checked) => expect(checked).to.be.true);
                cy.get('input[name=firstPartyRole]').should(
                    'have.value',
                    mockUpdatedText
                );
                cy.get('input[name=secondPartyRole]').should(
                    'have.value',
                    mockUpdatedText
                );
                cy.get('@testFirstPartyTypeOption2')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].firstPartyTypeId"]'
                        ).should('have.value', value);
                    });
                cy.get('@testSecondPartyTypeOption2')
                    .invoke('prop', 'value')
                    .then((value) => {
                        cy.get(
                            'select[name="rules[0].secondPartyTypeId"]'
                        ).should('have.value', value);
                    });
            });
        });
    });
});
