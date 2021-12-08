import React from "react";
import { fireEvent, render } from "@testing-library/react-native";

import { AddRecordScreen } from "../../src/components/screens/AddRecordScreen";
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
    fields: Array(10)
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

export const WrappedComponent = props => {
    const { control, formState } = useForm();
    return <AddRecordScreen {...props} control={control} formState={formState} />;
};

describe(AddRecordScreen.name, () => {
    test("renders correctly", () => {
        const { toJSON } = render(<WrappedComponent {...mockProps} />);
        expect(toJSON()).toMatchSnapshot();
    });

    test("displays all of a form's fields", () => {
        const { getAllByTestId } = render(<WrappedComponent {...mockProps} />);
        expect(getAllByTestId(testIds.formControl).length).toEqual(mockForm.fields.length);
    });

    test("displays a loading message", () => {
        const props = { ...mockProps, isLoading: true };
        const { getByText } = render(<WrappedComponent {...props} />);
        expect(getByText("Loading...")).toBeTruthy();
    });

    test("displays a message if there is local data", () => {
        const props = { ...mockProps, hasLocalData: true };
        const { getByText } = render(<WrappedComponent {...props} />);
        expect(getByText("There is locally stored data for this individual.")).toBeTruthy();
    });

    test("displays a message with a button if there is local data that can be submitted", () => {
        const props = { ...mockProps, hasLocalData: true };
        const { getByText, getByA11yLabel } = render(<WrappedComponent {...props} />);
        expect(getByText("Do you want to upload it?")).toBeTruthy();
        expect(getByA11yLabel("Submit local data")).toBeTruthy;
    });

    test("user pressing submit local data button triggers callback", () => {
        const props = { ...mockProps, hasLocalData: true };
        const { getByText, getByA11yLabel } = render(<WrappedComponent {...props} />);
        expect(getByText("Do you want to upload it?")).toBeTruthy();
        fireEvent.press(getByA11yLabel("Submit local data"));
        expect(mockOnSubmit).toHaveBeenCalled();
    });
});
