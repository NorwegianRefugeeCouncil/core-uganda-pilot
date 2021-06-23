let accessToken
describe("Obtain access token", () => {
    it("gets an access token", () => {
        cy.request({
            method: "POST",
            url: "http://localhost:8080/auth/realms/nrc/protocol/openid-connect/token",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: "client_id=api&client_secret=e6486272-039d-430f-b3c7-47887aa9e206&grant_type=password&username=admin&password=admin&scope=openid",
        }).then(response => {
            accessToken = response.body.access_token
            console.log(accessToken)
        })
    })
// cy.server({
//     onAnyRequest: (route, proxy) => {
//         proxy.xhr.setRequestHeader("Authorization", "Bearer " + response.body.access_token)
//     }
// })
})
// TODO reset and seed the database before running tests

beforeEach(() => {
    cy.intercept(
        {url: "*"},
        req => req.headers["Authorization"] = `Bearer ${accessToken}`
    )
})


describe("User log in", () => {
    it('is accessible', () => {
        cy.visit('/')
        cy.get("input[name=username]").type("admin")
        cy.get("input[name=password]").type("admin{enter}")
    })
})
