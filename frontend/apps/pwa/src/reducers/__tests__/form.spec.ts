import {
  findSubForm,
  hasSubFormWithId,
  selectByFolderOrDBId,
  selectFormOrSubFormById,
  selectRootForm,
} from '../form';
import { RootState } from '../../app/store';

describe('form reducer', () => {
  describe('findSubForm', () => {
    it('should find the first level sub form', () => {
      const found = findSubForm('subformId', [
        {
          fieldType: { subForm: { fields: [] } },
          id: 'subformId',
          code: 'code',
          name: 'name',
          description: 'description',
          required: false,
          key: false,
        },
      ]);
      expect(found?.id).toEqual('subformId');
    });
    it('should find a nested sub form', () => {
      const found = findSubForm('nested', [
        {
          fieldType: {
            subForm: {
              fields: [
                {
                  fieldType: {
                    subForm: {
                      fields: [],
                    },
                  },
                  id: 'nested',
                  code: 'code',
                  name: 'name',
                  description: 'description',
                  required: false,
                  key: false,
                },
              ],
            },
          },
          id: 'mainSubform',
          code: 'code',
          name: 'name',
          description: 'description',
          required: false,
          key: false,
        },
      ]);
      expect(found?.id).toEqual('nested');
    });
    it('should return undefined if not found', () => {
      const found = findSubForm('bla', [
        {
          fieldType: { subForm: { fields: [] } },
          id: 'foo',
          code: 'code',
          name: 'name',
          description: 'description',
          required: false,
          key: false,
        },
      ]);
      expect(found?.id).toBeUndefined();
    });
  });
  describe('hasSubFormWithId', () => {
    it('should return true if has subform with id', () => {
      const found = hasSubFormWithId('subformId', [
        {
          fieldType: { subForm: { fields: [] } },
          id: 'subformId',
          code: 'code',
          name: 'name',
          description: 'description',
          required: false,
          key: false,
        },
      ]);
      expect(found).toBeTruthy();
    });
    it('should return false if has subform with no such id', () => {
      const found = hasSubFormWithId('bla', [
        {
          fieldType: { subForm: { fields: [] } },
          id: 'subformId',
          code: 'code',
          name: 'name',
          description: 'description',
          required: false,
          key: false,
        },
      ]);
      expect(found).toBeFalsy();
    });
    it('should return true with a nested subform', () => {
      const found = hasSubFormWithId('nested', [
        {
          fieldType: {
            subForm: {
              fields: [
                {
                  fieldType: {
                    subForm: {
                      fields: [],
                    },
                  },
                  id: 'nested',
                  code: 'code',
                  name: 'name',
                  description: 'description',
                  required: false,
                  key: false,
                },
              ],
            },
          },
          id: 'mainSubform',
          code: 'code',
          name: 'name',
          description: 'description',
          required: false,
          key: false,
        },
      ]);
      expect(found).toBeTruthy();
    });
  });
  describe('selectRootForm', () => {
    it('should return the root form id if given the root form id', () => {
      const state = {
        forms: {
          ids: ['form'],
          entities: {
            form: { id: 'form' },
          },
        },
      } as unknown;
      const found = selectRootForm(state as RootState, 'form');
      expect(found?.id).toEqual('form');
    });
    it('should return the root form id if given the child form id', () => {
      const state = {
        forms: {
          ids: ['form'],
          entities: {
            form: {
              id: 'form',
              fields: [
                {
                  fieldType: { subForm: {} },
                  id: 'subform',
                },
              ],
            },
          },
        },
      } as unknown;
      const found = selectRootForm(state as RootState, 'subform');
      expect(found?.id).toEqual('form');
    });
  });
  describe('selectFormOrSubFormById', () => {
    it('should return the root form if given a root form id', () => {
      const state = {
        forms: {
          ids: ['form'],
          entities: {
            form: {
              id: 'form',
              fields: [{ fieldType: { subForm: { id: 'subform' } } }],
            },
          },
        },
      } as unknown;
      const found = selectFormOrSubFormById(state as RootState, 'form');
      expect(found?.id).toEqual('form');
    });
    it('should return the root form if given a child form id', () => {
      const state = {
        forms: {
          ids: ['form'],
          entities: {
            form: {
              id: 'form',
              fields: [
                {
                  fieldType: { subForm: {} },
                  id: 'subform',
                },
              ],
            },
          },
        },
      } as unknown;
      const found = selectFormOrSubFormById(state as RootState, 'subform');
      expect(found?.id).toEqual('subform');
    });
  });
  describe('selectByFolderOrDBId', () => {
    const state = {
      forms: {
        ids: ['form0', 'form1', 'form2'],
        entities: {
          form0: {
            id: 'form0',
            databaseId: 'dbId',
            folderId: 'folderId',
          },
          form1: {
            id: 'form1',
            databaseId: 'dbId',
          },
          form2: {
            id: 'form2',
            folderId: 'folderId',
          },
        },
      },
    } as unknown;
    it('should return the forms that have a folderId and a dbId', () => {
      const found = selectByFolderOrDBId(state as RootState, {
        dbId: 'dbId',
        folderId: 'folderId',
      });
      expect(found?.length).toEqual(1);
      expect(found[0].id).toEqual('form0');
    });
    it('should return the forms that have only a dbId', () => {
      const found = selectByFolderOrDBId(state as RootState, {
        dbId: 'dbId',
      });
      expect(found?.length).toEqual(1);
      expect(found[0].id).toEqual('form1');
    });
    it('should return the forms that have a folderId', () => {
      const found = selectByFolderOrDBId(state as RootState, {
        folderId: 'folderId',
      });
      expect(found?.length).toEqual(2);
      expect(found[0].id).toEqual('form0');
      expect(found[1].id).toEqual('form2');
    });
  });
});
