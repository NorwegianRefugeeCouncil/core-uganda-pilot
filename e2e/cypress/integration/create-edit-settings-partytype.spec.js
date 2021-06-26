const mockText = "mock"
const mockUpdatedText = "update"

describe("Create & edit partytype", function () {

    describe.only("Create", () => {

        it("registers a new partytype", () => {
            cy.visit('/settings/partytypes/new')
            cy.get('input[name=name]').type(mockText)
            cy.get('input[name=isBuiltIn]').check()
            cy.get('button').contains('Save').click()
            cy.wait(500)
        })
        it("saved the partytype", () => {
            cy.visit('/settings/partytypes')
            cy.get('a[data-cy=partytype]').last().click({force: true})
            cy.get('input[name=name]').should('have.value', mockText)
            cy.get('input[name=isBuiltIn]').invoke('prop', 'checked').then(checked => expect(checked).to.be.true)
        })
    })
    describe("Update", () => {
        it("updates the partytype", () => {

            cy.get('input[data-cy=personal-info-chkbx').uncheck()

            cy.get('input[data-cy=translation-long').clear().type(mockUpdatedText)
            cy.get('input[data-cy=translation-short').clear().type(mockUpdatedText)
            cy.get('button').contains('Save').click()
        })
        it("saved the updated partytype", () => {
            cy.visit('/settings/partytypes')
            cy.get('a[data-cy=partytype]').last().click({force: true})
            cy.get('input[data-cy=personal-info-chkbx').invoke('prop', 'checked').then(checked => expect(checked).to.be.false)

            cy.get('input[data-cy=translation-long').should('have.value', mockUpdatedText)
            cy.get('input[data-cy=translation-short').should('have.value', mockUpdatedText)
        })
    })
})
