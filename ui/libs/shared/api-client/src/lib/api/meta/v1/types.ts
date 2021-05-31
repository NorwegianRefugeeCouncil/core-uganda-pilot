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

export type StatusType = 'Failure' | 'Success'

export type StatusReason =
  ''
  | 'Unauthorized'
  | 'Forbidden'
  | 'NotFound'
  | 'AlreadyExists'
  | 'Conflict'
  | 'Gone'
  | 'Invalid'
  | 'ServerTimeout'
  | 'Timeout'
  | 'TooManyRequests'
  | 'BadRequest'
  | 'MethodNotAllowed'
  | 'NotAcceptable'
  | 'RequestEntityTooLarge'
  | 'UnsupportedMediaType'
  | 'InternalError'
  | 'Expired'
  | 'ServiceUnavailable';

export type CauseType =
  'FieldValueNotFound'
  | 'ValueRequired'
  | 'FieldValueDuplicate'
  | 'FieldValueInvalid'
  | 'FieldValueNotSupported'
  | 'UnexpectedServerResponse'
  | 'ResourceVersionTooLarge'


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
