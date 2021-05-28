export interface TypeMeta {
  apiVersion: string
  kind: string
}

export interface ObjectMeta {
  name?: string
  labels?: { [key: string]: string }
  annotations?: { [key: string]: string }
  uid?: string
  creationTimestamp?: string
}

export interface ListMeta {
  resourceVersion: string
  continue: string
  remainingItemCount: number
}

export enum StatusType {
  Failure = 'Failure',
  Success = 'Success'
}

export enum StatusReason {
  Unknown = '',
  Unauthorized = 'Unauthorized',
  Forbidden = 'Forbidden',
  NotFound = 'NotFound',
  AlreadyExists = 'AlreadyExists',
  Conflict = 'Conflict',
  Gone = 'Gone',
  Invalid = 'Invalid',
  ServerTimeout = 'ServerTimeout',
  Timeout = 'Timeout',
  TooManyRequests = 'TooManyRequests',
  BadRequest = 'BadRequest',
  MethodNotAllowed = 'MethodNotAllowed',
  NotAcceptable = 'NotAcceptable',
  RequestEntityTooLarge = 'RequestEntityTooLarge',
  UnsupportedMediaType = 'UnsupportedMediaType',
  InternalError = 'InternalError',
  Expired = 'Expired',
  ServiceUnavailable = 'ServiceUnavailable'
}

export enum CauseType {
  FieldValueNotFound = 'FieldValueNotFound',
  FieldValueRequired = 'ValueRequired',
  FieldValueDuplicate = 'FieldValueDuplicate',
  FieldValueInvalid = 'FieldValueInvalid',
  FieldValueNotSupported = 'FieldValueNotSupported',
  UnexpectedServerResponse = 'UnexpectedServerResponse',
  ResourceVersionTooLarge = 'ResourceVersionTooLarge'
}

export interface Status extends TypeMeta {
  metadata?: ListMeta
  status?: StatusType
  message?: string
  reason?: StatusReason
  details?: StatusDetails
  code?: number
}

export interface StatusCause {
  type?: CauseType
  message?: string
  field?: string
}

export interface StatusDetails {
  name?: string
  group?: string
  kind?: string
  uid?: string
  causes?: StatusCause[]
  retryAfterSeconds?: number
}
