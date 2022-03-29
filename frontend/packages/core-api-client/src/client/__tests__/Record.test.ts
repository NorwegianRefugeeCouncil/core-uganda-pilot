import {
  FieldDefinition,
  FieldType,
  FieldValue,
  FormDefinition,
  FormType,
  FormWithRecord,
  Record,
} from '../../types';
import { Client } from '../Client';

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

  const mockClientImplementation = (record: Record) => (request) =>
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

      const createdSubrecordA: Record = {
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

      const createdSubrecordB: Record = {
        id: 'created-sub-record-id-1',
        databaseId: inputData.form.databaseId,
        formId: inputData.form.fields[0].id,
        ownerId: createdRecord.id,
        values: (
          inputData.record.values[0] as {
            fieldId: string;
            value: FieldValue[][];
          }
        ).value[1],
      };

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

      const createdSubrecordA: Record = {
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

      const createdSubrecordB: Record = {
        id: 'created-sub-record-id-1',
        databaseId: inputData.form.databaseId,
        formId: inputData.form.fields[1].id,
        ownerId: createdRecord.id,
        values: (
          inputData.record.values[1] as {
            fieldId: string;
            value: FieldValue[][];
          }
        ).value[0],
      };

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
