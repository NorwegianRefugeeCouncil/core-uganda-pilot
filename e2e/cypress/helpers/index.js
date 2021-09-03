export const testId = id => `[data-testid=${id}]`;
export const nameAttr = name => `[name=${name}]`;

export const URL = {
    individuals: '/individuals',
    cases: '/cases',
    newCase: '/cases/new',
    teams: '/teams',
    attributes: '/settings/attributes',
    newAttribute: '/settings/attributes/new',
    newCasetype: '/settings/casetypes/new',
    new_individual: '/individuals/new',
    new_partytype: '/settings/partytypes/new',
    newRelationshipType: '/settings/relationshiptypes/new',
    casetypes: '/settings/casetypes',
    partytypes: '/settings/partytypes',
    relationshiptTypes: '/settings/relationshiptypes',
};

export const TEST_CASE_TEMPLATE_FIELD = {
    text: testId('test-text'),
    EMAIL: testId('test-email'),
    PHONE: testId('test-tel'),
    URL: testId('test-url'),
    DATE: testId('test-date'),
    TEXTAREA: testId('test-textarea'),
    DROPDOWN: testId('test-dropdown'),
    CHECKBOX: testId('test-checkbox'),
    RADIO: testId('test-radio'),
};
