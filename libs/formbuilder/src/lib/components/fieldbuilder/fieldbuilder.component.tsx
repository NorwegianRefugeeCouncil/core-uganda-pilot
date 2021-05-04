import React, { FunctionComponent } from 'react';
import {
    Card,
    FormLabel,
    FormSelect
} from '@nrc.no/ui-toolkit'

type FieldBuilderProps = {
}

export const FieldBuilder: FunctionComponent<FieldBuilderProps> = (props) => {
    return (
        <Card>
            <FormLabel>Type:</FormLabel>
            <FormSelect>
                <option>Short text</option>
                <option>Long text</option>
                <option>Checkbox</option>
                <option>Number</option>
            </FormSelect>
        </Card>
    )
}