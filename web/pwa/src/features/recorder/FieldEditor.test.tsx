import { render } from "@testing-library/react";
import { FieldDefinition } from "../../types/types";
import { WeekFieldEditor, MonthFieldEditor } from "./FieldEditor";

const noOp = () => {}

const getWeekFieldDefinition = () => {
    const fd = new FieldDefinition()
    fd.id = "TEST_ID"
    fd.description = "TEST_DESCRIPTION"
    fd.name = "TEST_NAME"
    return fd
}

test("Week Field Editor - Null Value", () => {
    const mockSetValue = jest.fn()
    const fd = getWeekFieldDefinition()
    const {container, debug} = render(
        <WeekFieldEditor 
            field={fd}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    )
    debug()
})

test("Week Field Editor - Non-null Value", () => {
    const mockSetValue = jest.fn()
    const fd = getWeekFieldDefinition()
    const testValue = new Date(2021, 0, 1)
    const {container, debug} = render(
        <WeekFieldEditor 
            field={fd}
            value={testValue} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    )
    debug()
})

test("Month Field Editor", () => {
    const mockSetValue = jest.fn()

    const monthFieldEditor = render(
        <MonthFieldEditor 
            field={new FieldDefinition}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    )
})