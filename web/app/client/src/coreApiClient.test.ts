import {expect} from 'chai'
import {of, OperatorFunction, Subject} from 'rxjs';
import {switchMap} from 'rxjs/operators'
import {IAMClient, CMSClient, prepareRequestOptions, Request, Response} from './coreApiClient';
import {Control, Form, RelationshipTypeRule, Time} from "./types/models";

const defaults = {
    global: {
        scheme: 'testscheme',
        host: 'testhost',
        headers: {
            'X-Authenticated-User-Subject': ['test@user.email']
        }
    },
    cases: {
        case: {
            id: 'TESTCASEID',
            caseTypeId: 'TESTCASETYPEID',
            partyId: 'TESTPARTYID',
            teamId: 'TESTTEAMID',
            creatorId: 'TESTCREATORID',
            parentId: 'TESTPARENTID',
            intakeCase: false,
            form: {} as Form,
            formData: {} as FormData,
            done: false
        },
        caseListOptions: null,
        caseGetRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/cases/TESTCASEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        caseListRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/cases',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        caseUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/cases/TESTCASEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCASEID',
                caseTypeId: 'TESTCASETYPEID',
                partyId: 'TESTPARTYID',
                teamId: 'TESTTEAMID',
                creatorId: 'TESTCREATORID',
                parentId: 'TESTPARENTID',
                intakeCase: false,
                form: {} as Form,
                formData: {} as FormData,
                done: false
            }
        },
        caseCreateRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/cases',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCASEID',
                caseTypeId: 'TESTCASETYPEID',
                partyId: 'TESTPARTYID',
                teamId: 'TESTTEAMID',
                creatorId: 'TESTCREATORID',
                parentId: 'TESTPARENTID',
                intakeCase: false,
                form: {} as Form,
                formData: {} as FormData,
                done: false
            }
        }
    },
    casetypes: {
        casetype: {
            id: 'TESTCASETYPEID',
            name: 'TESTCASETYPENAME',
            partyTypeId: 'TESTPARTYTYPEID',
            teamId: 'TESTTEAMID',
            form: {} as Form,
            intakeCaseType: false
        },
        casetypeListOptions: null,
        casetypeGetRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/casetypes/TESTCASETYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        casetypeListRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/casetypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        casetypeUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/casetypes/TESTCASETYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCASETYPEID',
                name: 'TESTCASETYPENAME',
                partyTypeId: 'TESTPARTYTYPEID',
                teamId: 'TESTTEAMID',
                form: {} as Form,
                intakeCaseType: false
            }
        },
        casetypeCreateRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/casetypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCASETYPEID',
                name: 'TESTCASETYPENAME',
                partyTypeId: 'TESTPARTYTYPEID',
                teamId: 'TESTTEAMID',
                form: {} as Form,
                intakeCaseType: false
            }
        }
    },
    comments: {
        comment: {
            id: 'TESTCOMMENTID',
            caseId: 'TESTCASEID',
            authorId: 'TESTAUTHORID',
            body: 'TESTBODY',
            createdAt: {} as Time,
            updatedAt: {} as Time,
        },
        commentListOptions: null,
        commentGetRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/comments/TESTCOMMENTID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        commentDeleteRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/comments/TESTCOMMENTID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'DELETE',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        commentListRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/comments',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        commentUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/comments/TESTCOMMENTID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCOMMENTID',
                caseId: 'TESTCASEID',
                authorId: 'TESTAUTHORID',
                body: 'TESTBODY',
                createdAt: {} as Time,
                updatedAt: {} as Time,
            }
        },
        commentCreateRequestOptions: {
            url: 'testscheme://testhost/apis/cms/v1/comments',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCOMMENTID',
                caseId: 'TESTCASEID',
                authorId: 'TESTAUTHORID',
                body: 'TESTBODY',
                createdAt: {} as Time,
                updatedAt: {} as Time,
            }
        }
    },
    countries: {
        country: {
            id: 'TESTCOUNTRYID',
            name: 'TESTCOUNTRYNAME'
        },
        countryListOptions: null,
        countryGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/countries/TESTCOUNTRYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        countryListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/countries',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        countryUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/countries/TESTCOUNTRYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCOUNTRYID',
                name: 'TESTCOUNTRYNAME'
            }
        },
        countryCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/countries',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTCOUNTRYID',
                name: 'TESTCOUNTRYNAME'
            }
        }
    },
    identificationdocuments: {
        identificationdocument: {
            id: 'TESTIDENTIFICATIONDOCUMENTID',
            partyId: 'TESTPARTYID',
            documentNumber: 'TESTDOCUMENTNUMBER',
            identificationDocumentTypeId: 'TESTIDENTIFICATIONDOCUMENTTYPE',
        },
        identificationdocumentListOptions: null,
        identificationdocumentGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocuments/TESTIDENTIFICATIONDOCUMENTID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        identificationdocumentDeleteRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocuments/TESTIDENTIFICATIONDOCUMENTID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'DELETE',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        identificationdocumentListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocuments',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        identificationdocumentUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocuments/TESTIDENTIFICATIONDOCUMENTID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTIDENTIFICATIONDOCUMENTID',
                partyId: 'TESTPARTYID',
                documentNumber: 'TESTDOCUMENTNUMBER',
                identificationDocumentTypeId: 'TESTIDENTIFICATIONDOCUMENTTYPE',
            }
        },
        identificationdocumentCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocuments',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTIDENTIFICATIONDOCUMENTID',
                partyId: 'TESTPARTYID',
                documentNumber: 'TESTDOCUMENTNUMBER',
                identificationDocumentTypeId: 'TESTIDENTIFICATIONDOCUMENTTYPE',
            }
        }
    },
    identificationdocumenttypes: {
        identificationdocumenttype: {
            id: 'TESTIDENTIFICATIONDOCUMENTTYPEID',
            name: 'TESTIDENTIFICATIONDOCUMENTTYPENAME'
        },
        identificationdocumenttypeListOptions: null,
        identificationdocumenttypeGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocumenttypes/TESTIDENTIFICATIONDOCUMENTTYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        identificationdocumenttypeListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocumenttypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        identificationdocumenttypeUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocumenttypes/TESTIDENTIFICATIONDOCUMENTTYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTIDENTIFICATIONDOCUMENTTYPEID',
                name: 'TESTIDENTIFICATIONDOCUMENTTYPENAME'
            }
        },
        identificationdocumenttypeCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/identificationdocumenttypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTIDENTIFICATIONDOCUMENTTYPEID',
                name: 'TESTIDENTIFICATIONDOCUMENTTYPENAME'
            }
        }
    },
    individuals: {
        individual: {
            id: 'TESTINDIVIDUALID',
            partyTypeIds: [
                'TESTPARTYTYPEID'
            ],
            attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
        },
        individualListOptions: null,
        individualGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/individuals/TESTINDIVIDUALID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        individualListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/individuals',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        individualUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/individuals/TESTINDIVIDUALID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTINDIVIDUALID',
                partyTypeIds: [
                    'TESTPARTYTYPEID'
                ],
                attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
            }
        },
        individualCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/individuals',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTINDIVIDUALID',
                partyTypeIds: [
                    'TESTPARTYTYPEID'
                ],
                attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
            }
        }
    },
    memberships: {
        membership: {
            id: 'TESTMEMBERSHIPID',
            teamId: 'TESTTEAMID',
            individualId: 'TESTINDIVIDUALID'
        },
        membershipListOptions: null,
        membershipGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/memberships/TESTMEMBERSHIPID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        membershipListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/memberships',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        membershipUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/memberships/TESTMEMBERSHIPID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTMEMBERSHIPID',
                teamId: 'TESTTEAMID',
                individualId: 'TESTINDIVIDUALID'
            }
        },
        membershipCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/memberships',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTMEMBERSHIPID',
                teamId: 'TESTTEAMID',
                individualId: 'TESTINDIVIDUALID'
            }
        }
    },
    nationalities: {
        nationality: {
            id: 'TESTNATIONALITYID',
            countryId: 'TESTCOUNTRYID',
            teamId: 'TESTTEAMID'
        },
        nationalityListOptions: null,
        nationalityGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/nationalities/TESTNATIONALITYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        nationalityListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/nationalities',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        nationalityUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/nationalities/TESTNATIONALITYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTNATIONALITYID',
                countryId: 'TESTCOUNTRYID',
                teamId: 'TESTTEAMID'
            }
        },
        nationalityCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/nationalities',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTNATIONALITYID',
                countryId: 'TESTCOUNTRYID',
                teamId: 'TESTTEAMID'
            }
        }
    },
    attributes: {
        partyattributedefinition: {
            id: 'TESTPARTYATTRIBUTEDEFINITIONID',
            countryId: 'TESTCOUNTRYID',
            partyTypeIds: ['TESTPARTYTYPEID'],
            isPii: false,
            formControl: {} as Control
        },
        partyattributedefinitionListOptions: null,
        partyattributedefinitionGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/attributes/TESTPARTYATTRIBUTEDEFINITIONID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partyattributedefinitionListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/attributes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partyattributedefinitionUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/attributes/TESTPARTYATTRIBUTEDEFINITIONID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTPARTYATTRIBUTEDEFINITIONID',
                countryId: 'TESTCOUNTRYID',
                partyTypeIds: ['TESTPARTYTYPEID'],
                isPii: false,
                formControl: {} as Control
            }
        },
        partyattributedefinitionCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/attributes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTPARTYATTRIBUTEDEFINITIONID',
                countryId: 'TESTCOUNTRYID',
                partyTypeIds: ['TESTPARTYTYPEID'],
                isPii: false,
                formControl: {} as Control
            }
        }
    },
    parties: {
        party: {
            id: 'TESTPARTYID',
            partyTypeIds: [
                'TESTPARTYTYPEID'
            ],
            attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
        },
        partyListOptions: null,
        partyGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/parties/TESTPARTYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partyListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/parties',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partyUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/parties/TESTPARTYID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTPARTYID',
                partyTypeIds: [
                    'TESTPARTYTYPEID'
                ],
                attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
            }
        },
        partyCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/parties',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTPARTYID',
                partyTypeIds: [
                    'TESTPARTYTYPEID'
                ],
                attributes: {'TESTATTRIBUTEKEY': ['TESTATTRIBUTEVALUE']}
            }
        }
    },
    partytypes: {
        partytype: {
            id: 'TESTPARTYTYPEID',
            name: 'TESTPARTYTYPENAME',
            isBuiltIn: false
        },
        partytypeListOptions: null,
        partytypeGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/partytypes/TESTPARTYTYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partytypeListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/partytypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        partytypeUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/partytypes/TESTPARTYTYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTPARTYTYPEID',
                name: 'TESTPARTYTYPENAME',
                isBuiltIn: false
            }
        },
        partytypeCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/partytypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTPARTYTYPEID',
                name: 'TESTPARTYTYPENAME',
                isBuiltIn: false
            }
        }
    },
    relationships: {
        relationship: {
            id: 'TESTRELATIONSHIPID',
            relationshipTypeId: 'TESTRELATIONSHIPTYPEID',
            firstParty: 'TESTFIRSTPARTY',
            secondParty: 'TESTSECONDPARTY'
        },
        relationshipListOptions: null,
        relationshipGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationships/TESTRELATIONSHIPID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        relationshipDeleteRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationships/TESTRELATIONSHIPID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'DELETE',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        relationshipListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationships',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        relationshipUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationships/TESTRELATIONSHIPID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTRELATIONSHIPID',
                relationshipTypeId: 'TESTRELATIONSHIPTYPEID',
                firstParty: 'TESTFIRSTPARTY',
                secondParty: 'TESTSECONDPARTY'
            }
        },
        relationshipCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationships',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTRELATIONSHIPID',
                relationshipTypeId: 'TESTRELATIONSHIPTYPEID',
                firstParty: 'TESTFIRSTPARTY',
                secondParty: 'TESTSECONDPARTY'
            }
        }
    },
    relationshiptypes: {
        relationshiptype: {
            id: 'TESTRELATIONSHIPTYPEID',
            isDirectional: false,
            name: 'TESTRELATIONSHIPTYPENAME',
            firstPartyRole: 'TESTFIRSTPARTYROLE',
            secondPartyRole: 'TESTSECONDPARTYROLE',
            rules: [] as RelationshipTypeRule[]
        },
        relationshiptypeListOptions: null,
        relationshiptypeGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationshiptypes/TESTRELATIONSHIPTYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        relationshiptypeListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationshiptypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        relationshiptypeUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationshiptypes/TESTRELATIONSHIPTYPEID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTRELATIONSHIPTYPEID',
                isDirectional: false,
                name: 'TESTRELATIONSHIPTYPENAME',
                firstPartyRole: 'TESTFIRSTPARTYROLE',
                secondPartyRole: 'TESTSECONDPARTYROLE',
                rules: [] as RelationshipTypeRule[]
            }
        },
        relationshiptypeCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/relationshiptypes',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTRELATIONSHIPTYPEID',
                isDirectional: false,
                name: 'TESTRELATIONSHIPTYPENAME',
                firstPartyRole: 'TESTFIRSTPARTYROLE',
                secondPartyRole: 'TESTSECONDPARTYROLE',
                rules: [] as RelationshipTypeRule[]
            }
        }
    },
    teams: {
        team: {
            id: 'TESTTEAMID',
            name: 'TESTTEAMNAME'
        },
        teamListOptions: null,
        teamGetRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/teams/TESTTEAMID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        teamListRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/teams',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: undefined
        },
        teamUpdateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/teams/TESTTEAMID',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTTEAMID',
                name: 'TESTTEAMNAME'
            }
        },
        teamCreateRequestOptions: {
            url: 'testscheme://testhost/apis/iam/v1/teams',
            headers: {
                'X-Authenticated-User-Subject': ['test@user.email'],
                'Content-Type': ['application/json'],
                Accept: ['application/json']
            },
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            body: {
                id: 'TESTTEAMID',
                name: 'TESTTEAMNAME'
            }
        }
    },
}

