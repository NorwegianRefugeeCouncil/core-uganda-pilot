export const Urls = {
    CASES_URL: '/cases',
    TEAMS_URL: '/teams',
    NEW_CASE_URL: '/cases/new',
    NEW_ATTRIBUTE_URL: '/settings/attributes/new',
    NEW_CASETYPE_URL: '/settings/casetypes/new',
    NEW_INDIVIDUAL_URL: '/individuals/new',
    NEW_PARTYTYPE_URL: '/settings/partytypes/new',
    NEW_RELATIONSHIPTYPE_URL: '/settings/relationshiptypes/new',
    ATTRIBUTE_URL: '/settings/attributes',
    CASETYPE_URL: '/settings/casetypes',
    INDIVIDUALS_URL: '/individuals',
    PARTYTYPE_URL: '/settings/partytypes',
    RELATIONSHIPTYPE_URL: '/settings/relationshiptypes',
};

export const credentials = {
    username: Cypress.env('USERNAME'),
    password: Cypress.env('PASSWORD'),
};

export const caseTypeTemplate = `{
  "formElements": [
    {
      "type": "dropdown",
      "attributes": {
        "label": "Legal satus",
        "id": "legalStatus",
        "description": "What is the beneficiary's current legal status?",
        "placeholder": "",
        "value": null,
        "multiple": false,
        "options": [
          "Citizen",
          "Permanent resident",
          "Accepted refugee",
          "Asylum seeker",
          "Undetermined"
        ],
        "checkboxOptions": null
      },
      "validation": {
        "required": true
      }
    },
    {
      "type": "checkbox",
      "attributes": {
        "label": "Qualified services",
        "id": "qualifiedServices",
        "description": "What services does the beneficiary qualify for?",
        "placeholder": "",
        "value": null,
        "multiple": false,
        "options": null,
        "checkboxOptions": [
          {
            "label": "Counselling",
            "required": false
          },
          {
            "label": "Representation",
            "required": false
          },
          {
            "label": "Arbitration",
            "required": false
          }
        ]
      },
      "validation": {
        "required": true
      }
    },
    {
      "type": "textarea",
      "attributes": {
        "label": "Notes",
        "id": "notes",
        "description": "Additional information, observations, concerns, etc.",
        "placeholder": "Type here",
        "value": null,
        "multiple": false,
        "options": null,
        "checkboxOptions": null
      },
      "validation": {
        "required": false
      }
    },
    {
      "type": "textinput",
      "attributes": {
        "label": "Project number",
        "id": "projectNumber",
        "description": "Enter the beneficiaries project number, if any",
        "placeholder": "",
        "value": null,
        "multiple": false,
        "options": null,
        "checkboxOptions": null
      },
      "validation": {
        "required": false
      }
    }
  ]
}`;
