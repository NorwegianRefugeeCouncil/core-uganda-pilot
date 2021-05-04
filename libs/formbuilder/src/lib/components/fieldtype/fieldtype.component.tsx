import React, { Fragment, FunctionComponent } from 'react';
import {
    FormLabel,
    FormCheck,
    FormCheckLabel,
    FormRadioInput
} from '@nrc.no/ui-toolkit';

type FieldTypeProps = {
}

export enum FieldTypes {
    text = "Text",
    textarea = "Textarea",
    integer = "Integer",
    checkbox = "Checkbox",
    dropdown = "Dropdown",
    multidropdown = "Dropdown (multi)"
}

const makeFieldTypeRadios = () => {
    const returnList = []
    for (const option in FieldTypes){
        returnList.push(
            <FormCheck>
                <FormCheckLabel>{FieldTypes[option]}</FormCheckLabel>
                <FormRadioInput name={'fieldtype'} value={FieldTypes[option]}/>
            </FormCheck>
        )
    }
    return returnList
}

export const FieldType: FunctionComponent<FieldTypeProps> = (props) => {
    return (
        <Fragment>
            {
                makeFieldTypeRadios().map((radio) => radio)
            }
        </Fragment>
    )
}