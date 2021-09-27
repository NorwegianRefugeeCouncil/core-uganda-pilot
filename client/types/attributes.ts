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
  // type: string;
  // attributes: FormElementAttributes;
  // validation: FormElementValidation;
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