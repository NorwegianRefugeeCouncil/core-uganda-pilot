/* Do not change, this code is generated from Golang structs */


export class Error {
    type: string;
    field: string;
    badValue: any;
    detail: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.type = source["type"];
        this.field = source["field"];
        this.badValue = source["badValue"];
        this.detail = source["detail"];
    }
}
export class FormElementValidation {
    required: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.required = source["required"];
    }
}
export class CheckboxOption {
    label: string;
    required: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.label = source["label"];
        this.required = source["required"];
    }
}
export class FormElementAttributes {
    label: string;
    name: string;
    value: string[];
    description: string;
    placeholder: string;
    multiple: boolean;
    options: string[];
    checkboxOptions: CheckboxOption[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.label = source["label"];
        this.name = source["name"];
        this.value = source["value"];
        this.description = source["description"];
        this.placeholder = source["placeholder"];
        this.multiple = source["multiple"];
        this.options = source["options"];
        this.checkboxOptions = this.convertValues(source["checkboxOptions"], CheckboxOption);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class AttributeTranslation {
    locale: string;
    longFormulation: string;
    shortFormulation: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.locale = source["locale"];
        this.longFormulation = source["longFormulation"];
        this.shortFormulation = source["shortFormulation"];
    }
}
export class Attribute {
    id: string;
    name: string;
    countryId: string;
    partyTypeIds: string[];
    isPii: boolean;
    translations: AttributeTranslation[];
    type: string;
    attributes: FormElementAttributes;
    validation: FormElementValidation;
    errors?: Error[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.name = source["name"];
        this.countryId = source["countryId"];
        this.partyTypeIds = source["partyTypeIds"];
        this.isPii = source["isPii"];
        this.translations = this.convertValues(source["translations"], AttributeTranslation);
        this.type = source["type"];
        this.attributes = this.convertValues(source["attributes"], FormElementAttributes);
        this.validation = this.convertValues(source["validation"], FormElementValidation);
        this.type = source["type"];
        this.attributes = this.convertValues(source["attributes"], FormElementAttributes);
        this.validation = this.convertValues(source["validation"], FormElementValidation);
        this.errors = this.convertValues(source["errors"], Error);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}

export class AttributeList {
    items: Attribute[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Attribute);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class PartyAttributes {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class Party {
    id: string;
    partyTypeIds: string[];
    attributes: {[key: string]: string[]};

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.partyTypeIds = source["partyTypeIds"];
        this.attributes = source["attributes"];
    }
}
export class PartyList {
    items: Party[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Party);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class PartyType {
    id: string;
    name: string;
    isBuiltIn: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.name = source["name"];
        this.isBuiltIn = source["isBuiltIn"];
    }
}
export class PartyTypeList {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class Relationship {
    id: string;
    relationshipTypeId: string;
    firstParty: string;
    secondParty: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.relationshipTypeId = source["relationshipTypeId"];
        this.firstParty = source["firstParty"];
        this.secondParty = source["secondParty"];
    }
}
export class RelationshipList {
    items: Relationship[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Relationship);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class PartyTypeRule {
    firstPartyTypeId: string;
    secondPartyTypeId: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.firstPartyTypeId = source["firstPartyTypeId"];
        this.secondPartyTypeId = source["secondPartyTypeId"];
    }
}
/*export class RelationshipTypeRule {
    ?: PartyTypeRule;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this. = this.convertValues(source[""], PartyTypeRule);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}*/

export class RelationshipType {
    id: string;
    isDirectional: boolean;
    name: string;
    firstPartyRole: string;
    secondPartyRole: string;
    //rules: RelationshipTypeRule[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.isDirectional = source["isDirectional"];
        this.name = source["name"];
        this.firstPartyRole = source["firstPartyRole"];
        this.secondPartyRole = source["secondPartyRole"];
        //this.rules = this.convertValues(source["rules"], RelationshipTypeRule);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class RelationshipTypeList {
    items: RelationshipType[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], RelationshipType);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Individual {
    id: string;
    partyTypeIds: string[];
    attributes: {[key: string]: string[]};

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.partyTypeIds = source["partyTypeIds"];
        this.attributes = source["attributes"];
    }
}
export class Links {
    first: string;
    previous: string;
    self: string;
    next: string;
    last: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.first = source["first"];
        this.previous = source["previous"];
        this.self = source["self"];
        this.next = source["next"];
        this.last = source["last"];
    }
}
export class Pagination {
    page: number;
    perPage: number;
    pageCount: number;
    totalCount: number;
    sort: string;
    links: Links;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.page = source["page"];
        this.perPage = source["perPage"];
        this.pageCount = source["pageCount"];
        this.totalCount = source["totalCount"];
        this.sort = source["sort"];
        this.links = this.convertValues(source["links"], Links);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class IndividualList {
    items: Individual[];
    metadata: Pagination;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Individual);
        this.metadata = this.convertValues(source["metadata"], Pagination);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Team {
    id: string;
    name: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.name = source["name"];
    }
}
export class TeamList {
    items: Team[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Team);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Country {
    id: string;
    name: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.name = source["name"];
    }
}
export class CountryList {
    items: Country[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Country);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Staff {
    id: string;
    individualId: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.individualId = source["individualId"];
    }
}
export class StaffList {
    items: Staff[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Staff);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Membership {
    id: string;
    teamId: string;
    individualId: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.teamId = source["teamId"];
        this.individualId = source["individualId"];
    }
}
export class MembershipList {
    items: Membership[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Membership);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Nationality {
    id: string;
    CountryId: string;
    teamId: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.CountryId = source["CountryId"];
        this.teamId = source["teamId"];
    }
}
export class NationalityList {
    items: Nationality[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Nationality);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class PartyListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class PartySearchOptions {
    partyIds: string[];
    partyTypeIds: string[];
    attributes: {[key: string]: string};
    searchParam: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.partyIds = source["partyIds"];
        this.partyTypeIds = source["partyTypeIds"];
        this.attributes = source["attributes"];
        this.searchParam = source["searchParam"];
    }
}
export class PartyTypeListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class RelationshipListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class RelationshipTypeListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class AttributeListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class TeamListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class CountryListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class StaffListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class MembershipListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class NationalityListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class IndividualListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class FormElement {
    type: string;
    attributes: FormElementAttributes;
    validation: FormElementValidation;
    errors?: Error[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.type = source["type"];
        this.attributes = this.convertValues(source["attributes"], FormElementAttributes);
        this.validation = this.convertValues(source["validation"], FormElementValidation);
        this.errors = this.convertValues(source["errors"], Error);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CaseTemplate {
    formElements: FormElement[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.formElements = this.convertValues(source["formElements"], FormElement);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Case {
    id: string;
    caseTypeId: string;
    partyId: string;
    done: boolean;
    parentId: string;
    teamId: string;
    creatorId: string;
    template?: CaseTemplate;
    intakeCase: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.caseTypeId = source["caseTypeId"];
        this.partyId = source["partyId"];
        this.done = source["done"];
        this.parentId = source["parentId"];
        this.teamId = source["teamId"];
        this.creatorId = source["creatorId"];
        this.template = this.convertValues(source["template"], CaseTemplate);
        this.intakeCase = source["intakeCase"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CaseList {
    items: Case[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Case);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CaseType {
    id: string;
    name: string;
    partyTypeId: string;
    teamId: string;
    template?: CaseTemplate;
    intakeCaseType: boolean;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.name = source["name"];
        this.partyTypeId = source["partyTypeId"];
        this.teamId = source["teamId"];
        this.template = this.convertValues(source["template"], CaseTemplate);
        this.intakeCaseType = source["intakeCaseType"];
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CaseTypeList {
    items: CaseType[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], CaseType);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class Time {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class Comment {
    id: string;
    caseId: string;
    authorId: string;
    body: string;
    createdAt: Time;
    updatedAt: Time;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.id = source["id"];
        this.caseId = source["caseId"];
        this.authorId = source["authorId"];
        this.body = source["body"];
        this.createdAt = this.convertValues(source["createdAt"], Time);
        this.updatedAt = this.convertValues(source["updatedAt"], Time);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}
export class CommentList {
    items: Comment[];

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.items = this.convertValues(source["items"], Comment);
    }

	convertValues(a: any, classs: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice) {
	        return (a as any[]).map(elem => this.convertValues(elem, classs));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = new classs(a[key]);
	            }
	            return a;
	        }
	        return new classs(a);
	    }
	    return a;
	}
}

export class CaseListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class CaseTypeListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}
export class CommentListOptions {


    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);

    }
}