import React from "react";
import { render } from "@testing-library/react-native";
import { FormsScreen } from "../../src/components/screens/FormsScreen";
import testIds from "../../src/constants/testIds";

const mockForm = {
    formId: "blabla",
    name: "whatever",
};
const mockForms = Array(10)
    .fill(mockForm)
    .map((form, index) => ({ ...form, id: index.toString(), code: index.toString() }));
const mockIsLoading = false;
const mockNavigation = { navigate: jest.fn() };

// https://stackoverflow.com/questions/52569447/how-to-mock-react-navigations-navigation-prop-for-unit-tests-with-typescript-in
const mockProps: any = {
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
        const { getAllByTestId } =  render(<FormsScreen {...mockProps} />);
        const formItems = getAllByTestId(testIds.formListItem);
        expect(formItems.length).toBe(mockForms.length);
    });

    test("displays a message when loading", () => {
        const props = {...mockProps, isLoading: true}
        const { getByText } = render(<FormsScreen  {...props} />);;
        expect(getByText("Loading...")).toBeTruthy();
      })
});
