export const testId = id => `[data-testid=${id}]`;
export const nameAttr = name => `[name=${name}]`;

export const URL = {
    INDIVIDUALS: '/individuals',
    CASES: '/cases',
    NEW_CASE: '/cases/new',
    TEAMS: '/teams',
    ATTRIBUTES: '/settings/attributes',
    NEW_ATTRIBUTE: '/settings/attributes/new',
    NEW_CASETYPE: '/settings/casetypes/new',
    NEW_INDIVIDUAL: '/individuals/new',
    NEW_PARTYTYPE: '/settings/partytypes/new',
    NEW_RELATIONSHIPTYPE: '/settings/relationshiptypes/new',
    CASETYPES: '/settings/casetypes',
    PARTYTYPES: '/settings/partytypes',
    RELATIONSHIPTYPES: '/settings/relationshiptypes',
};

export const TEST_CASE_TEMPLATE_FIELD = {
    TEXT: testId('test-text'),
    EMAIL: testId('test-email'),
    PHONE: testId('test-tel'),
    URL: testId('test-url'),
    DATE: testId('test-date'),
    TEXTAREA: testId('test-textarea'),
    DROPDOWN: testId('test-dropdown'),
    CHECKBOX: testId('test-checkbox'),
    RADIO: testId('test-radio'),
};
