import { renderHook } from '@testing-library/react-hooks';
import { Resource, ResourceReference, setAllByTypeOptions, useEntityStore } from './entity-store';
import * as React from 'react';

type TestResource = Resource & {
  prop: string
}

const aResource = (name: string, apiVersion = 'bla.com/v1', kind = 'Bla'): TestResource => {
  return {
    apiVersion: apiVersion,
    kind: kind,
    metadata: {
      name: name
    },
    prop: 'abc'
  };
};

const getRef = (res: Resource): ResourceReference => {
  return {
    apiVersion: res.apiVersion,
    kind: res.kind,
    name: res.metadata.name
  };
};

describe('entityStore', () => {

  let insertOne: (res: TestResource) => TestResource;
  let insertMany: (res: TestResource[]) => TestResource[];
  let getOne: (ref: ResourceReference) => TestResource;
  let getByType: (type: { apiVersion: string, kind: string }) => TestResource[];
  let updateOne: (res: TestResource) => TestResource;
  let updateWithCallback: (ref: ResourceReference, cb: (old: TestResource) => TestResource) => TestResource;
  let deleteOne: (res: ResourceReference) => TestResource;
  let setAll: (res: TestResource[]) => TestResource[];
  let setAllByType: (resources: TestResource[], options: setAllByTypeOptions) => TestResource[];
  let syncDiscovery: () => void;
  let resource: TestResource;
  let ref: ResourceReference;

  beforeEach(() => {
    for (let localStorageKey in localStorage) {
      localStorage.removeItem(localStorageKey);
    }
    const { result } = renderHook(() => useEntityStore<TestResource>());
    insertOne = result.current.insertOne;
    insertMany = result.current.insertMany;
    getOne = result.current.getOne;
    getByType = result.current.getByType;
    updateOne = result.current.updateOne;
    updateWithCallback = result.current.updateWithCallback;
    deleteOne = result.current.deleteOne;
    setAll = result.current.setAll;
    setAllByType = result.current.setAllByType;
    syncDiscovery = result.current.syncDiscovery;
    resource = aResource('example');
    ref = getRef(resource);
  });

  it('should allow to put and get objects from the store', function() {
    insertOne(resource);
    const obtained = getOne(ref);
    expect(obtained).toEqual(resource);
  });

  it('should allow to get objects by type', function() {
    const a = aResource('a', 'api.com', 'Kind1');
    const b = aResource('b', 'api.com', 'Kind2');
    insertOne(a);
    insertOne(b);
    expect(getByType({ apiVersion: 'api.com', kind: 'Kind1' })).toContainEqual(a);
    expect(getByType({ apiVersion: 'api.com', kind: 'Kind1' })).toHaveLength(1);
    expect(getByType({ apiVersion: 'api.com', kind: 'Kind2' })).toContainEqual(b);
    expect(getByType({ apiVersion: 'api.com', kind: 'Kind2' })).toHaveLength(1);
  });

  it('should allow to insert many', function() {
    insertMany([
      aResource('a'),
      aResource('b'),
      aResource('c')
    ]);
    expect(getOne(getRef(aResource('a')))).toBeTruthy();
    expect(getOne(getRef(aResource('b')))).toBeTruthy();
    expect(getOne(getRef(aResource('c')))).toBeTruthy();
    expect(getOne(getRef(aResource('d')))).toBeNull();
  });

  it('should allow to update objects', function() {
    insertOne(resource);
    const updated = { ...resource, prop: 'new' } as Resource;
    updateOne({ ...resource, prop: 'new' });
    const got = getOne(ref);
    expect(got).toEqual(updated);
  });

  it('should allow to update an object with a callback', function() {
    insertOne(resource);
    updateWithCallback(ref, old => {
      const updated = { ...old };
      updated.prop = 'withCallback';
      return updated;
    });
    expect(getOne(ref).prop).toEqual('withCallback');
  });

  it('should allow to delete objects', function() {
    insertOne(resource);
    expect(getOne(ref)).toEqual(resource);
    deleteOne(ref);
    expect(getOne(ref)).toBeNull();
  });

  it('should not complain if deleting a non-present objet', function() {
    expect(() => {
      deleteOne(ref);
    }).not.toThrow();
  });

  it('should allow to set all values', () => {
    insertOne(resource);
    expect(getOne(ref)).toBeTruthy();
    setAll([
      aResource('a'),
      aResource('b'),
      aResource('c')
    ]);
    expect(getOne(ref)).toBeNull();
    expect(getOne(getRef(aResource('a')))).toBeTruthy();
    expect(getOne(getRef(aResource('b')))).toBeTruthy();
    expect(getOne(getRef(aResource('c')))).toBeTruthy();
    expect(getOne(getRef(aResource('d')))).toBeNull();
  });

  it('should allow to set all values by type', function() {
    insertMany([
      aResource('a', 'api1', 'Kind1'),
      aResource('b', 'api2', 'Kind2'),
      aResource('c', 'api2', 'Kind2')
    ]);
    expect(getOne(getRef(aResource('a', 'api1', 'Kind1')))).toBeTruthy();
    expect(getOne(getRef(aResource('b', 'api2', 'Kind2')))).toBeTruthy();
    expect(getOne(getRef(aResource('c', 'api2', 'Kind2')))).toBeTruthy();

    setAllByType([
      aResource('d', 'api2', 'Kind2'),
      aResource('e', 'api2', 'Kind2')
    ], { apiVersion: 'api2', kind: 'Kind2' });

    expect(getOne(getRef(aResource('a', 'api1', 'Kind1')))).toBeTruthy();
    expect(getOne(getRef(aResource('b', 'api2', 'Kind2')))).toBeFalsy();
    expect(getOne(getRef(aResource('c', 'api2', 'Kind2')))).toBeFalsy();
    expect(getOne(getRef(aResource('d', 'api2', 'Kind2')))).toBeTruthy();
    expect(getOne(getRef(aResource('e', 'api2', 'Kind2')))).toBeTruthy();
  });

  it('should sync discovery', function() {
    syncDiscovery()
  });
});
