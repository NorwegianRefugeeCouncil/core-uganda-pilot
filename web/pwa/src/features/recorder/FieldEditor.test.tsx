import { render } from "@testing-library/react";
import { FieldDefinition } from "../../types/types";
import { WeekFieldEditor, MonthFieldEditor } from "./FieldEditor";

const noOp = () => {}

test("Week Field Editor", () => {
    const mockSetValue = jest.fn()

    const weekFieldEditor = render(
        <WeekFieldEditor 
            field={new FieldDefinition}
            value={null} 
            setValue={mockSetValue}
            addSubRecord={noOp}
            selectSubRecord={noOp}
            subRecords={undefined}/>
    )
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