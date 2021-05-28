import { NewPath, PathFrom } from './index';

describe('path', () => {
  describe('pathFrom', () => {

    it('should parse a simple path', function() {
      const actual = PathFrom('spec');
      const expected = NewPath('spec');
      expect(actual).toEqual(expected);
    });

    it('should parse a path with a sub property', function() {
      const actual = PathFrom('spec.bla');
      const expected = NewPath('spec').child('bla');
      expect(actual).toEqual(expected);
    });

    it('should parse a path with an indexer', function() {
      const actual = PathFrom('spec[0].bla');
      const expected = NewPath('spec').index(0).child('bla');
      expect(actual).toEqual(expected);
    });

    it('should parse a path ending with an indexer', function() {
      const actual = PathFrom('spec.bla[0]');
      const expected = NewPath('spec', 'bla').index(0);
      expect(actual).toEqual(expected);
    });

    it('should parse a path with a string indexer', function() {
      const actual = PathFrom('spec[snip].bla');
      const expected = NewPath('spec').key('snip').child('bla');
      expect(actual).toEqual(expected);
    });

    it('should parse a path ending with a string indexer', function() {
      const actual = PathFrom('spec.bla[snip]');
      const expected = NewPath('spec', 'bla').key('snip');
      expect(actual).toEqual(expected);
    });

  });
});
