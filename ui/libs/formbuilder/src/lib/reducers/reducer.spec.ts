import { addField, formBuilderReducer, INITIAL_STATE } from './index';

describe('reducer', function() {
  describe('addField', () => {
    it('should fail if the path is empty', function() {

      const state = INITIAL_STATE;
      expect(() => {
        formBuilderReducer(INITIAL_STATE, addField({
          path: '',
          field: {}
        }));
      }).toThrow();

    });
  });
});
