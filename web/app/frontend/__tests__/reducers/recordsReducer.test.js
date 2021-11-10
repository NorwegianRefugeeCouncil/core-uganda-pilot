import {
    initialRecordsState,
    recordsReducer,
} from '../../src/reducers/recordsReducers';

describe('records reducer', () => {
    test('should return the initial state', () => {
        const action = { type: 'UNKNOWN_ACTION', payload: {} };
        expect(recordsReducer(null, action)).toEqual(initialRecordsState);
    });
    test.todo('should handle GET_RECORDS');
    test.todo('should handle GET_LOCAL_RECORDS');
    test.todo('should handle ADD_LOCAL_RECORD');
});