describe('Unit Tests: Case Management System Client Set, Request Creation', () => {
    describe('While using the Cases Client', () => {
        let subject
        let cmsClient
        let caseGetClient
        let caseUpdateClient
        let caseCreateClient
        let caseListClient
        let caseId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            caseId = defaults.cases.case.id
            subject = new Subject()
            cmsClient = new CMSClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = cmsClient.Cases()
            client.execute = execute
            return client
        }

        describe('When making a request to get a case', () => {
            before(() => {
                caseGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseGetRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(caseGetClient.Get()).subscribe((kase) => {
                    done()
                })
                subject.next(caseId)
            })
        })

        describe('When making a request to list cases', () => {
            before(() => {
                caseListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseListRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(caseListClient.List()).subscribe((kase) => {
                    done()
                })
                subject.next(defaults.cases.caseListOptions)
            })
        })

        describe('When making a request to update a case', () => {
            before(() => {
                caseUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseUpdateRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(caseUpdateClient.Update()).subscribe((kase) => {
                    done()
                })
                subject.next(defaults.cases.case)
            })
        })

        describe('When making a request to create a case', () => {
            before(() => {
                caseCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseCreateRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(caseCreateClient.Create()).subscribe((kase) => {
                    done()
                })
                subject.next(defaults.cases.case)
            })
        })
    })

    describe('While using the Case Types Client', () => {
        let subject
        let cmsClient
        let caseGetClient
        let caseUpdateClient
        let caseCreateClient
        let caseListClient
        let caseId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            caseId = defaults.cases.case.id
            subject = new Subject()
            cmsClient = new CMSClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = cmsClient.Cases()
            client.execute = execute
            return client
        }

        describe('When making a request to get a case', () => {
            before(() => {
                caseGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseGetRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(caseGetClient.Get()).subscribe((kase) => {
                    done()
                })
                subject.next(caseId)
            })
        })

        describe('When making a request to list cases', () => {
            before(() => {
                caseListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseListRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(caseListClient.List()).subscribe((kase) => {
                    done()
                })
                subject.next(defaults.cases.caseListOptions)
            })
        })

        describe('When making a request to update a case', () => {
            before(() => {
                caseUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseUpdateRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(caseUpdateClient.Update()).subscribe((kase) => {
                    done()
                })
                subject.next(defaults.cases.case)
            })
        })

        describe('When making a request to create a case', () => {
            before(() => {
                caseCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.cases.caseCreateRequestOptions)
                                return of(new Response(defaults.cases.case))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(caseCreateClient.Create()).subscribe((kase) => {
                    done()
                })
                subject.next(defaults.cases.case)
            })
        })
    })

    describe('While using the Comments Client', () => {
        let subject
        let cmsClient
        let commentGetClient
        let commentUpdateClient
        let commentCreateClient
        let commentListClient
        let commentDeleteClient
        let commentId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            commentId = defaults.comments.comment.id
            subject = new Subject()
            cmsClient = new CMSClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = cmsClient.Comments()
            client.execute = execute
            return client
        }

        describe('When making a request to get a comment', () => {
            before(() => {
                commentGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.comments.commentGetRequestOptions)
                                return of(new Response(defaults.comments.comment))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(commentGetClient.Get()).subscribe((comment) => {
                    done()
                })
                subject.next(commentId)
            })
        })

        describe('When making a request to list comments', () => {
            before(() => {
                commentListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.comments.commentListRequestOptions)
                                return of(new Response(defaults.comments.comment))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(commentListClient.List()).subscribe((comment) => {
                    done()
                })
                subject.next(defaults.comments.commentListOptions)
            })
        })

        describe('When making a request to update a comment', () => {
            before(() => {
                commentUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.comments.commentUpdateRequestOptions)
                                return of(new Response(defaults.comments.comment))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(commentUpdateClient.Update()).subscribe((comment) => {
                    done()
                })
                subject.next(defaults.comments.comment)
            })
        })

        describe('When making a request to create a comment', () => {
            before(() => {
                commentCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.comments.commentCreateRequestOptions)
                                return of(new Response(defaults.comments.comment))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(commentCreateClient.Create()).subscribe((comment) => {
                    done()
                })
                subject.next(defaults.comments.comment)
            })
        })

        describe('When making a request to delete a comment', () => {
            before(() => {
                commentDeleteClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.comments.commentDeleteRequestOptions)
                                return of(new Response({}))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a delete request', (done) => {
                subject.pipe(commentDeleteClient.Delete()).subscribe((comment) => {
                    done()
                })
                subject.next(commentId)
            })
        })
    })
})

