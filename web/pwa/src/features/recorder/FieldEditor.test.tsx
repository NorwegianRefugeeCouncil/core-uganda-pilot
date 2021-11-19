import { render } from "@testing-library/react";
import { FieldDefinition } from "../../types/types";
import { WeekFieldEditor, MonthFieldEditor } from "./FieldEditor";
import { screen } from '@testing-library/dom'
import Enzyme, { shallow } from 'enzyme';
import Adapter from '@wojtekmaj/enzyme-adapter-react-17';

Enzyme.configure({ adapter: new Adapter() });

const noOp = () => {}

const getGenericFieldDefinition = () => {
    const fd = new FieldDefinition()
    fd.id = "TEST_ID"
    fd.description = "TEST_DESCRIPTION"
    fd.name = "TEST_NAME"
    return fd
}

const getComputedStyleWorkaround = () => {
    /*
    This method stubs out the `getComputedStyle` method on the
    window object.

    Without this, the tests pass but jest throws an error due to
    the above property not being present.

    The error which is thrown is unrelated to the actual component
    being tested and has more to do with the environment in which
    the tests are run.

    A github thread of a related library introduced this workaround:
    https://github.com/nickcolley/jest-axe/issues/147
    */
    const { getComputedStyle } = window;
    window.getComputedStyle = (elt) => getComputedStyle(elt);
}

test("Week Field Editor should set attributes correctly when initialised with a null value", () => {
    getComputedStyleWorkaround();
    const mockSetValue = jest.fn();
    const fd = getGenericFieldDefinition();
    const {container, debug} = render(
        <WeekFieldEditor 
            field={fd}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    );
    const input = screen.getByRole('input', { name: fd.name });
    expect(input.getAttribute("value")).toBe("");
    expect(input.getAttribute("name")).toBe(fd.name);
    expect(input.getAttribute("id")).toBe(fd.id);
})

test("Week Field Editor should set attributes correctly when initialised with a null value", () => {
    getComputedStyleWorkaround();
    const mockSetValue = jest.fn();
    const fd = getGenericFieldDefinition();
    const testValue = new Date(2021, 0, 1);
    const rendered = render(
        <WeekFieldEditor 
            field={fd}
            value={testValue} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    );
    const input = screen.getByRole('input', { name: fd.name });
    expect(input.getAttribute("value")).toBe("2021-W01");
    expect(input.getAttribute("name")).toBe(fd.name);
    expect(input.getAttribute("id")).toBe(fd.id);
})

test("Week Field Editor should behave as expected when entering correct data", () => {
    getComputedStyleWorkaround();
    const mockSetValue = jest.fn();
    const fd = getGenericFieldDefinition();
    const shallowed = shallow(
        <WeekFieldEditor 
            field={fd}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    );
    shallowed.simulate("change", { target: { value: "2021-W01" } })
    expect(mockSetValue.mock.calls.length).toEqual(1);
})

test("Month Field Editor should set atttributes correctly when initialised with a null value", () => {
    getComputedStyleWorkaround()
    const mockSetValue = jest.fn()
    const fd = getGenericFieldDefinition();
    const monthFieldEditor = render(
        <MonthFieldEditor 
            field={fd}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    )
    const input = screen.getByRole('input', { name: fd.name });
    expect(input.getAttribute("value")).toBe("");
    expect(input.getAttribute("name")).toBe(fd.name);
    expect(input.getAttribute("id")).toBe(fd.id);
})

test("Month Field Editor should set attributes correctly when initialised with a non-null value", () => {
    getComputedStyleWorkaround()
    const mockSetValue = jest.fn()
    const fd = getGenericFieldDefinition();
    const testValue = new Date(2021, 0, 1);
    const monthFieldEditor = render(
        <MonthFieldEditor 
            field={fd}
            value={testValue} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    )
    const input = screen.getByRole('input', { name: fd.name });
    expect(input.getAttribute("value")).toBe("2021-01");
    expect(input.getAttribute("name")).toBe(fd.name);
    expect(input.getAttribute("id")).toBe(fd.id);
})

test("Month Field Editor should behave as expected when entering correct data", () => {
    getComputedStyleWorkaround();
    const mockSetValue = jest.fn();
    const fd = getGenericFieldDefinition();
    const shallowed = shallow(
        <MonthFieldEditor 
            field={fd}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    );
    shallowed.simulate("change", { target: { value: "2021-01" } })
    expect(mockSetValue.mock.calls.length).toEqual(1);
})

