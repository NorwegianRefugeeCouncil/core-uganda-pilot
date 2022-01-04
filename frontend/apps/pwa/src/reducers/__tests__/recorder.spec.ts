import { selectPostRecords } from '../recorder';

describe('selectPostRecords', () => {
  it('should return the ordered records', function () {
    const got = selectPostRecords({
      recorder: {
        ids: ['record2', 'record1'],
        entities: {
          record1: {
            recordId: 'record1',
            formId: 'form1',
            ownerRecordId: undefined,
            values: {
              a: 'b',
            },
          },
          record2: {
            recordId: 'record2',
            formId: 'form2',
            ownerRecordId: 'record1',
            values: {
              a: 'b',
            },
          },
        },
        baseFormId: 'form1',
      },
      forms: {
        ids: ['form1', 'form2'],
        entities: {
          form1: {
            id: 'form1',
          },
          form2: {
            id: 'form2',
          },
        },
      },
    } as any);

    expect(got).toHaveLength(2);
    expect(got[0].id).toEqual('record1');
    expect(got[1].id).toEqual('record2');
  });
});
