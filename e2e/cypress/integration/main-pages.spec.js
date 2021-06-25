
describe('Individuals', () => {
  it('should show the page', () => {
    cy.visit('/');
  });
});

describe('Individual', () => {
  it('should show the page', () => {
    cy.visit('/individuals')
    cy.get('[data-cy=individual]').first().click()
  });
});

describe('Cases', () => {
  it('should show the page', () => {
    cy.visit('/cases');
  });
});

describe('Teams', () => {
  it('should show the page', () => {
    cy.visit('/teams');
  });
});

describe('Settings', () => {
  it('should show the page', () => {
    cy.visit('/settings');
  });
});

describe('Attribute settings', function() {
  it('should show the page', () => {
    cy.visit('/settings/attributes');
  });
});

describe('Entity Type settings', function() {
  it('should show the page', () => {
    cy.visit('/settings/partytypes');
  });
});


describe('Case Type settings', function() {
  it('should show the page', () => {
    cy.visit('/settings/casetypes');
  });
});


describe('Relationship Type settings', function() {
  it('should show the page', () => {
    cy.visit('/settings/relationshiptypes');
  });
});

describe('Country Type settings', function() {
  it('should show the page', () => {
    cy.visit('/settings/countries');
  });
});

