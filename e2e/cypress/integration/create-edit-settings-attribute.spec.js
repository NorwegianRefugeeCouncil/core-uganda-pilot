const mockText = 'mock';
const mockUpdatedText = 'update';

describe('Create & edit attribute', function () {
    describe('Create', () => {
        it('registers a new attribute', () => {
            cy.visit('/settings/attributes/new');
            cy.get('input[data-cy=name]').type(mockText);
            // TODO cy.get('select[data-cy=type]').select('String')
            // TODO cy.get('select[data-cy=subject]').select('Individual')
            cy.get('input[data-cy=personal-info-chkbx').check();

            // add translation
            cy.get('select[data-cy=language]').select('en');
            cy.get('input[data-cy=translation-long').type(mockText);
            cy.get('input[data-cy=translation-short').type(mockText);

            cy.get('button').contains('Save').click();
        });
        it('saved the attribute', () => {
            cy.visit('/settings/attributes');
            cy.get('a[data-cy=attribute]').last().click({ force: true });
            cy.get('input[data-cy=name]').should('have.value', mockText);
            // TODO cy.get('select[data-cy=type]').should(...)
            // TODO cy.get('select[data-cy=subject]').should(...)
            cy.get('input[data-cy=personal-info-chkbx')
                .invoke('prop', 'checked')
                .then((checked) => expect(checked).to.be.true);

            cy.get('input[data-cy=translation-long').should(
                'have.value',
                mockText
            );
            cy.get('input[data-cy=translation-short').should(
                'have.value',
                mockText
            );
        });
    });
    describe('Update', () => {
        it('updates the attribute', () => {
            cy.get('input[data-cy=personal-info-chkbx').uncheck();

            cy.get('input[data-cy=translation-long')
                .clear()
                .type(mockUpdatedText);
            cy.get('input[data-cy=translation-short')
                .clear()
                .type(mockUpdatedText);
            cy.get('button').contains('Save').click();
        });
        it('saved the updated attribute', () => {
            cy.visit('/settings/attributes');
            cy.get('a[data-cy=attribute]').last().click({ force: true });
            cy.get('input[data-cy=personal-info-chkbx')
                .invoke('prop', 'checked')
                .then((checked) => expect(checked).to.be.false);

            cy.get('input[data-cy=translation-long').should(
                'have.value',
                mockUpdatedText
            );
            cy.get('input[data-cy=translation-short').should(
                'have.value',
                mockUpdatedText
            );
        });
    });
});
