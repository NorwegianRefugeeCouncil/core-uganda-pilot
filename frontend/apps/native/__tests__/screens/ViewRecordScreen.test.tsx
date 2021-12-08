import { render } from "@testing-library/react-native";
import React from "react";
import { Control, useForm } from "react-hook-form";
import { FieldDefinition, FormDefinition } from "../../../client/lib/esm";
import ViewRecordScreen from "../../src/components/screens/ViewRecordScreen";
import testIds from "../../src/constants/testIds";

const mockFormId = "a form id";
const mockFolderId = "a folder id";
const mockDatabaseId = "a database id";
const mockFormCode = "a form code";
const mockFormName = "a form name";
const mockField: FieldDefinition = {
    id: "",
    code: "",
    key: false,
    name: "a field name",
    description: "a field description",
    options: [],
    required: false,
    fieldType: {text: {}},
}
const mockFields: FieldDefinition[] = Array(10).fill(mockField).map((f, i) => ({...f, id: i.toString(), code: i.toString()}));

const mockForm: FormDefinition = {
    id: mockFormId,
    folderId: mockFolderId,
    databaseId: mockDatabaseId,
    code: mockFormCode,
    name: mockFormName,
    fields: mockFields,
}


const mockProps: any = {
    isLoading: false,
    form: mockForm,
}


const HookedComponent = (props: any) => {
    const {control} = useForm();
    return <ViewRecordScreen {...props} control={control} />
}

describe(ViewRecordScreen.name, () => {

    test("renders correctly", () => {
        const { toJSON } = render(<HookedComponent {...mockProps} />);
        expect(toJSON()).toMatchSnapshot();
    })

    test("shows a message when loading", async () => {
        const { findByText } = render(<HookedComponent {...mockProps} isLoading={true} />);
        const message = await findByText("Loading...");

        expect(message).toBeTruthy();
    })

    test("shows a form's fields", async () => {
        const {findAllByTestId} = render(<HookedComponent {...mockProps} />);
        const fields = await findAllByTestId(testIds.formControl);
        expect(fields.length).toEqual(mockFields.length);
    })
})
