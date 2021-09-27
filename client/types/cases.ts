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