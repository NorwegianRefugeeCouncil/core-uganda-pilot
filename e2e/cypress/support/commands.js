// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
Cypress.Commands.add('login', (username, password) => {
  return
  cy.request({
    method: 'POST',
    url: 'http://localhost:8080/auth/realms/nrc/protocol/openid-connect/token',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    body: `client_id=api&client_secret=e6486272-039d-430f-b3c7-47887aa9e206&grant_type=password&username=${username}&password=${password}&scope=openid`
  }).then(response => {
    cy.intercept(
      { url: '*' },
      req => req.headers['Authorization'] = `Bearer ${response.body.access_token}`
    );
  });
});
