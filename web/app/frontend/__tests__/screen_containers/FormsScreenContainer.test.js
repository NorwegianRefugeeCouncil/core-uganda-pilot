import { render } from "@testing-library/react-native";
import React from "react";
import { FormsScreenContainer } from "../../src/components/screen_containers/FormsScreenContainer";

const mockNavigation = jest.fn();
const mockResponse = { response: { items: [] } };
const mockListForms = jest.fn(() => Promise.resolve(mockResponse));

jest.mock("../../src/utils/clients", () => {
    return () => ({
        listForms: () => mockListForms(),
    });
});

describe(FormsScreenContainer.name, () => {
    beforeAll(() => {
        render(<FormsScreenContainer navigation={mockNavigation} />);
    });
    test("fetches forms", () => {
        expect(mockListForms).toHaveBeenCalled();
    });
});
