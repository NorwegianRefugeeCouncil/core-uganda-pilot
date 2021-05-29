import { ListMeta, ObjectMeta, TypeMeta } from '../../meta/v1';

export interface APIService extends TypeMeta {
  metadata: ObjectMeta
  spec: APIServiceSpec
}

export interface APIServiceList extends TypeMeta {
  metadata: ListMeta
  items: APIService[]
}

export interface APIServiceSpec {
  group: string
  version: string
}

export interface APIGroup extends TypeMeta {
  name: string
  versions: GroupVersionForDiscovery[]
  preferredVersion: GroupVersionForDiscovery
}

export interface APIGroupList extends TypeMeta {
  groups: APIGroup[]
}

export interface GroupVersionForDiscovery {
  groupVersion: string
  version: string
}

export interface APIResource {
  name: string
  singularName: string
  namespaced: boolean
  group: string
  version: string
  kind: string
  verbs: string[]
}

export interface APIResourceList extends TypeMeta {
  groupVersion: string
  apiResources: APIResource[]
}
