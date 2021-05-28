import { newPath, pathFrom } from './index';

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
  });
  describe('getValue', () => {
    it('should return the value at the given path', function() {
      expect(pathFrom('prop1.prop2').getValue({ 'prop1': { 'prop2': 'value' } })).toEqual('value');
    });
    it('should return the value at the given path using indexers', function() {
      expect(pathFrom('prop1[0].prop2').getValue({ 'prop1': [{ 'prop2': 'value' }] })).toEqual('value');
    });
    it('should return undefined if the indexed path does not exist', function() {
      expect(pathFrom('prop1[1].prop2').getValue({ 'prop1': [{ 'prop2': 'value' }] })).toBeUndefined();
    });
    it('should return undefined if the path does not exist', function() {
      expect(pathFrom('snapper.snippers').getValue({ 'prop1': [{ 'prop2': 'value' }] })).toBeUndefined();
    });
  });
});