describe('Unit Tests: Identity & Access Management Client Set, Request Creation', () => {
    describe('While using the Countries Client', () => {
        let subject
        let iamClient
        let countryGetClient
        let countryUpdateClient
        let countryCreateClient
        let countryListClient
        let countryId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            countryId = defaults.countries.country.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Countries()
            client.execute = execute
            return client
        }

        describe('When making a request to get a country', () => {
            before(() => {
                countryGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.countries.countryGetRequestOptions)
                                return of(new Response(defaults.countries.country))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(countryGetClient.Get()).subscribe((country) => {
                    done()
                })
                subject.next(countryId)
            })
        })

        describe('When making a request to list countries', () => {
            before(() => {
                countryListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.countries.countryListRequestOptions)
                                return of(new Response(defaults.countries.country))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(countryListClient.List()).subscribe((country) => {
                    done()
                })
                subject.next(defaults.countries.countryListOptions)
            })
        })

        describe('When making a request to update a country', () => {
            before(() => {
                countryUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.countries.countryUpdateRequestOptions)
                                return of(new Response(defaults.countries.country))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(countryUpdateClient.Update()).subscribe((country) => {
                    done()
                })
                subject.next(defaults.countries.country)
            })
        })

        describe('When making a request to create a country', () => {
            before(() => {
                countryCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.countries.countryCreateRequestOptions)
                                return of(new Response(defaults.countries.country))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(countryCreateClient.Create()).subscribe((country) => {
                    done()
                })
                subject.next(defaults.countries.country)
            })
        })
    })

    describe('While using the Identification Documents Client', () => {
        let subject
        let iamClient
        let identificationDocumentGetClient
        let identificationDocumentUpdateClient
        let identificationDocumentCreateClient
        let identificationDocumentListClient
        let identificationDocumentDeleteClient
        let identificationDocumentId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            identificationDocumentId = defaults.identificationdocuments.identificationdocument.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.IdentificationDocuments()
            client.execute = execute
            return client
        }

        describe('When making a request to get an identification document', () => {
            before(() => {
                identificationDocumentGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocuments.identificationdocumentGetRequestOptions)
                                return of(new Response(defaults.identificationdocuments.identificationdocument))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(identificationDocumentGetClient.Get()).subscribe((identificationdocument) => {
                    done()
                })
                subject.next(identificationDocumentId)
            })
        })

        describe('When making a request to list identification documents', () => {
            before(() => {
                identificationDocumentListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocuments.identificationdocumentListRequestOptions)
                                return of(new Response(defaults.identificationdocuments.identificationdocument))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(identificationDocumentListClient.List()).subscribe((identificationdocument) => {
                    done()
                })
                subject.next(defaults.identificationdocuments.identificationdocumentListOptions)
            })
        })

        describe('When making a request to update an identification document', () => {
            before(() => {
                identificationDocumentUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocuments.identificationdocumentUpdateRequestOptions)
                                return of(new Response(defaults.identificationdocuments.identificationdocument))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(identificationDocumentUpdateClient.Update()).subscribe((identificationdocument) => {
                    done()
                })
                subject.next(defaults.identificationdocuments.identificationdocument)
            })
        })

        describe('When making a request to create an identification document', () => {
            before(() => {
                identificationDocumentCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocuments.identificationdocumentCreateRequestOptions)
                                return of(new Response(defaults.identificationdocuments.identificationdocument))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(identificationDocumentCreateClient.Create()).subscribe((identificationdocument) => {
                    done()
                })
                subject.next(defaults.identificationdocuments.identificationdocument)
            })
        })

        describe('When making a request to delete an identification document', () => {
            before(() => {
                identificationDocumentDeleteClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocuments.identificationdocumentDeleteRequestOptions)
                                return of(new Response({}))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a delete request', (done) => {
                subject.pipe(identificationDocumentDeleteClient.Delete()).subscribe((identificationdocument) => {
                    done()
                })
                subject.next(identificationDocumentId)
            })
        })
    })

    describe('While using the Identification Document Types Client', () => {
        let subject
        let iamClient
        let identificationdocumenttypeGetClient
        let identificationdocumenttypeUpdateClient
        let identificationdocumenttypeCreateClient
        let identificationdocumenttypeListClient
        let identificationdocumenttypeId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            identificationdocumenttypeId = defaults.identificationdocumenttypes.identificationdocumenttype.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.IdentificationDocumentTypes()
            client.execute = execute
            return client
        }

        describe('When making a request to get an identification document type', () => {
            before(() => {
                identificationdocumenttypeGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocumenttypes.identificationdocumenttypeGetRequestOptions)
                                return of(new Response(defaults.identificationdocumenttypes.identificationdocumenttype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(identificationdocumenttypeGetClient.Get()).subscribe((identificationdocumenttype) => {
                    done()
                })
                subject.next(identificationdocumenttypeId)
            })
        })

        describe('When making a request to list identification document types', () => {
            before(() => {
                identificationdocumenttypeListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocumenttypes.identificationdocumenttypeListRequestOptions)
                                return of(new Response(defaults.identificationdocumenttypes.identificationdocumenttype))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(identificationdocumenttypeListClient.List()).subscribe((identificationdocumenttype) => {
                    done()
                })
                subject.next(defaults.identificationdocumenttypes.identificationdocumenttypeListOptions)
            })
        })

        describe('When making a request to update an identification document type', () => {
            before(() => {
                identificationdocumenttypeUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocumenttypes.identificationdocumenttypeUpdateRequestOptions)
                                return of(new Response(defaults.identificationdocumenttypes.identificationdocumenttype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(identificationdocumenttypeUpdateClient.Update()).subscribe((identificationdocumenttype) => {
                    done()
                })
                subject.next(defaults.identificationdocumenttypes.identificationdocumenttype)
            })
        })

        describe('When making a request to create a country', () => {
            before(() => {
                identificationdocumenttypeCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.identificationdocumenttypes.identificationdocumenttypeCreateRequestOptions)
                                return of(new Response(defaults.identificationdocumenttypes.identificationdocumenttype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(identificationdocumenttypeCreateClient.Create()).subscribe((identificationdocumenttype) => {
                    done()
                })
                subject.next(defaults.identificationdocumenttypes.identificationdocumenttype)
            })
        })
    })

    describe('While using the Individuals Client', () => {
        let subject
        let iamClient
        let individualGetClient
        let individualUpdateClient
        let individualCreateClient
        let individualListClient
        let individualId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            individualId = defaults.individuals.individual.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Individuals()
            client.execute = execute
            return client
        }

        describe('When making a request to get an individual', () => {
            before(() => {
                individualGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.individuals.individualGetRequestOptions)
                                return of(new Response(defaults.individuals.individual))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(individualGetClient.Get()).subscribe((individual) => {
                    done()
                })
                subject.next(individualId)
            })
        })

        describe('When making a request to list individuals', () => {
            before(() => {
                individualListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.individuals.individualListRequestOptions)
                                return of(new Response(defaults.individuals.individual))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(individualListClient.List()).subscribe((individual) => {
                    done()
                })
                subject.next(defaults.individuals.individualListOptions)
            })
        })

        describe('When making a request to update a party', () => {
            before(() => {
                individualUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.individuals.individualUpdateRequestOptions)
                                return of(new Response(defaults.individuals.individual))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(individualUpdateClient.Update()).subscribe((individual) => {
                    done()
                })
                subject.next(defaults.individuals.individual)
            })
        })

        describe('When making a request to create a party', () => {
            before(() => {
                individualCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.individuals.individualCreateRequestOptions)
                                return of(new Response(defaults.individuals.individual))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(individualCreateClient.Create()).subscribe((individual) => {
                    done()
                })
                subject.next(defaults.individuals.individual)
            })
        })
    })

    describe('While using the Memberships Client', () => {
        let subject
        let iamClient
        let membershipGetClient
        let membershipUpdateClient
        let membershipCreateClient
        let membershipListClient
        let membershipId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            membershipId = defaults.memberships.membership.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Memberships()
            client.execute = execute
            return client
        }

        describe('When making a request to get a membership', () => {
            before(() => {
                membershipGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.memberships.membershipGetRequestOptions)
                                return of(new Response(defaults.memberships.membership))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(membershipGetClient.Get()).subscribe((membership) => {
                    done()
                })
                subject.next(membershipId)
            })
        })

        describe('When making a request to list memberships', () => {
            before(() => {
                membershipListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.memberships.membershipListRequestOptions)
                                return of(new Response(defaults.memberships.membership))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(membershipListClient.List()).subscribe((membership) => {
                    done()
                })
                subject.next(defaults.memberships.membershipListOptions)
            })
        })

        describe('When making a request to update a membership', () => {
            before(() => {
                membershipUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.memberships.membershipUpdateRequestOptions)
                                return of(new Response(defaults.memberships.membership))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(membershipUpdateClient.Update()).subscribe((membership) => {
                    done()
                })
                subject.next(defaults.memberships.membership)
            })
        })

        describe('When making a request to create a membership', () => {
            before(() => {
                membershipCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.memberships.membershipCreateRequestOptions)
                                return of(new Response(defaults.memberships.membership))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(membershipCreateClient.Create()).subscribe((membership) => {
                    done()
                })
                subject.next(defaults.memberships.membership)
            })
        })
    })

    describe('While using the Nationalities Client', () => {
        let subject
        let iamClient
        let nationalityGetClient
        let nationalityUpdateClient
        let nationalityCreateClient
        let nationalityListClient
        let nationalityId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            nationalityId = defaults.nationalities.nationality.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Nationalities()
            client.execute = execute
            return client
        }

        describe('When making a request to get a nationality', () => {
            before(() => {
                nationalityGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.nationalities.nationalityGetRequestOptions)
                                return of(new Response(defaults.nationalities.nationality))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(nationalityGetClient.Get()).subscribe((nationality) => {
                    done()
                })
                subject.next(nationalityId)
            })
        })

        describe('When making a request to list nationalities', () => {
            before(() => {
                nationalityListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.nationalities.nationalityListRequestOptions)
                                return of(new Response(defaults.nationalities.nationality))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(nationalityListClient.List()).subscribe((nationality) => {
                    done()
                })
                subject.next(defaults.nationalities.nationalityListOptions)
            })
        })

        describe('When making a request to update a nationality', () => {
            before(() => {
                nationalityUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.nationalities.nationalityUpdateRequestOptions)
                                return of(new Response(defaults.nationalities.nationality))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(nationalityUpdateClient.Update()).subscribe((nationality) => {
                    done()
                })
                subject.next(defaults.nationalities.nationality)
            })
        })

        describe('When making a request to create a membership', () => {
            before(() => {
                nationalityCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.nationalities.nationalityCreateRequestOptions)
                                return of(new Response(defaults.nationalities.nationality))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(nationalityCreateClient.Create()).subscribe((nationality) => {
                    done()
                })
                subject.next(defaults.nationalities.nationality)
            })
        })
    })

    describe('While using the Party Attribute Definitions Client', () => {
        let subject
        let iamClient
        let partyattributedefinitionGetClient
        let partyattributedefinitionUpdateClient
        let partyattributedefinitionCreateClient
        let partyattributedefinitionListClient
        let partyattributedefinitionId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            partyattributedefinitionId = defaults.attributes.partyattributedefinition.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.PartyAttributeDefinitions()
            client.execute = execute
            return client
        }

        describe('When making a request to get a party attribute definition', () => {
            before(() => {
                partyattributedefinitionGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.attributes.partyattributedefinitionGetRequestOptions)
                                return of(new Response(defaults.attributes.partyattributedefinition))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(partyattributedefinitionGetClient.Get()).subscribe((partyattributedefinition) => {
                    done()
                })
                subject.next(partyattributedefinitionId)
            })
        })

        describe('When making a request to list party attribute definitions', () => {
            before(() => {
                partyattributedefinitionListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.attributes.partyattributedefinitionListRequestOptions)
                                return of(new Response(defaults.attributes.partyattributedefinition))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(partyattributedefinitionListClient.List()).subscribe((partyattributedefinition) => {
                    done()
                })
                subject.next(defaults.attributes.partyattributedefinitionListOptions)
            })
        })

        describe('When making a request to update a party attribute definition', () => {
            before(() => {
                partyattributedefinitionUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.attributes.partyattributedefinitionUpdateRequestOptions)
                                return of(new Response(defaults.attributes.partyattributedefinition))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(partyattributedefinitionUpdateClient.Update()).subscribe((partyattributedefinition) => {
                    done()
                })
                subject.next(defaults.attributes.partyattributedefinition)
            })
        })

        describe('When making a request to create a party attribute definition', () => {
            before(() => {
                partyattributedefinitionCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.attributes.partyattributedefinitionCreateRequestOptions)
                                return of(new Response(defaults.attributes.partyattributedefinition))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(partyattributedefinitionCreateClient.Create()).subscribe((partyattributedefinition) => {
                    done()
                })
                subject.next(defaults.attributes.partyattributedefinition)
            })
        })
    })

    describe('While using the Parties Client', () => {
        let subject
        let iamClient
        let partyGetClient
        let partyUpdateClient
        let partyCreateClient
        let partyListClient
        let partyId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            partyId = defaults.parties.party.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Parties()
            client.execute = execute
            return client
        }

        describe('When making a request to get a party', () => {
            before(() => {
                partyGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.parties.partyGetRequestOptions)
                                return of(new Response(defaults.parties.party))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(partyGetClient.Get()).subscribe((party) => {
                    done()
                })
                subject.next(partyId)
            })
        })

        describe('When making a request to list parties', () => {
            before(() => {
                partyListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.parties.partyListRequestOptions)
                                return of(new Response(defaults.parties.party))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(partyListClient.List()).subscribe((party) => {
                    done()
                })
                subject.next(defaults.parties.partyListOptions)
            })
        })

        describe('When making a request to update a party', () => {
            before(() => {
                partyUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.parties.partyUpdateRequestOptions)
                                return of(new Response(defaults.parties.party))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(partyUpdateClient.Update()).subscribe((party) => {
                    done()
                })
                subject.next(defaults.parties.party)
            })
        })

        describe('When making a request to create a party', () => {
            before(() => {
                partyCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.parties.partyCreateRequestOptions)
                                return of(new Response(defaults.parties.party))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(partyCreateClient.Create()).subscribe((party) => {
                    done()
                })
                subject.next(defaults.parties.party)
            })
        })
    })

    describe('While using the Party Types Client', () => {
        let subject
        let iamClient
        let partytypeGetClient
        let partytypeUpdateClient
        let partytypeCreateClient
        let partytypeListClient
        let partytypeId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            partytypeId = defaults.partytypes.partytype.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.PartyTypes()
            client.execute = execute
            return client
        }

        describe('When making a request to get a party type', () => {
            before(() => {
                partytypeGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.partytypes.partytypeGetRequestOptions)
                                return of(new Response(defaults.partytypes.partytype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(partytypeGetClient.Get()).subscribe((partytype) => {
                    done()
                })
                subject.next(partytypeId)
            })
        })

        describe('When making a request to list party types', () => {
            before(() => {
                partytypeListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.partytypes.partytypeListRequestOptions)
                                return of(new Response(defaults.partytypes.partytype))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(partytypeListClient.List()).subscribe((partytype) => {
                    done()
                })
                subject.next(defaults.partytypes.partytypeListOptions)
            })
        })

        describe('When making a request to update a party type', () => {
            before(() => {
                partytypeUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.partytypes.partytypeUpdateRequestOptions)
                                return of(new Response(defaults.partytypes.partytype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(partytypeUpdateClient.Update()).subscribe((partytype) => {
                    done()
                })
                subject.next(defaults.partytypes.partytype)
            })
        })

        describe('When making a request to create a party', () => {
            before(() => {
                partytypeCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.partytypes.partytypeCreateRequestOptions)
                                return of(new Response(defaults.partytypes.partytype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(partytypeCreateClient.Create()).subscribe((partytype) => {
                    done()
                })
                subject.next(defaults.partytypes.partytype)
            })
        })
    })

    describe('While using the Relationships Client', () => {
        let subject
        let iamClient
        let relationshipGetClient
        let relationshipUpdateClient
        let relationshipCreateClient
        let relationshipListClient
        let relationshipDeleteClient
        let relationshipId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            relationshipId = defaults.relationships.relationship.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Relationships()
            client.execute = execute
            return client
        }

        describe('When making a request to get a relationship', () => {
            before(() => {
                relationshipGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationships.relationshipGetRequestOptions)
                                return of(new Response(defaults.relationships.relationship))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(relationshipGetClient.Get()).subscribe((relationship) => {
                    done()
                })
                subject.next(relationshipId)
            })
        })

        describe('When making a request to list relationships', () => {
            before(() => {
                relationshipListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationships.relationshipListRequestOptions)
                                return of(new Response(defaults.relationships.relationship))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(relationshipListClient.List()).subscribe((relationship) => {
                    done()
                })
                subject.next(defaults.relationships.relationshipListOptions)
            })
        })

        describe('When making a request to update a relationship', () => {
            before(() => {
                relationshipUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationships.relationshipUpdateRequestOptions)
                                return of(new Response(defaults.relationships.relationship))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(relationshipUpdateClient.Update()).subscribe((relationship) => {
                    done()
                })
                subject.next(defaults.relationships.relationship)
            })
        })

        describe('When making a request to create a relationship', () => {
            before(() => {
                relationshipCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationships.relationshipCreateRequestOptions)
                                return of(new Response(defaults.relationships.relationship))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(relationshipCreateClient.Create()).subscribe((relationship) => {
                    done()
                })
                subject.next(defaults.relationships.relationship)
            })
        })

        describe('When making a request to delete a relationship', () => {
            before(() => {
                relationshipDeleteClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationships.relationshipDeleteRequestOptions)
                                return of(new Response({}))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a delete request', (done) => {
                subject.pipe(relationshipDeleteClient.Delete()).subscribe((relationship) => {
                    done()
                })
                subject.next(relationshipId)
            })
        })
    })

    describe('While using the Relationship Types Client', () => {
        let subject
        let iamClient
        let relationshiptypeGetClient
        let relationshiptypeUpdateClient
        let relationshiptypeCreateClient
        let relationshiptypeListClient
        let relationshiptypeId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            relationshiptypeId = defaults.relationshiptypes.relationshiptype.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.RelationshipTypes()
            client.execute = execute
            return client
        }

        describe('When making a request to get a relationship type', () => {
            before(() => {
                relationshiptypeGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationshiptypes.relationshiptypeGetRequestOptions)
                                return of(new Response(defaults.relationshiptypes.relationshiptype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(relationshiptypeGetClient.Get()).subscribe((relationshiptype) => {
                    done()
                })
                subject.next(relationshiptypeId)
            })
        })

        describe('When making a request to list relationship types', () => {
            before(() => {
                relationshiptypeListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationshiptypes.relationshiptypeListRequestOptions)
                                return of(new Response(defaults.relationshiptypes.relationshiptype))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(relationshiptypeListClient.List()).subscribe((relationshiptype) => {
                    done()
                })
                subject.next(defaults.relationshiptypes.relationshiptypeListOptions)
            })
        })

        describe('When making a request to update a relationship types', () => {
            before(() => {
                relationshiptypeUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationshiptypes.relationshiptypeUpdateRequestOptions)
                                return of(new Response(defaults.relationshiptypes.relationshiptype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(relationshiptypeUpdateClient.Update()).subscribe((relationshiptype) => {
                    done()
                })
                subject.next(defaults.relationshiptypes.relationshiptype)
            })
        })

        describe('When making a request to create a relationship type', () => {
            before(() => {
                relationshiptypeCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.relationshiptypes.relationshiptypeCreateRequestOptions)
                                return of(new Response(defaults.relationshiptypes.relationshiptype))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(relationshiptypeCreateClient.Create()).subscribe((relationshiptype) => {
                    done()
                })
                subject.next(defaults.relationshiptypes.relationshiptype)
            })
        })
    })

    describe('While using the Teams Client', () => {
        let subject
        let iamClient
        let teamGetClient
        let teamUpdateClient
        let teamCreateClient
        let teamListClient
        let teamId

        // sets up test variables and injects mock execute function
        const init = (execute: () => OperatorFunction<Request, Response>) => {
            teamId = defaults.teams.team.id
            subject = new Subject()
            iamClient = new IAMClient(
                defaults.global.scheme,
                defaults.global.host,
                defaults.global.headers
            )
            let client = iamClient.Teams()
            client.execute = execute
            return client
        }

        describe('When making a request to get a team', () => {
            before(() => {
                teamGetClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.teams.teamGetRequestOptions)
                                return of(new Response(defaults.teams.team))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a get request', (done) => {
                subject.pipe(teamGetClient.Get()).subscribe((team) => {
                    done()
                })
                subject.next(teamId)
            })
        })

        describe('When making a request to list teams', () => {
            before(() => {
                teamListClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.teams.teamListRequestOptions)
                                return of(new Response(defaults.teams.team))
                            })
                        )
                    }
                })
            })
            it('it should correctly a build list request', (done) => {
                subject.pipe(teamListClient.List()).subscribe((team) => {
                    done()
                })
                subject.next(defaults.teams.teamListOptions)
            })
        })

        describe('When making a request to update a team', () => {
            before(() => {
                teamUpdateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.teams.teamUpdateRequestOptions)
                                return of(new Response(defaults.teams.team))
                            })
                        )
                    }
                })
            })
            it('it should correctly build an update request', (done) => {
                subject.pipe(teamUpdateClient.Update()).subscribe((team) => {
                    done()
                })
                subject.next(defaults.teams.team)
            })
        })

        describe('When making a request to create a team', () => {
            before(() => {
                teamCreateClient = init((): OperatorFunction<Request, Response> => {
                    return source => {
                        return source.pipe(
                            switchMap(req => {
                                let ro = prepareRequestOptions(req)
                                expect(ro).to.deep.equal(defaults.teams.teamCreateRequestOptions)
                                return of(new Response(defaults.teams.team))
                            })
                        )
                    }
                })
            })
            it('it should correctly build a create request', (done) => {
                subject.pipe(teamCreateClient.Create()).subscribe((team) => {
                    done()
                })
                subject.next(defaults.teams.team)
            })
        })
    })
})
