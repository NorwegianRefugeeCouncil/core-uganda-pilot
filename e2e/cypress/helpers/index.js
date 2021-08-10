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
    DROPDOWN: 'testDropdown',
    CHECKBOX: 'testCheckbox',
    TEXTINPUT: 'testTextInput',
    TEXTAREA: 'testTextarea',
};

export const testId = (id) => `[data-testid=${id}]`;
