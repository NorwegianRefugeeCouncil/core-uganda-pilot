import { Client } from '../Client';
import { FormGetResponse, FormType, RecordGetResponse } from '../../types';

describe('Recipient', () => {
  const client = new Client('https://www.testUrl.no');
  const recipientClient = client.Recipient;

  it('should have the expected properties', () => {
    expect(recipientClient).toHaveProperty('formClient');
    expect(recipientClient).toHaveProperty('recordClient');
  });

  it('should have the expected functions', () => {
    expect(typeof recipientClient.get).toBe('function');
    expect(typeof recipientClient.create).toBe('function');
    expect(typeof recipientClient.list).toBe('function');
  });

  describe('get', () => {
    const formGetErrorResponse: FormGetResponse = {
      error: 'formGetError',
      status: 'status',
      statusCode: 404,
      request: { id: 'formId' },
      success: false,
      response: undefined,
    };
    const formGetResponse: FormGetResponse = {
      error: undefined,
      status: 'status',
      statusCode: 200,
      request: { id: 'formId' },
      success: true,
      response: {
        id: 'id',
        fields: [],
        code: '',
        name: 'name',
        databaseId: 'databaseId',
        formType: FormType.RecipientFormType,
        folderId: 'folderId',
      },
    };

    const recordGetErrorResponse: RecordGetResponse = {
      error: 'recordGetError',
      status: 'status',
      statusCode: 404,
      request: {
        formId: 'formId',
        databaseId: 'databaseId',
        recordId: 'recordId',
      },
      success: false,
      response: undefined,
    };
    const recordGetResponse: RecordGetResponse = {
      error: undefined,
      status: 'status',
      statusCode: 200,
      request: {
        formId: 'formId',
        databaseId: 'databaseId',
        recordId: 'recordId',
      },
      success: true,
      response: {
        id: 'id',
        values: [],
        databaseId: 'databaseId',
        formId: 'formId',
        ownerId: undefined,
      },
    };

    const getFormSpy = jest.spyOn(client.Form, 'get');
    const getRecordSpy = jest.spyOn(client.Record, 'get');

    beforeEach(() => {
      getFormSpy.mockClear();
      getRecordSpy.mockClear();
    });

    it('should call formClient.get', () => {
      recipientClient.get({
        recordId: 'recordId',
        formId: 'formId',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledWith({ id: 'formId' });
    });

    it('should throw an error if formClient.get returns an error response', async () => {
      getFormSpy.mockResolvedValueOnce(formGetErrorResponse);
      try {
        await recipientClient.get({
          recordId: 'recordId',
          formId: 'formId',
          databaseId: 'databaseId',
        });
      } catch (e) {
        expect(e).toEqual(new Error('formGetError'));
      }
    });

    it('should throw an error if recordClient.get returns an error response', async () => {
      getFormSpy.mockResolvedValueOnce(formGetResponse);
      getRecordSpy.mockResolvedValueOnce(recordGetErrorResponse);
      try {
        await recipientClient.get({
          recordId: 'recordId',
          formId: 'formId',
          databaseId: 'databaseId',
        });
        expect(getFormSpy).toHaveBeenCalledWith({ id: 'formId' });
        expect(getRecordSpy).toHaveBeenCalledWith({
          recordId: 'recordId',
          formId: 'formId',
          databaseId: 'databaseId',
        });
        expect(getFormSpy).toHaveBeenCalledTimes(1);
        expect(getRecordSpy).toHaveBeenCalledTimes(1);
      } catch (e) {
        expect(e).toEqual(new Error('recordGetError'));
      }
    });

    it('should return one form and record pair if no referenceKey found in form', async () => {
      getFormSpy.mockResolvedValueOnce(formGetResponse);
      getRecordSpy.mockResolvedValueOnce(recordGetResponse);
      const result = await recipientClient.get({
        recordId: 'recordId',
        formId: 'formId',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledWith({ id: 'formId' });
      expect(getRecordSpy).toHaveBeenCalledWith({
        recordId: 'recordId',
        formId: 'formId',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledTimes(1);
      expect(getRecordSpy).toHaveBeenCalledTimes(1);
      expect(result).toEqual([
        { form: formGetResponse.response, record: recordGetResponse.response },
      ]);
    });

    it('should throw an error in case of a broken reference', async () => {
      const form1: FormGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: { id: 'form1' },
        success: true,
        response: {
          id: 'form1',
          fields: [
            {
              id: 'form1Field1',
              fieldType: {
                reference: {
                  formId: 'form2',
                  databaseId: 'databaseId',
                },
              },
              code: '',
              name: 'name',
              required: true,
              key: true,
              description: '',
            },
          ],
          code: '',
          name: 'name',
          databaseId: 'databaseId',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
        },
      };

      const record1: RecordGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: {
          formId: 'form1',
          databaseId: 'databaseId',
          recordId: 'record1',
        },
        success: true,
        response: {
          id: 'record1',
          values: [
            {
              value: 'record2',
              fieldId: 'form1Field1',
            },
          ],
          databaseId: 'databaseId',
          formId: 'form1',
          ownerId: undefined,
        },
      };
      const form2: FormGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: { id: 'form2' },
        success: true,
        response: {
          id: 'form2',
          fields: [
            {
              id: 'form2Field1',
              fieldType: {
                text: {},
              },
              code: '',
              name: 'name',
              required: true,
              key: true,
              description: '',
            },
          ],
          code: '',
          name: 'name',
          databaseId: 'databaseId',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
        },
      };

      const record2: RecordGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: {
          formId: 'form2',
          databaseId: 'databaseId',
          recordId: 'record2',
        },
        success: true,
        response: {
          id: 'record2',
          values: [],
          databaseId: 'databaseId',
          formId: 'form2',
          ownerId: undefined,
        },
      };

      getFormSpy.mockResolvedValueOnce(form1).mockResolvedValueOnce(form2);
      getRecordSpy
        .mockResolvedValueOnce(record1)
        .mockResolvedValueOnce(record2);

      try {
        await recipientClient.get({
          recordId: 'record1',
          formId: 'form1',
          databaseId: 'databaseId',
        });
      } catch (e) {
        expect(e).toEqual(new Error('broken reference'));
      }
      expect(getFormSpy).toHaveBeenCalledWith({ id: 'form1' });
      expect(getRecordSpy).toHaveBeenCalledWith({
        recordId: 'record1',
        formId: 'form1',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledWith({ id: 'form2' });
      expect(getRecordSpy).toHaveBeenCalledWith({
        recordId: 'record2',
        formId: 'form2',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledTimes(2);
      expect(getRecordSpy).toHaveBeenCalledTimes(2);
    });

    it('should return all ancestors', async () => {
      const form1: FormGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: { id: 'form1' },
        success: true,
        response: {
          id: 'form1',
          fields: [
            {
              id: 'form1Field1',
              fieldType: {
                reference: {
                  formId: 'form2',
                  databaseId: 'databaseId',
                },
              },
              code: '',
              name: 'name',
              required: true,
              key: true,
              description: '',
            },
          ],
          code: '',
          name: 'name',
          databaseId: 'databaseId',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
        },
      };

      const record1: RecordGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: {
          formId: 'form1',
          databaseId: 'databaseId',
          recordId: 'record1',
        },
        success: true,
        response: {
          id: 'record1',
          values: [
            {
              value: 'record2',
              fieldId: 'form1Field1',
            },
          ],
          databaseId: 'databaseId',
          formId: 'form1',
          ownerId: undefined,
        },
      };

      const form2: FormGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: { id: 'form2' },
        success: true,
        response: {
          id: 'form2',
          fields: [
            {
              id: 'form2Field1',
              fieldType: {
                text: {},
              },
              code: '',
              name: 'name',
              required: true,
              key: true,
              description: '',
            },
          ],
          code: '',
          name: 'name',
          databaseId: 'databaseId',
          formType: FormType.RecipientFormType,
          folderId: 'folderId',
        },
      };

      const record2: RecordGetResponse = {
        error: undefined,
        status: 'status',
        statusCode: 200,
        request: {
          formId: 'form2',
          databaseId: 'databaseId',
          recordId: 'record2',
        },
        success: true,
        response: {
          id: 'record2',
          values: [
            {
              value: 'record2',
              fieldId: 'form2Field1',
            },
          ],
          databaseId: 'databaseId',
          formId: 'form2',
          ownerId: undefined,
        },
      };

      getFormSpy.mockResolvedValueOnce(form1).mockResolvedValueOnce(form2);
      getRecordSpy
        .mockResolvedValueOnce(record1)
        .mockResolvedValueOnce(record2);

      const result = await recipientClient.get({
        recordId: 'record1',
        formId: 'form1',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledWith({ id: 'form1' });
      expect(getRecordSpy).toHaveBeenCalledWith({
        recordId: 'record1',
        formId: 'form1',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledWith({ id: 'form2' });
      expect(getRecordSpy).toHaveBeenCalledWith({
        recordId: 'record2',
        formId: 'form2',
        databaseId: 'databaseId',
      });
      expect(getFormSpy).toHaveBeenCalledTimes(2);
      expect(getRecordSpy).toHaveBeenCalledTimes(2);
      expect(result).toEqual([
        { form: form2.response, record: record2.response },
        { form: form1.response, record: record1.response },
      ]);
    });
  });
});
