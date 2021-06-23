describe("The home page", () => {
    // TODO reset and seed the database before running tests
    it("successfully loads", () => {
        cy.visit("/")
    })
    it("redirects to auth", () => {
        cy.url().should('include', "auth/realms/nrc/protocol/openid-connect/auth")
    })
    it("lets a user log in", () => {
        cy.get("input[name=username]").type("admin")
        cy.get("input[name=password]").type("admin{enter}")
    })
})
