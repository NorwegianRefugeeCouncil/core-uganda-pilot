import React from "react";
import { NetworkStateType } from "expo-network";
import { fireEvent, render, waitFor } from "@testing-library/react-native";
import { AddRecordScreenContainer } from "../../src/components/screen_containers/AddRecordScreenContainer";
import { RECORD_ACTIONS } from "../../src/reducers/recordsReducers";

const mockGetNetworkState = jest.fn();
jest.mock("../../src/utils/getNetworkState", () => {
    return {
        getNetworkState: () => mockGetNetworkState(),
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
    beforeEach(() => {
        // reset this before each test so we don't have to reset it manually when setting a
        // different NetworkStateType for a test
        mockGetNetworkState.mockReset().mockResolvedValue(NetworkStateType.CELLULAR);
      })

    test("renders correctly", async () => {
        const { toJSON } = await waitFor(async () => render(<Component {...mockProps} />));
        expect(toJSON()).toMatchSnapshot();
    });

    test("gets network state", async () => {
        await waitFor(async () => render(<Component {...mockProps} />));
        expect(mockGetNetworkState).toHaveBeenCalled();
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
        expect(mockGetEncryptedLocalData).toHaveBeenCalled();
    });

    test("stores local data on user submit when offline", async () => {
        mockGetNetworkState.mockReset().mockResolvedValue(NetworkStateType.NONE);

        const { getByA11yLabel } = await waitFor(async () => render(<Component {...mockProps} />));

        await fireEvent.press(getByA11yLabel("Submit"), mockLocalData);

        expect(mockHandleSubmit).toHaveBeenCalled();
        expect(mockGetEncryptionKey).toHaveBeenCalled();
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

    test("creates a db record on user submit when online", async () => {
        const { getByA11yLabel } = await waitFor(async () => render(<Component {...mockProps} />));

        await fireEvent.press(getByA11yLabel("Submit"), mockLocalData);

        expect(mockHandleSubmit).toHaveBeenCalled();

        const expectedCreateRecordArg = {
            object: { formId: mockRoute.params.formId, values: mockLocalData },
        };

        await waitFor(() => expect(mockCreateRecord).toHaveBeenCalledWith(expectedCreateRecordArg));
    });
});
