const mockText = "mock"
const mockUpdatedText = "update"

describe("Create & edit case", function () {
    beforeEach(() => {
        cy.visit('/cases/new')
        cy.get("option[data-cy=casetype-option]").first()
            .as('caseTypeName');
    })
    describe("Create", () => {

        it("registers a new case", () => {
            cy.visit('/cases/new')
            cy.get("option[data-cy=casetype-option]").first()
                .invoke('attr', 'value')
                .then((value) => cy.get("select[data-cy=casetype-select]").select(value))


            cy.get("option[data-cy=party-option]").first()
                .invoke('attr', 'value')
                .then(value => cy.get("select[data-cy=party-select]").select(value))

            cy.get("textarea[data-cy=description]").type(mockText)
            cy.get("button[type=submit]").click()
        })
        it("saved the case", () => {
            let name
            cy.get('@caseTypeName').invoke('text').then(t => name = t)
            cy.visit('/cases')
            cy.get('tr').last().within($row => {
                cy.wrap($row).should('contain.text', name)
            })
            cy.get('tr').last().click({force: true})
            cy.get("textarea[data-cy=description]").should('contain.text', mockText)
            cy.get("input[data-cy=done-check]").invoke('prop', 'checked').then(checked => expect(checked).to.be.false)
        })
    })
    describe("Update", () => {
        it("updates the case", () => {
            cy.visit('/cases')

            cy.get('tr').last().within($row => {
                cy.wrap($row).should('contain.text', name)
            })
            cy.get('tr').last().click({force: true})
            cy.get("textarea[data-cy=description]").clear().type(mockUpdatedText)
            cy.get("input[data-cy=done-check]").check()
            cy.get("button[type=submit]").click()
        })
        it("saved the updated case", () => {
            cy.visit('/cases')

            cy.get('tr').last().within($row => {
                cy.wrap($row).should('contain.text', name)
            })
            cy.get('tr').last().click({force: true})
            cy.get("textarea[data-cy=description]").should('contain.text', mockUpdatedText)
            cy.get("input[data-cy=done-check]").invoke('prop', 'checked').then(checked => expect(checked).to.be.true)
        })
    })
})
