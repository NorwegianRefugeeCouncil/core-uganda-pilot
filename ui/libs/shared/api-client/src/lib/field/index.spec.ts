import { newPath, pathFrom } from './index';
import exp = require('constants');

describe('path', () => {
  describe('pathFrom', () => {
    it('should parse a simple path', function() {
      const actual = pathFrom('spec');
      const expected = newPath('spec');
      expect(actual).toEqual(expected);
    });

    it('should parse a path with a sub property', function() {
      const actual = pathFrom('spec.bla');
      const expected = newPath('spec').child('bla');
      expect(actual).toEqual(expected);
    });

    it('should parse a path with an indexer', function() {
      const actual = pathFrom('spec[0].bla');
      const expected = newPath('spec').index(0).child('bla');
      expect(actual).toEqual(expected);
    });

    it('should parse a path ending with an indexer', function() {
      const actual = pathFrom('spec.bla[0]');
      const expected = newPath('spec', 'bla').index(0);
      expect(actual).toEqual(expected);
    });

    it('should parse a path with a string indexer', function() {
      const actual = pathFrom('spec[snip].bla');
      const expected = newPath('spec').key('snip').child('bla');
      expect(actual).toEqual(expected);
    });

    it('should parse a path ending with a string indexer', function() {
      const actual = pathFrom('spec.bla[snip]');
      const expected = newPath('spec', 'bla').key('snip');
      expect(actual).toEqual(expected);
    });

    it('should parse a path ending with multiple indexers', function() {
      const actual = pathFrom('spec[0][0]');
      const expected = newPath('spec').index(0).index(0);
      expect(actual).toEqual(expected);
    });

    it('should parse a path ending with a mix of indexers and non-indexers', function() {
      const actual = pathFrom('spec[0][0][0].a.def[0]');

      const expected = newPath('spec')
        .index(0)
        .index(0)
        .index(0)
        .child('a')
        .child('def')
        .index(0);

      expect(actual).toEqual(expected);
    });
  });
  describe('get', () => {
    it('should return the value at the given path', function() {
      expect(pathFrom('prop1.prop2').get({ 'prop1': { 'prop2': 'value' } })).toEqual('value');
    });
    it('should return the value at the given path using indexers', function() {
      expect(pathFrom('prop1[0].prop2').get({ 'prop1': [{ 'prop2': 'value' }] })).toEqual('value');
    });
    it('should return undefined if the indexed path does not exist', function() {
      expect(pathFrom('prop1[1].prop2').get({ 'prop1': [{ 'prop2': 'value' }] })).toBeUndefined();
    });
    it('should return undefined if the path does not exist', function() {
      expect(pathFrom('snapper.snippers').get({ 'prop1': [{ 'prop2': 'value' }] })).toBeUndefined();
    });
    it('should return the value of a subpath', function() {
      expect(pathFrom('prop1').get({ 'prop1': [{ 'prop2': 'value' }] })).toEqual([{ 'prop2': 'value' }]);
    });
  });
  describe('set', () => {
    it('should set the value at the given path', function() {
      const obj = { prop1: {} } as any;
      pathFrom('prop1.prop2').set(obj, 'abc');
      expect(obj.prop1.prop2).toEqual('abc');
    });
    it('should set the value at the given path if the path does not yet exist', function() {
      const obj = { prop1: {} } as any;
      pathFrom('prop1.prop2.prop3[0].abc').set(obj, 'abc');
      expect(obj.prop1.prop2.prop3[0].abc).toEqual('abc');
    });
    it('should set the value at the indexed path', function() {
      const obj = { prop1: [] } as any;
      pathFrom('prop1[0]').set(obj, 'abc');
      expect(obj).toEqual({ prop1: ['abc'] });
    });
    it('should set the value at the indexed path index', function() {
      const obj = { prop1: ['abc'] } as any;
      pathFrom('prop1[1]').set(obj, 'def');
      expect(obj).toEqual({ prop1: ['abc', 'def'] });
    });
    it('should replace the value at the indexed path index', function() {
      const obj = { prop1: ['abc', 'def'] } as any;
      pathFrom('prop1[0]').set(obj, 'new');
      expect(obj).toEqual({ prop1: ['new', 'def'] });
    });
    it('should set the value at the root', function() {
      const obj = {} as any;
      pathFrom('prop1').set(obj, 'new');
      expect(obj).toEqual({ prop1: 'new' });
    });
  });
  describe('remove', () => {
    it('should remove the value at the given root path', function() {
      const obj = { prop1: {} };
      pathFrom('prop1').remove(obj);
      expect(obj).toEqual({});
    });
    it('should remove the value at the given child path', function() {
      const obj = { prop1: { prop2: 'abc' } };
      pathFrom('prop1.prop2').remove(obj);
      expect(obj).toEqual({ prop1: {} });
    });
    it('should remove the value at the given child indexed path', function() {
      const obj = { prop1: { prop2: ['abc'] } };
      pathFrom('prop1.prop2[0]').remove(obj);
      expect(obj).toEqual({ prop1: { prop2: [] } });
    });
    it('should remove the value at the given root indexed path', function() {
      const obj = ['abc'];
      pathFrom('[0]').remove(obj);
      expect(obj).toEqual([]);
    });
  });
  describe('add', () => {
    it('should add an item to an array', function() {
      const obj = { abc: [] };
      pathFrom('abc').add(obj, 'def');
      expect(obj).toEqual({ abc: ['def'] });
    });
    it('should add an item to a nested array', function() {
      const obj = { abc: { def: [] } };
      pathFrom('abc.def').add(obj, 'new');
      expect(obj).toEqual({ abc: { def: ['new'] } });
    });
    it('should add an item to a consecutive array', function() {
      const obj = {};
      pathFrom('abc[0][0]').add(obj, 'new');
      expect(obj).toEqual({ abc: [[['new']]] });
    });
    it('should add an item to an array even if the array does not exist', function() {
      const obj = {};
      pathFrom('abc').add(obj, 'def');
      expect(obj).toEqual({ abc: ['def'] });
    });
  });
  describe('insert', () => {
    it('should insert an item at the specific index', function() {
      const obj = { abc: ['1', '2'] };
      pathFrom('abc').insert(obj, 0, 'new');
      expect(obj).toEqual({ abc: ['new', '1', '2'] });
    });
  });
  describe('setIndexed', () => {
    it('should set the indexed value', function() {
      const obj = { abc: [{ key: 'a', value: '123' }] };
      pathFrom('abc').setIndexed(obj, 'key', { key: 'a', value: '234' });
      expect(obj).toEqual({ abc: [{ key: 'a', value: '234' }] });
    });
    it('should add a value if it doesnt exist', function() {
      const obj = { abc: [{ key: 'a', value: '123' }] };
      pathFrom('abc').setIndexed(obj, 'key', { key: 'b', value: '234' });
      expect(obj).toEqual({ abc: [{ key: 'a', value: '123' }, { key: 'b', value: '234' }] });
    });
  });
  describe('ensurePath', () => {
    it('should not do anything if the path exists', function() {
      const actual = { prop1: { prop2: 'abc' } };
      const expected = { prop1: { prop2: 'abc' } };
      pathFrom('prop1.prop2.abc').ensurePath(actual);
      expect(actual).toEqual(expected);
    });
    it('should add a property if the path does not exist', function() {
      const actual = { prop1: { prop2: {} } };
      const expected = { prop1: { prop2: { prop3: { prop4: {} } } } };
      pathFrom('prop1.prop2.prop3.prop4.def').ensurePath(actual);
      expect(actual).toEqual(expected);
    });
    it('should add an array if the path does not exist', function() {
      const actual = {};
      const expected = { prop1: [{ prop2: {} }] };
      pathFrom('prop1[0].prop2.abc').ensurePath(actual);
      expect(actual).toEqual(expected);
    });
    it('should add an array if the path partially exist', function() {
      const actual = { prop1: [] };
      const expected = { prop1: [{ prop2: {} }] };
      pathFrom('prop1[0].prop2.abc').ensurePath(actual);
      expect(actual).toEqual(expected);
    });
    it('should add an array as the last property', function() {
      const actual = {};
      const expected = { prop1: { prop2: [] } };
      pathFrom('prop1.prop2[0]').ensurePath(actual);
      expect(actual).toEqual(expected);
    });
    it('should handle multiple consecutive arrays', function() {
      const actual = {};
      const expected = { prop1: [[]] };
      pathFrom('prop1[0][0]').ensurePath(actual);
      expect(actual).toEqual(expected);
    });
  });

});
