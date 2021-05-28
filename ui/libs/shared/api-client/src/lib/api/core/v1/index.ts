import { ListMeta, ObjectMeta, TypeMeta } from '../../meta/v1';

export interface FormDefinition extends TypeMeta {
  metadata: ObjectMeta
  spec: FormDefinitionSpec
}

export interface FormDefinitionList extends TypeMeta {
  metadata: ListMeta
  items: FormDefinition[]
}

export interface FormDefinitionSpec {
  group: string
  names: FormDefinitionNames
  versions: FormDefinitionVersion[]
}

export interface FormDefinitionNames {
  plural: string
  singular: string
  kind: string
}

export interface FormDefinitionVersion {
  name: string
  served: boolean
  storage: boolean
  schema: FormDefinitionValidation
}

export interface FormDefinitionValidation {
  formSchema: FormRoot
}

export interface FormRoot {
  root: FormElement
}

export type FieldType = 'text' | 'integer' | 'float' | 'checkbox' | 'radio' | 'select' | 'multiselect' | 'section'

export interface FormElement {
  key?: string
  name?: TranslatedStrings
  description?: TranslatedStrings
  type?: FieldType
  label?: TranslatedStrings
  tooltip?: TranslatedStrings
  help?: TranslatedStrings
  required?: boolean
  min?: number
  max?: number
  minLength?: number
  maxLength?: number
  children?: FormElement[]
}

export interface TranslatedString {
  locale: string
  value: string
}

export type TranslatedStrings = TranslatedString[]

export interface CustomResourceDefinition extends TypeMeta {
  metadata: ObjectMeta
}

export interface CustomResourceDefinitionList extends TypeMeta {
  metadata: ListMeta
  items: CustomResourceDefinition[]
}

export interface CustomResourceDefinitionSpec {
  group: string
  names: CustomResourceDefinitionNames
}

export interface CustomResourceDefinitionNames {
  plural: string
  singular: string
  kind: string
}

export interface CustomResourceDefinitionVersion {
  name: string
  served: boolean
  storage: boolean
  schema: CustomResourceDefinitionValidation
}

export interface CustomResourceDefinitionValidation {
  openAPIV3Schema: JSONSchemaProps
}

export interface JSONSchemaProps {
  id: string
  $ref?: string
  description: string
  type: string
  format: Format
  title: string
  default?: any
  maximum?: number
  exclusiveMaximum: boolean
  minimum?: number
  exclusiveMinimum: boolean
  maxLength?: number
  minLength?: number
  pattern: string
  maxItems?: number
  minItems?: number
  uniqueItems: boolean
  multipleOf?: number
  enum: any[]
  maxProperties?: number
  minProperties?: number
  required: string[]
  items?: JSONSchemaProps | JSONSchemaProps[]
  allOf: JSONSchemaProps[]
  oneOf: JSONSchemaProps[]
  anyOf: JSONSchemaProps[]
  not?: JSONSchemaProps
  properties: { [key: string]: JSONSchemaProps }
  additionalProperties: JSONSchemaProps | boolean
  patternProperties: { [key: string]: JSONSchemaProps }
  dependencies: JSONSchemaProps | string[]
  additionalItems?: JSONSchemaProps | boolean
  definitions: { [key: string]: JSONSchemaProps }
  externalDocs?: ExternalDocumentation
  example?: any
  nullable: boolean
  'x-kubernetes-preserve-unknown-fields'?: boolean
  'x-kubernetes-embedded-resource': boolean
  'x-kubernetes-int-or-string'?: boolean
  'x-kubernetes-list-map': string[]
  'x-kubernetes-list-type'?: ListType
  'x-kubernetes-map-type'?: MapType
}

export enum Format {
  BSONObjectID = 'bsonobjectid',
  URI = 'uri',
  Email = 'email',
  Hostname = 'hostname',
  IPV4 = 'ipv4',
  IPV6 = 'ipv6',
  CIDR = 'cidr',
  MAC = 'mac',
  UUID = 'uuid',
  UUID3 = 'uuid3',
  UUID4 = 'uuid4',
  UUID5 = 'uuid5',
  ISBN = 'isbn',
  ISBN10 = 'isbn10',
  ISBN13 = 'isbn13',
  CreditCard = 'creditcard',
  SSN = 'ssn',
  HexColor = 'hexcolor',
  RGBColor = 'rgbcolor',
  Byte = 'byte',
  Password = 'password',
  Date = 'date',
  Duration = 'duration',
  DateTime = 'datetime'
}

export enum ListType {
  Atomic = 'atomic',
  Set = 'set',
  Map = 'map'
}

export enum MapType {
  Granular = 'granular',
  Atomic = 'atomic'
}

export interface ExternalDocumentation {
  description: string
  url: string
}
