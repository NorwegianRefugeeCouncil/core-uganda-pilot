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
    ATTRIBUTE: '/settings/attributes',
    CASETYPE: '/settings/casetypes',
    PARTYTYPE: '/settings/partytypes',
    RELATIONSHIPTYPE: '/settings/relationshiptypes',
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

export const testId = (id) => `[data-testid=${id}]`;
