import { addFormElement, formBuilderSlice, INITIAL_STATE, removeFormElement, replaceFormElement } from './index';

describe('reducer', function() {

  const reducer = formBuilderSlice.reducer;

  describe('addField', () => {

    it('should fail if the path is empty', function() {
      // empty path
      expect(() => {
        reducer(INITIAL_STATE, addFormElement({ path: '', field: {} }));
      }).toThrow('path cannot be empty');
    });

    it('should fail if the path is not starting with root', function() {
      // empty path
      expect(() => {
        reducer(INITIAL_STATE, addFormElement({ path: 'notroot', field: {} }));
      }).toThrow('first part of the path must be equal to "root", got "notroot"');
    });


    it('should add a new field to the root', function() {

      const current = {
        root: {}
      };

      const expected = {
        root: {
          children: [{ key: 'bla' }]
        }
      };

      const action = addFormElement({
        path: 'root',
        field: {
          key: 'bla'
        }
      });

      expect(reducer(current, action)).toStrictEqual(expected);
    });

    it('should add a new field to specified path', function() {

      const current = {
        root: {
          children: [{ key: 'prop1' }]
        }
      };

      const expected = {
        root: {
          children: [{ key: 'prop1', children: [{ key: 'newElement' }] }]
        }
      };

      const action = addFormElement({
        path: 'root/children/0',
        field: {
          key: 'newElement'
        }
      });

      expect(reducer(current, action)).toStrictEqual(expected);

    });
  });

  describe('removeField', () => {
    it('should remove the field', function() {

      const current = {
        root: {
          children: [{ key: 'prop1' }]
        }
      };

      const expected = {
        root: {
          children: []
        }
      };

      const action = removeFormElement({ path: '/root/children/0' });
      expect(reducer(current, action)).toStrictEqual(expected);

    });

    it('should remove the nested field', function() {

      const current = {
        root: {
          children: [
            { key: 'prop1' },
            { key: 'prop2', children: [{ key: 'prop2' }] }
          ]
        }
      };

      const expected = {
        root: {
          children: [
            { key: 'prop1' },
            { key: 'prop2', children: [] }
          ]
        }
      };

      const action = removeFormElement({ path: '/root/children/1/children/0' });
      expect(reducer(current, action)).toStrictEqual(expected);

    });

  });


  describe('updateField', () => {

    it('should update the field', function() {
      const current = { root: { 'key': 'root' } };
      const expected = { root: { 'key': 'bla' } };
      const action = replaceFormElement({ path: '/root', field: { key: 'bla' } });
      expect(reducer(current, action)).toStrictEqual(expected);
    });

    it('should update the nested field', function() {
      const current = { root: { children: [{ key: 'snip' }] } };
      const expected = { root: { children: [{ key: 'snap' }] } };
      const action = replaceFormElement({ path: '/root/children/0', field: { key: 'snap' } });
      expect(reducer(current, action)).toStrictEqual(expected);
    });

  });

});
