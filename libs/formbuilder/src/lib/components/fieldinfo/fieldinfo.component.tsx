import React, { Fragment, FunctionComponent } from 'react';
import {
    FormLabel,
    FormInput
} from '@nrc.no/ui-toolkit';

type FieldInfoProps = {
    name: string;
    description: string;
}

export const FieldInfo: FunctionComponent<FieldInfoProps> = (props) => {
    return (
        <Fragment>
            <FormLabel>Name:</FormLabel>
            <FormInput name="name" value={props.name} />
            <br />
            <FormLabel>Description:</FormLabel>
            <FormInput name="description" value={props.description} />
        </Fragment>
    )
}