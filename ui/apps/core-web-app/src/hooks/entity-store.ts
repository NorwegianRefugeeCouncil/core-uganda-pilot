import { HasObjectMeta, Status, TypeMeta } from '@core/api-client';
import { useCallback } from 'react';
import Api from '../data/api';

export type Resource = HasObjectMeta & TypeMeta

export type ResourceReference = {
  apiVersion: string
  kind: string
  name: string
}

export type insertOneOptions = {}

export type getOneOptions = {}

export type getByTypeOptions = {
  apiVersion: string,
  kind: string
}

export type updateOneOptions = {}

export type updateWithCallbackOptions = {}

export type deleteOneOptions = {}

export type insertManyOptions = {}

export type setAllOptions = {}

export type setAllByTypeOptions = {
  apiVersion: string,
  kind: string
}

type storeInterface<T extends Resource> = {
  insertOne: (resource: T, options?: insertOneOptions) => T
  insertMany: (resources: T[], options?: insertManyOptions) => T[]
  getOne: (ref: ResourceReference, options?: getOneOptions) => T | undefined
  getByType: (options: getByTypeOptions) => T[]
  updateOne: (resource: T, options?: updateOneOptions) => T
  updateWithCallback: (resourceRef: ResourceReference, cb: (old: T) => T, options?: updateWithCallbackOptions) => T
  deleteOne: (ref: ResourceReference, options?: deleteOneOptions) => T | undefined
  setAll: (resources: T[], options?: setAllOptions) => T[]
  setAllByType: (resources: T[], options: setAllByTypeOptions) => T[]
  syncDiscovery: () => void
}

const getKey = (resource: Resource): string => {
  const apiVersionParts = resource.apiVersion.split('/');
  return apiVersionParts[0] + '/' + resource.kind + '/' + resource.metadata.name;
};

const getKeyFromRef = (resourceRef: ResourceReference): string => {
  const apiVersionParts = resourceRef.apiVersion.split('/');
  return apiVersionParts[0] + '/' + resourceRef.kind + '/' + resourceRef.name;
};

type EnvelopeProps = {
  isNew: boolean
  status: Status
}

type Envelope<T> = EnvelopeProps & {
  resource: T
}

const withEnvelope = <T extends Resource>(resource: T, options?: Partial<EnvelopeProps>): Envelope<T> => {
  const props: EnvelopeProps = {
    isNew: true,
    status: {
      apiVersion: 'v1',
      kind: 'Status',
      status: 'Success',
      code: 200
    },
    ...options
  };
  return {
    resource,
    ...props
  };
};

const withoutEnvelope = <T>(payload: any): T => {
  const e = payload as Envelope<T>;
  return e ? e.resource : null;
};

export const useEntityStore = <T extends Resource>() => {

  const insertOne = useCallback<(resource: T) => T>(resource => {
    const key = getKey(resource);
    localStorage.setItem(key, JSON.stringify(withEnvelope(resource)));
    return resource;
  }, []);

  const insertMany = useCallback<(resources: T[]) => T[]>(resources => {
    for (let resource of resources) {
      const key = getKey(resource);
      if (localStorage.getItem(key)) {
        throw new Error(`entity with key ${key} already exists`);
      }
    }
    const result: T[] = [];
    for (let resource of resources) {
      const key = getKey(resource);
      localStorage.setItem(key, JSON.stringify(withEnvelope(resource)));
      result.push(resource);
    }
    return result;
  }, []);

  const getOne = useCallback<(ref: ResourceReference) => T | undefined>((ref) => {
    const key = getKeyFromRef(ref);
    return withoutEnvelope(JSON.parse(localStorage.getItem(key)));
  }, []);

  const getByType = useCallback<(type: { apiVersion: string, kind: string }) => T[]>(type => {
    const apiVersionParts = type.apiVersion.split('/');
    const prefix = apiVersionParts[0] + '/' + type.kind + '/';
    const result: T[] = [];
    for (let localStorageKey in localStorage) {
      if (localStorageKey.startsWith(prefix)) {
        result.push(withoutEnvelope(JSON.parse(localStorage.getItem(localStorageKey))));
      }
    }
    return result;
  }, []);

  const updateOne = useCallback<(res: T) => T>(res => {
    const key = getKey(res);
    if (!localStorage.getItem(key)) {
      throw new Error(`entity with key ${key} not found in store`);
    }
    localStorage.setItem(key, JSON.stringify(withEnvelope(res)));
    return res;
  }, []);

  const setAll = useCallback<(resources: T[]) => T[]>(resources => {
    for (let localStorageKey in localStorage) {
      localStorage.removeItem(localStorageKey);
    }
    const ret: T[] = [];
    for (let resource of resources) {
      const key = getKey(resource);
      ret.push(resource);
      localStorage.setItem(key, JSON.stringify(withEnvelope(resource)));
    }
    return ret;
  }, []);

  const setAllByType = useCallback<(resources: T[], options: setAllByTypeOptions) => T[]>((resources, options) => {
    const apiVersionParts = options.apiVersion.split('/');
    const prefix = apiVersionParts[0] + '/' + options.kind + '/';
    for (let localStorageKey in localStorage) {
      if (localStorageKey.startsWith(prefix)) {
        localStorage.removeItem(localStorageKey);
      }
    }
    const ret: T[] = [];
    for (let resource of resources) {
      const key = getKey(resource);
      if (!key.startsWith(prefix)) {
        console.warn(`trying to set all values for type ${options.apiVersion}/${options.kind} but got item with type ${resource.apiVersion}/${resource.kind}`);
      }
      if (!resource.apiVersion) {
        resource.apiVersion = options.apiVersion;
      }
      if (!resource.kind) {
        resource.kind = options.kind;
      }
      ret.push(resource);
      localStorage.setItem(key, JSON.stringify(withEnvelope(resource)));
    }
    return ret;
  }, []);

  const deleteOne = useCallback<(ref: ResourceReference) => T>(ref => {
    const key = getKeyFromRef(ref);
    const ret = withoutEnvelope<T>(JSON.parse(localStorage.getItem(key)));
    localStorage.removeItem(key);
    return ret;
  }, []);

  const updateWithCallback = useCallback<(ref: ResourceReference, cb: (old: T) => T) => T>((ref, cb) => {
    const key = getKeyFromRef(ref);
    const res = withoutEnvelope<T>(JSON.parse(localStorage.getItem(key)));
    const updated = cb(res);
    localStorage.setItem(key, JSON.stringify(withEnvelope(updated)));
    return updated;
  }, []);

  const syncDiscovery = useCallback(() => {
    Api.discovery().v1().apiServices().list().toPromise().then(value => {
      setAllByType((value.items as any) as T[], { apiVersion: value.apiVersion, kind: 'APIService' });
    });
  }, []);

  const syncFromServer = useCallback((args: { group: string, resource: string }) => {

  }, []);

  const ret: storeInterface<T> = {
    insertOne,
    insertMany,
    getOne,
    getByType,
    updateOne,
    updateWithCallback,
    deleteOne,
    setAll,
    setAllByType,
    syncDiscovery
  };

  return ret;
};
