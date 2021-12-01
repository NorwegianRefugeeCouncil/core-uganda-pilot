import React from "react";
import { render } from "@testing-library/react-native";
import { FormsScreen } from "../../src/components/screens/FormsScreen";
import testIds from "../../src/constants/testIds";

const mockForm = {
    formId: "blabla",
    code: "something else",
    name: "whatever",
};
const mockForms = Array(10)
    .fill(mockForm)
    .map((form, index) => ({ ...form, id: index }));
const mockIsLoading = false;
const mockNavigation = { navigate: jest.fn() };
const mockProps = {
    forms: mockForms,
    navigation: mockNavigation,
    isLoading: mockIsLoading,
};

describe(FormsScreen.name, () => {
    it("renders correctly", () => {
        const { toJSON } = render(<FormsScreen {...mockProps} />);
        expect(toJSON()).toMatchSnapshot();
    });

    it("renders form items", () => {
        const { getAllByTestId } = render(<FormsScreen {...mockProps} />);
        const formItems = getAllByTestId(testIds.formListItem);
        expect(formItems.length).toBe(mockForms.length);
    });
});
