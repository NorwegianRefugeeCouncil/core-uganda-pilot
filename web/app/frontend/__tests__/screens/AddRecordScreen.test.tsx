import React from "react";
import { render } from "@testing-library/react-native";

import AddRecordScreen from "../../src/components/screens/AddRecordScreen";
import testIds from "../../src/constants/testIds";
import { useForm } from "react-hook-form";

const mockField = {
    name: "something",
    descriptions: "something",
    required: true,
    options: [],
    key: false,
    fieldType: { text: {} },
};

const mockForm = {
    id: "something",
    code: "something",
    databaseId: "something",
    folderId: "something",
    name: "something",
    fields: Array(2)
        .fill(mockField)
        .map((field, index) => ({ ...field, id: index.toString(), code: index.toString() })),
};

const mockOnSubmit = jest.fn();

const mockProps = {
    form: mockForm,
    onSubmit: mockOnSubmit,
    isWeb: true,
    hasLocalData: true,
    isConnected: true,
    isLoading: false,
};

export const WrappedComponent = () => {
    const { control, formState } = useForm();
    const props = { ...mockProps, control, formState };
    return <AddRecordScreen {...props} />;
};

describe(AddRecordScreen.name, () => {
    test("renders correctly", () => {
        const { toJSON } = render(<WrappedComponent />);
        expect(toJSON()).toMatchSnapshot();
    });

    test("displays all of a form's fields", () => {
        const { getAllByTestId } = render(<WrappedComponent />);
        expect(getAllByTestId(testIds.formControl).length).toEqual(mockForm.fields.length);
    });
});
