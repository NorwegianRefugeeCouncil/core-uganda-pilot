import {
  FieldDefinition,
  FieldType,
  FieldValue,
  FormDefinition,
  FormType,
  FormWithRecord,
  Record,
  RecordListRequest,
} from '../../types';
import { Client } from '../Client';
import { makeField, makeForm, makeRecord } from '../../testUtils/mockData';

const client = new Client('https://www.testUrl.no');

afterEach(() => {
  jest.clearAllMocks();
});

describe('buildDefaultRecord', () => {
  it('should set the correct ids', () => {
    const record = client.Record.buildDefaultRecord({
      id: 'form-id',
      databaseId: 'database-id',
      code: 'form-code',
      folderId: 'folder-id',
      name: 'form-name',
      formType: FormType.DefaultFormType,
      fields: [],
    });
    expect(record.id).toBe('');
    expect(record.databaseId).toBe('database-id');
    expect(record.formId).toBe('form-id');
    expect(record.ownerId).toBeUndefined();
  });

  it('should set the correct field defaults', () => {
    const makeField = (fieldType: FieldType, i: number): FieldDefinition => ({
      id: `field-id-${i}`,
      code: '',
      name: `field-name-${i}`,
      description: '',
      required: false,
      key: false,
      fieldType,
    });

    const record = client.Record.buildDefaultRecord({
      id: 'form-id',
      databaseId: 'database-id',
      code: 'form-code',
      folderId: 'folder-id',
      name: 'form-name',
      formType: FormType.DefaultFormType,
      fields: [
        makeField({ text: {} }, 1),
        makeField({ multilineText: {} }, 2),
        makeField(
          { reference: { databaseId: 'database-id', formId: 'other-form-id' } },
          3,
        ),
        makeField({ subForm: { fields: [] } }, 4),
        makeField({ date: {} }, 5),
        makeField({ month: {} }, 6),
        makeField({ week: {} }, 7),
        makeField({ quantity: {} }, 8),
        makeField({ singleSelect: { options: [] } }, 9),
        makeField({ multiSelect: { options: [] } }, 10),
        makeField({ checkbox: {} }, 11),
        makeField({ quantity: {} }, 12),
      ],
    });

    expect(record.values).toEqual([
      { fieldId: 'field-id-1', value: '' },
      { fieldId: 'field-id-2', value: '' },
      { fieldId: 'field-id-3', value: null },
      { fieldId: 'field-id-4', value: [] },
      { fieldId: 'field-id-5', value: null },
      { fieldId: 'field-id-6', value: null },
      { fieldId: 'field-id-7', value: null },
      { fieldId: 'field-id-8', value: '' },
      { fieldId: 'field-id-9', value: null },
      { fieldId: 'field-id-10', value: [] },
      { fieldId: 'field-id-11', value: 'false' },
      { fieldId: 'field-id-12', value: '' },
    ]);
  });
});

