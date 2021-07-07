

export const Urls = {
    CASES_URL: '/cases',
    NEW_CASE_URL: '/cases/new',
    NEW_ATTRIBUTE_URL: '/settings/attributes/new',
    NEW_CASETYPE_URL: '/settings/casetypes/new',
    NEW_INDIVIDUAL_URL: '/individuals/new',
    NEW_PARTYTYPE_URL: '/settings/partytypes/new',
    NEW_RELATIONSHIPTYPE_URL: '/settings/relationshiptypes/new',
};

export const credentials = {
    username: Cypress.env('USERNAME'),
    password: Cypress.env('PASSWORD'),
};
