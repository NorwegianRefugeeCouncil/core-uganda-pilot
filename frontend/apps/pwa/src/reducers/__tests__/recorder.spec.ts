import { selectPostRecords } from '../Recorder/recorder.selectors';

describe('selectPostRecords', () => {
  it('should return the ordered records', function () {
    const got = selectPostRecords({
      recorder: {
        ids: ['record2', 'record1'],
        entities: {
          record1: {
            id: 'record1',
            formId: 'form1',
            ownerId: undefined,
            values: {
              a: 'b',
            },
          },
          record2: {
            id: 'record2',
            formId: 'form2',
            ownerId: 'record1',
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