describe('createWithSubRecords', () => {
  const makeSubFormFields = (count: number): FieldDefinition[] =>
    new Array(count).fill(0).map((_, i) => ({
      id: `field-id-${i}`,
      name: `field-name-${i}`,
      description: '',
      code: '',
      required: false,
      key: false,
      fieldType: {
        subForm: {
          fields: [
            {
              id: `sub-field-id-${i}`,
              name: `sub-field-name-${i}`,
              description: '',
              code: '',
              required: false,
              key: false,
              fieldType: { text: {} },
            },
          ],
        },
      },
    }));

  const makeSubFormValue = (i: number, count: number): FieldValue => ({
    fieldId: `field-id-${i}`,
    value: new Array(count)
      .fill(0)
      .map((_, j) => [
        { fieldId: `sub-field-id-${j}`, value: `sub-field-value-${j}` },
      ]),
  });

  const makeFormWithRecord = (
    fieldCount: number,
    valueCount: number,
  ): FormWithRecord<Record> => {
    const form: FormDefinition = {
      id: 'form-id',
      name: 'form-name',
      databaseId: 'database-id',
      code: '',
      folderId: '',
      formType: FormType.DefaultFormType,
      fields: makeSubFormFields(fieldCount),
    };

    const record = {
      ...client.Record.buildDefaultRecord(form),
      values: new Array(fieldCount)
        .fill(0)
        .map((_, i) => makeSubFormValue(i, valueCount)),
    };

    return {
      form,
      record,
    };
  };

  const addRecordId = (record: Record, i: number): Record => ({
    ...record,
    id: `record-id-${i}`,
  });

  const mockClientImplementation = (record: Record) => (request: any) =>
    Promise.resolve({
      success: true,
      error: undefined,
      response: record,
      request,
      status: 'ok',
      statusCode: 200,
    });

  const makeCreateRecordSpy = (records: Record[]) => {
    const createRecordSpy = jest.spyOn(client.Record, 'create');
    records.forEach((record) => {
      createRecordSpy.mockImplementationOnce(mockClientImplementation(record));
    });
    return createRecordSpy;
  };

  const makeGetRecordSpy = (record: Record) =>
    jest
      .spyOn(client.Record, 'get')
      .mockImplementationOnce(mockClientImplementation(record));

  const assertThings = (
    createRecordSpy: jest.SpyInstance,
    getRecordSpy: jest.SpyInstance,
    records: Record[],
    createdRecord: Record,
    result: Record,
  ) => {
    expect(createRecordSpy).toHaveBeenCalledTimes(records.length);
    records.forEach((r) =>
      expect(createRecordSpy).toHaveBeenCalledWith({ object: r }),
    );

    expect(getRecordSpy).toHaveBeenCalledTimes(1);
    expect(getRecordSpy).toHaveBeenCalledWith({
      recordId: createdRecord.id,
      databaseId: createdRecord.databaseId,
      formId: createdRecord.formId,
    });

    expect(result).toEqual(createdRecord);
  };

  // Sonar is an annoying little shit that wants deduplication at the expense of readability
  const makeTwoSubrecordsToAppeaseSonar = (
    inputData: FormWithRecord<Record>,
    createdRecord: Record,
    sameField: boolean,
  ) => [
    {
      id: 'created-sub-record-id-0',
      databaseId: inputData.form.databaseId,
      formId: inputData.form.fields[0].id,
      ownerId: createdRecord.id,
      values: (
        inputData.record.values[0] as {
          fieldId: string;
          value: FieldValue[][];
        }
      ).value[0],
    },
    {
      id: 'created-sub-record-id-1',
      databaseId: inputData.form.databaseId,
      formId: inputData.form.fields[sameField ? 0 : 1].id,
      ownerId: createdRecord.id,
      values: (
        inputData.record.values[sameField ? 0 : 1] as {
          fieldId: string;
          value: FieldValue[][];
        }
      ).value[sameField ? 1 : 0],
    },
  ];

  describe('success', () => {
    it('should create a record without a sub record', async () => {
      const inputData = makeFormWithRecord(0, 0);
      const createdRecord = addRecordId(inputData.record, 0);

      const createRecordSpy = makeCreateRecordSpy([createdRecord]);

      const getRecordSpy = makeGetRecordSpy(createdRecord);

      const result = await client.Record.createWithSubRecords(inputData);

      assertThings(
        createRecordSpy,
        getRecordSpy,
        [inputData.record],
        createdRecord,
        result,
      );
    });

    it('should create a record with a sub record', async () => {
      const inputData = makeFormWithRecord(1, 1);

      const createdRecord: Record = addRecordId(inputData.record, 0);

      const createdSubrecord: Record = {
        id: 'created-sub-record-id-0',
        databaseId: inputData.form.databaseId,
        formId: inputData.form.fields[0].id,
        ownerId: createdRecord.id,
        values: (
          inputData.record.values[0] as {
            fieldId: string;
            value: FieldValue[][];
          }
        ).value[0],
      };

      const createRecordSpy = makeCreateRecordSpy([
        { ...createdRecord, values: [] },
        createdSubrecord,
      ]);

      const getRecordSpy = makeGetRecordSpy(createdRecord);

      const result = await client.Record.createWithSubRecords(inputData);

      assertThings(
        createRecordSpy,
        getRecordSpy,
        [
          { ...inputData.record, values: [] },
          { ...createdSubrecord, id: '' },
        ],
        createdRecord,
        result,
      );
    });

    it('should create a record with multiple sub records of the same field', async () => {
      const inputData = makeFormWithRecord(1, 2);

      const createdRecord: Record = addRecordId(inputData.record, 0);

      const [createdSubrecordA, createdSubrecordB] =
        makeTwoSubrecordsToAppeaseSonar(inputData, createdRecord, true);

      const createRecordSpy = makeCreateRecordSpy([
        { ...createdRecord, values: [] },
        createdSubrecordA,
        createdSubrecordB,
      ]);

      const getRecordSpy = makeGetRecordSpy(createdRecord);

      const result = await client.Record.createWithSubRecords(inputData);

      assertThings(
        createRecordSpy,
        getRecordSpy,
        [
          { ...inputData.record, values: [] },
          { ...createdSubrecordA, id: '' },
          { ...createdSubrecordB, id: '' },
        ],
        createdRecord,
        result,
      );
    });

    it('should create a record with multiple sub records of different field', async () => {
      const inputData = makeFormWithRecord(2, 1);

      const createdRecord: Record = addRecordId(inputData.record, 0);

      const [createdSubrecordA, createdSubrecordB] =
        makeTwoSubrecordsToAppeaseSonar(inputData, createdRecord, false);

      const createRecordSpy = makeCreateRecordSpy([
        { ...createdRecord, values: [] },
        createdSubrecordA,
        createdSubrecordB,
      ]);

      const getRecordSpy = makeGetRecordSpy(createdRecord);

      const result = await client.Record.createWithSubRecords(inputData);

      assertThings(
        createRecordSpy,
        getRecordSpy,
        [
          { ...inputData.record, values: [] },
          { ...createdSubrecordA, id: '' },
          { ...createdSubrecordB, id: '' },
        ],
        createdRecord,
        result,
      );
    });
  });
});

