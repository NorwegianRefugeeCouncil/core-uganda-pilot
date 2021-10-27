/* Do not change, this code is generated from Golang structs */


export interface Error {
    type: string;
    field: string;
    badValue: any;
    detail: string;
}
export interface ControlValidation {
    required: boolean;
}
export interface CheckboxOption {
    label: LocaleString[];
    value: string;
    required: boolean;
}
export interface LocaleString {
    locale: string;
    value: string;
}
export interface Control {
    name: string;
    type: string;
    placeholder: LocaleString[];
    multiple: boolean;
    readonly: boolean;
    label: LocaleString[];
    description: LocaleString[];
    defaultValue: string[];
    value: string[];
    options: LocaleString[][];
    checkboxOptions: CheckboxOption[];
    validation: ControlValidation;
    errors?: Error[];
}
export interface PartyAttributeDefinition {
    id: string;
    alias: string;
    countryId: string;
    partyTypeIds: string[];
    isPii: boolean;
    formControl: Control;
}
export interface PartyAttributeDefinitionList {
    items: PartyAttributeDefinition[];
}
export interface PartyAttributeDefinitionListOptions {
    partyTypeId: string[];
    countryId: string[];
}
export interface AttributeMap {

}
export interface Party {
    id: string;
    partyTypeIds: string[];
    attributes: {[key: string]: string[]};
}
export interface PartyList {
    items: Party[];
}
export interface PartyType {
    id: string;
    name: string;
    isBuiltIn: boolean;
}
export interface PartyTypeList {

}
export interface Relationship {
    id: string;
    relationshipTypeId: string;
    firstParty: string;
    secondParty: string;
}
export interface RelationshipList {
    items: Relationship[];
}
export interface PartyTypeRule {
    firstPartyTypeId: string;
    secondPartyTypeId: string;
}
export interface RelationshipTypeRule {
}

export interface RelationshipType {
    id: string;
    isDirectional: boolean;
    name: string;
    firstPartyRole: string;
    secondPartyRole: string;
    rules: RelationshipTypeRule[];
}
export interface RelationshipTypeList {
    items: RelationshipType[];
}
export interface Individual {
    id: string;
    partyTypeIds: string[];
    attributes: {[key: string]: string[]};
}
export interface Links {
    first: string;
    previous: string;
    self: string;
    next: string;
    last: string;
}
export interface Pagination {
    page: number;
    perPage: number;
    pageCount: number;
    totalCount: number;
    sort: string;
    links: Links;
}
export interface IndividualList {
    items: Individual[];
    metadata: Pagination;
}
export interface Team {
    id: string;
    name: string;
}
export interface TeamList {
    items: Team[];
}
export interface Country {
    id: string;
    name: string;
}
export interface CountryList {
    items: Country[];
}
export interface Staff {
    id: string;
    individualId: string;
}
export interface StaffList {
    items: Staff[];
}
export interface Membership {
    id: string;
    teamId: string;
    individualId: string;
}
export interface MembershipList {
    items: Membership[];
}
export interface Nationality {
    id: string;
    countryId: string;
    teamId: string;
}
export interface NationalityList {
    items: Nationality[];
}
export interface PartyListOptions {
    partyTypeId: string;
    searchParam: string;
}
export interface PartySearchOptions {
    partyIds: string[];
    partyTypeIds: string[];
    attributes: {[key: string]: string};
    searchParam: string;
}
export interface PartyTypeListOptions {

}
export interface RelationshipListOptions {
    relationshipTypeId: string;
    firstPartyId: string;
    secondPartyId: string;
    eitherPartyId: string;
}
export interface RelationshipTypeListOptions {

}
export interface TeamListOptions {

}
export interface CountryListOptions {

}
export interface StaffListOptions {
    individualId: string;
    organizationId: string;
}
export interface MembershipListOptions {
    individualId: string;
    teamId: string;
}
export interface NationalityListOptions {
    teamId: string;
    countryId: string;
}
export interface IndividualListOptions {
    partyTypeId: string[];
    attributes: {[key: string]: string};
    searchParam: string;
    page: number;
    perPage: number;
    sort: string;
}
export interface IdentificationDocument {
    id: string;
    partyId: string;
    documentNumber: string;
    identificationDocumentTypeId: string;
}
export interface IdentificationDocumentList {
    items: IdentificationDocument[];
}
export interface IdentificationDocumentListOptions {
    partyId: string[];
}
export interface IdentificationDocumentType {
    id: string;
    name: string;
}
export interface IdentificationDocumentTypeList {
    items: IdentificationDocumentType[];
}
export interface IdentificationDocumentTypeListOptions {

}
export interface Section {
    title: LocaleString[];
    controlNames: string[];
}
export interface Form {
    controls: Control[];
    sections: Section[];
    errors?: Error[];
}
export interface Case {
    id: string;
    caseTypeId: string;
    partyId: string;
    teamId: string;
    creatorId: string;
    parentId: string;
    intakeCase: boolean;
    urgency: string;
    form: Form;
    formData: {[key: string]: string[]};
    done: boolean;
}
export interface CaseList {
    items: Case[];
}
export interface CaseType {
    id: string;
    name: string;
    partyTypeId: string;
    teamId: string;
    form: Form;
    intakeCaseType: boolean;
}
export interface CaseTypeList {
    items: CaseType[];
}
export interface Time {

}
export interface Comment {
    id: string;
    caseId: string;
    authorId: string;
    body: string;
    createdAt: Time;
    updatedAt: Time;
}
export interface CommentList {
    items: Comment[];
}
export interface CaseListOptions {
    partyId: string[];
    teamId: string[];
    caseTypeId: string[];
    parentId: string;
    done?: boolean;
}
export interface CaseTypeListOptions {
    partyTypeId: string[];
    teamId: string[];
}
export interface CommentListOptions {
    caseId: string;
}
