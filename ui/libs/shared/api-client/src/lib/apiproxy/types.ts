export type APICallResourceInfo = {
  group: string
  version: string
  resource: string
  name: string
  method: Method
}

export type Method = 'list' | 'get' | 'create' | 'update' | 'delete' | 'unknown'

export type RequestContext = {
  key: string
  request: RequestInit
  url: string
  info: APICallResourceInfo
  payload: any
}