describe('list', () => {
  const textfield1 = makeField(1, false, false, { text: {} });
  const textfield2 = makeField(2, false, false, { text: {} });
  const form = makeForm(1, FormType.DefaultFormType, [textfield1, textfield2]);
  const record1 = makeRecord(1, form);
  const record2 = makeRecord(2, form);

  const formGetSpy = jest.spyOn(client.Form, 'get').mockResolvedValue({
    success: true,
    error: undefined,
    response: form,
    request: { id: form.id },
    status: 'ok',
    statusCode: 200,
  });
  const doSpy = jest.spyOn(client, 'do');

  describe('success', () => {
    beforeEach(() => {
      doSpy.mockReset();
    });

    const setDoSpySuccessResponse = (
      request: RecordListRequest,
      returnValue: any,
    ) => {
      return doSpy.mockResolvedValueOnce({
        success: true,
        error: undefined,
        response: returnValue,
        request,
        status: 'ok',
        statusCode: 200,
      });
    };

    it('should return an empty array when no records exist', async () => {
      setDoSpySuccessResponse(
        {
          formId: form.id,
          databaseId: form.databaseId,
        },
        { items: [] },
      );

      const result = await client.Record.list({
        formId: form.id,
        databaseId: form.databaseId,
        subforms: false,
      });

      expect(result.response?.items).toEqual([]);
      expect(doSpy).toHaveBeenCalledWith(
        {
          formId: form.id,
          databaseId: form.databaseId,
        },
        '/records?databaseId=databaseId&formId=form1',
        'get',
        undefined,
        200,
      );
      expect(formGetSpy).not.toHaveBeenCalled();
    });

    it('should return all records', async () => {
      setDoSpySuccessResponse(
        {
          formId: form.id,
          databaseId: form.databaseId,
        },
        { items: [record1] },
      );

      const result = await client.Record.list({
        formId: form.id,
        databaseId: form.databaseId,
        subforms: false,
      });

      expect(result.response?.items).toEqual([record1]);
      expect(doSpy).toHaveBeenCalledWith(
        {
          formId: form.id,
          databaseId: form.databaseId,
        },
        '/records?databaseId=databaseId&formId=form1',
        'get',
        undefined,
        200,
      );
      expect(formGetSpy).not.toHaveBeenCalled();
    });

    it('should return all records, including subrecords', async () => {
      setDoSpySuccessResponse(
        {
          formId: form.id,
          databaseId: form.databaseId,
        },
        { items: [record1, record2] },
      );

      const result = await client.Record.list({
        formId: form.id,
        databaseId: form.databaseId,
        subforms: true,
      });

      expect(result.response?.items).toEqual([record1, record2]);
      expect(doSpy).toHaveBeenCalledWith(
        {
          formId: form.id,
          databaseId: form.databaseId,
        },
        '/records?databaseId=databaseId&formId=form1',
        'get',
        undefined,
        200,
      );
      expect(formGetSpy).toHaveBeenCalledWith({ id: form.id });
      expect(formGetSpy).toBeCalledTimes(2);
      expect(result).toEqual({
        success: true,
        error: undefined,
        response: { items: [record1, record2] },
        request: { formId: form.id, databaseId: form.databaseId },
        status: 'ok',
        statusCode: 200,
      });
    });
  });

  describe('error', () => {
    beforeEach(() => {
      doSpy.mockReset();
    });

    const setDoSpyErrorResponse = (request: RecordListRequest) => {
      return doSpy.mockResolvedValueOnce({
        success: false,
        error: 'errorMessage',
        response: undefined,
        request,
        status: 'Error Status',
        statusCode: 500,
      });
    };

    it('should return an error response if client returns error', async () => {
      setDoSpyErrorResponse({
        formId: form.id,
        databaseId: form.databaseId,
      });

      const result = await client.Record.list({
        formId: form.id,
        databaseId: form.databaseId,
      });

      expect(doSpy).toHaveBeenCalled();
      expect(result).toEqual({
        success: false,
        error: 'errorMessage',
        response: undefined,
        request: {
          formId: form.id,
          databaseId: form.databaseId,
        },
        status: 'Error Status',
        statusCode: 500,
      });
    });
  });
});
