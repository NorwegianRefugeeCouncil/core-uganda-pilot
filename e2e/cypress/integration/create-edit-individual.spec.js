const mockText = "mock"
const mockUpdatedText = "update"

describe("Create & edit individual",function () {
    describe("Create", () => {
        it("registers a new individual", () => {
            cy.visit('/individuals/new')
            cy.get("input[data-cy=textAttribute]").each(($el) => {
                cy.wrap($el).type(mockText)
            })
            cy.contains("button", "Save").click()
        })
        it("saved the individual", () => {
            cy.visit('/individuals')
            cy.get("a[data-cy=individual]").last().should('contain.text', `${mockText.toUpperCase()}, ${mockText}`)
            cy.get("a[data-cy=individual]").last().click()
            cy.get("input[data-cy=textAttribute]").each($el => {
                cy.wrap($el).should('have.value', mockText)
            })
        })
    })
    describe("Update", () => {
        it ("updates the individual attributes", () => {
            cy.get("input[data-cy=textAttribute]").each($el => {
                cy.wrap($el).clear().type(mockUpdatedText)
            })
            cy.get("button").contains("Save").click({force: true})
        })
        it ("saved the updated individual", () => {
            cy.visit('/individuals')
            cy.get("a[data-cy=individual]").last().should('contain.text', `${mockUpdatedText.toUpperCase()}, ${mockUpdatedText}`)
            cy.get("a[data-cy=individual]").last().click()
            cy.get("input[data-cy=textAttribute]").each($el => {
                cy.wrap($el).should('have.value', mockUpdatedText)
            })
        })
    })
})
