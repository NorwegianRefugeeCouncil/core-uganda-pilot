import {Reducer} from "react";
import _ from 'lodash';
import {FormDefinition, Record} from "core-js-api-client/lib/types/types";

export type RecordsStoreProps = {
    formsById: {
        [key: string]: {
            definition: FormDefinition,
            records: Record[],
            recordsById: { [key: string]: Record; }
            localRecords: string[]
        }
    }
}

type RecordsAction = {
    type: string, // TODO use enum type
    payload: any
}

export const initialRecordsState: RecordsStoreProps = {
    formsById: {}
};

export enum RECORD_ACTIONS {
    GET_RECORDS = 'GET_RECORDS',
    GET_LOCAL_RECORDS = 'GET_LOCAL_RECORDS',
    ADD_LOCAL_RECORD = 'ADD_LOCAL_RECORD'
}

export const recordsReducer: Reducer<RecordsStoreProps, RecordsAction> = (state: RecordsStoreProps, action: RecordsAction) => {
    const {formId} = action.payload

    switch (action.type) {
        case RECORD_ACTIONS.GET_RECORDS:
            const {records} = action.payload
            return {
                formsById: {
                    ...state.formsById,
                    [formId]: {
                        ...state.formsById[formId],
                        records: records,
                        recordsById: _.keyBy(records, 'id')
                    }
                }
            }
        case RECORD_ACTIONS.GET_LOCAL_RECORDS:
            const {localRecords} = action.payload
            return {
                formsById: {
                    ...state.formsById,
                    [formId]: {
                        ...state.formsById[formId],
                        localRecords
                    }
                }
            }
        case RECORD_ACTIONS.ADD_LOCAL_RECORD:
            const currentList = state.formsById[formId].localRecords || []
            const newList = currentList.concat([action.payload.localRecord])
            return {
                formsById: {
                    ...state.formsById,
                    [formId]: {
                        ...state.formsById[formId],
                        localRecords: newList
                    }
                }
            }
        default:
            throw new Error();
    }
}
