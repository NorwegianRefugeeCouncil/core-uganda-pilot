import React from "react";
import { NetworkStateType } from "expo-network";
import { fireEvent, render, waitFor } from "@testing-library/react-native";
import { AddRecordScreenContainer } from "../../src/components/screen_containers/AddRecordScreenContainer";
import { RECORD_ACTIONS } from "../../src/reducers/recordsReducers";

const mockGetNetworkStateCellular = jest.fn().mockResolvedValue(NetworkStateType.CELLULAR);
const mockGetNetworkStateOffline = jest.fn().mockResolvedValue(NetworkStateType.NONE);
jest.mock("../../src/utils/getNetworkState", () => {
    return {
        getNetworkState: () => mockGetNetworkStateCellular(),
    };
});

const mockHandleSubmit = jest.fn(fn => fn());
const mockReset = jest.fn();
const mockFormState = {};
const mockControl = {};

jest.mock("react-hook-form", () => {
    const hf = jest.requireActual("react-hook-form");
    return {
        ...hf,
        useForm: () => ({
            handleSubmit: (fn: any) => mockHandleSubmit(fn),
            formState: mockFormState,
            reset: (data: any) => mockReset(data),
            control: mockControl,
        }),
    };
});

interface Data {
    [key: string]: any;
}

const mockForm = {};
const mockGetForm = jest.fn().mockResolvedValue(mockForm);
const mockCreateRecord = jest.fn().mockResolvedValue(true);
jest.mock("../../src/utils/clients", () => {
    return () => ({
        getForm: (query: any) => mockGetForm(query),
        createRecord: (record: any) => mockCreateRecord(record),
    });
});

const mockEncryptionKey = "super snazzy secret key";
const mockGetEncryptionKey = jest.fn().mockReturnValue(mockEncryptionKey);
jest.mock("../../src/utils/getEncryptionKey", () => {
    return {
        getEncryptionKey: () => mockGetEncryptionKey(),
    };
});

const mockLocalData: Data = { arbitrary: "data" };
const mockGetEncryptedLocalData = jest.fn().mockResolvedValue(mockLocalData);
const mockStoreEncryptedLocalData = jest.fn().mockResolvedValue(true);
jest.mock("../../src/utils/storage", () => {
    return {
        getEncryptedLocalData: () => mockGetEncryptedLocalData(),
        storeEncryptedLocalData: (recordId: any, key: any, data: any) =>
            mockStoreEncryptedLocalData(recordId, key, data),
    };
});

const mockRoute = {
    params: { formId: "randomId", recordId: "randomId" },
};

const mockDispatch = jest.fn();

const mockProps = {
    route: mockRoute,
    dispatch: (action: any) => mockDispatch(action),
};

const Component = (props: any) => <AddRecordScreenContainer {...props} />;

describe(AddRecordScreenContainer.name, () => {
    test("renders correctly", async () => {
        const { toJSON } = await waitFor(async () => render(<Component {...mockProps} />));
        expect(toJSON()).toMatchSnapshot();
    });

    test("gets network state", async () => {
        await waitFor(async () => render(<Component {...mockProps} />));
        expect(mockGetNetworkStateCellular).toHaveBeenCalledTimes(1);
    });

    test("fetches a form with the supplied formId", async () => {
        await waitFor(async () => render(<Component {...mockProps} />));
        expect(mockGetForm).toHaveBeenCalledWith({ id: mockRoute.params.formId });
    });

    test("calls useForm's reset method with local data", async () => {
        await waitFor(async () => render(<Component {...mockProps} />));
        expect(mockReset).toHaveBeenCalledWith(mockLocalData);
    });

    test("fetches locally stored date when on mobile", async () => {
        await waitFor(async () => render(<Component {...mockProps} />));
        expect(mockGetEncryptedLocalData).toHaveBeenCalledTimes(1);
    });

    test("stores local data on user submit when offline", async () => {
        jest.mock("../../src/utils/getNetworkState", () => ({
            getNetworkState: () => mockGetNetworkStateOffline(),
        }));

        const { getByA11yLabel } = await waitFor(async () => render(<Component {...mockProps} />));

        await fireEvent.press(getByA11yLabel("Submit"), mockLocalData);

        expect(mockHandleSubmit).toHaveBeenCalledTimes(1);
        expect(mockGetEncryptionKey).toHaveBeenCalledTimes(1);
        expect(mockStoreEncryptedLocalData).toHaveBeenCalledWith(
            mockRoute.params.recordId,
            mockEncryptionKey,
            mockLocalData
        );
        const expectedAddRecordAction = {
            type: RECORD_ACTIONS.ADD_LOCAL_RECORD,
            payload: {
                formId: mockRoute.params.formId,
                localRecord: mockRoute.params.recordId,
            },
        };
        await waitFor(() => expect(mockDispatch).toHaveBeenCalledWith(expectedAddRecordAction));
    });

    test.only("creates a db record on user submit when online", async () => {
        jest.mock("../../src/utils/getNetworkState", () => ({
            getNetworkState: () => mockGetNetworkStateCellular(),
        }));

        const { getByA11yLabel } = await waitFor(async () => render(<Component {...mockProps} />));

        await fireEvent.press(getByA11yLabel("Submit"), mockLocalData);

        expect(mockHandleSubmit).toHaveBeenCalledTimes(1);

        const expectedCreateRecordArg = {
            object: { formId: mockRoute.params.formId, values: mockLocalData },
        };

        await waitFor(() => expect(mockCreateRecord).toHaveBeenCalledWith(expectedCreateRecordArg));
        await waitFor(() => expect(mockCreateRecord).toHaveBeenCalledTimes(1));
    });
});
