import React, { Fragment, FunctionComponent } from 'react';
import {
    FormLabel,
    FormCheck,
    FormCheckLabel,
    FormCheckInput,
    FormInput
} from '@nrc.no/ui-toolkit';
import { FieldType } from '../fieldtype/fieldtype.component';

type FieldConfigProps = {
    fieldType: FieldType
    fieldProps: GenericFieldConfig
}   

interface DateConfig {
    mindate: string,
    maxdate: string,
    default: string,
    placeholder: string,
    required: boolean
}

const renderDateConfig = (props: DateConfig) => {
    return <Fragment>
        <FormLabel>Minimum date:</FormLabel>
        <FormInput name="mindate" type={"date"} value={props.mindate} />
        <FormLabel>Maximum Length:</FormLabel>
        <FormInput name="maxldate" type={"date"} value={props.maxdate} />
        <FormLabel>Default:</FormLabel>
        <FormInput name="default" type={"date"} value={props.default} />
        <FormLabel>Placeholder:</FormLabel>
        <FormInput name="placeholder" value={props.placeholder} />
        {
            props.required === true ? 
            <FormCheck>
                <FormCheckLabel>Required</FormCheckLabel>
                <FormCheckInput name="required" checked={true} />
            </FormCheck> :
            <FormCheck>
                <FormCheckLabel>Required</FormCheckLabel>
                <FormCheckInput name="required" />
            </FormCheck>
        }
        
    </Fragment>
}

interface TextConfig {
    minlength: number,
    maxlength: number,
    default: string,
    placeholder: string,
    required: boolean,
    regex: string
}

const renderTextConfig = (props: TextConfig) => {
    return <Fragment>
            <FormLabel>Minimum Length:</FormLabel>
            <FormInput name="minlength" type={"number"} value={props.minlength} />
            <FormLabel>Maximum Length:</FormLabel>
            <FormInput name="maxlength" type={"number"} value={props.maxlength} />
            <FormLabel>Default:</FormLabel>
            <FormInput name="default" value={props.default} />
            <FormLabel>Placeholder:</FormLabel>
            <FormInput name="placeholder" value={props.placeholder} />
            {
                props.required === true ? 
                <FormCheck>
                    <FormCheckLabel>Required</FormCheckLabel>
                    <FormCheckInput name="required" checked={true} />
                </FormCheck> :
                <FormCheck>
                    <FormCheckLabel>Required</FormCheckLabel>
                    <FormCheckInput name="required" />
                </FormCheck>
            }
            <FormLabel>Regex:</FormLabel>
            <FormInput name="regex" value={props.placeholder} />
    </Fragment>
}
const renderTextAreaConfig = (props: TextConfig) => {
    return renderTextConfig(props)
}

interface IntegerConfig {
    minval: number,
    maxval: number,
    default: number,
    placeholder: string,
    required: boolean
}

const renderIntegerConfig = (props: IntegerConfig) => {
    return <Fragment>
            <FormLabel>Minimum Value:</FormLabel>
            <FormInput name="minval" type={"number"} value={props.minval} />
            <FormLabel>Maximum Value:</FormLabel>
            <FormInput name="maxval" type={"number"} value={props.maxval} />
            <FormLabel>Default:</FormLabel>
            <FormInput name="default" type={"number"} value={props.default} />
            <FormLabel>Placeholder:</FormLabel>
            <FormInput name="placeholder" value={props.placeholder} />
            {
                props.required === true ? 
                <FormCheck>
                    <FormCheckLabel>Required</FormCheckLabel>
                    <FormCheckInput name="required" checked={true} />
                </FormCheck> :
                <FormCheck>
                    <FormCheckLabel>Required</FormCheckLabel>
                    <FormCheckInput name="required" />
                </FormCheck>
            }
    </Fragment>
}

interface CheckboxConfig {
    default: boolean
}

const renderCheckboxConfig = (props: CheckboxConfig) => {
    return <Fragment>
            {
                props.default === true ?
                <FormCheck>
                    <FormCheckLabel>Default</FormCheckLabel>
                    <FormCheckInput name="default" checked={true} />
                </FormCheck> :
                <FormCheck>
                    <FormCheckLabel>Default</FormCheckLabel>
                    <FormCheckInput name="default" />
                </FormCheck>
            }
    </Fragment>
}

interface DropdownConfig {
    default: string,
    placeholder: string,
    required: boolean
}

const renderDropdownConfig = (props: DropdownConfig) => {
    return <Fragment>
            <FormLabel>Default:</FormLabel>
            <FormInput name="default" type={"string"} value={props.default}/>
            <FormLabel>Placeholder:</FormLabel>
            <FormInput name="placeholder" value={props.placeholder} />
            {
                props.required === true ? 
                <FormCheck>
                    <FormCheckLabel>Required</FormCheckLabel>
                    <FormCheckInput name="required" checked={true} />
                </FormCheck> :
                <FormCheck>
                    <FormCheckLabel>Required</FormCheckLabel>
                    <FormCheckInput name="required" />
                </FormCheck>
            }
    </Fragment>
}
const renderMultiDropdownConfig = (props: DropdownConfig) => {
    return renderDropdownConfig(props)
}

const renderConfig = (fieldType: FieldType, fieldProps: GenericFieldConfig) => {
    switch (fieldType) {
        case FieldType.date:
            return renderDateConfig(fieldProps as DateConfig)
        case FieldType.text:
            return renderTextConfig(fieldProps as TextConfig)
        case FieldType.textarea:
            return renderTextAreaConfig(fieldProps as TextConfig)
        case FieldType.integer:
            return renderIntegerConfig(fieldProps as IntegerConfig)
        case FieldType.checkbox:
            return renderCheckboxConfig(fieldProps as CheckboxConfig)
        case FieldType.dropdown:
            return renderDropdownConfig(fieldProps as DropdownConfig)
        case FieldType.multidropdown:
            return renderMultiDropdownConfig(fieldProps as DropdownConfig)
    
        default:
            return (
                <p>No render function found</p>
            )
    }
}

export type GenericFieldConfig = DateConfig | TextConfig | IntegerConfig | CheckboxConfig | DropdownConfig

export const FieldConfig: FunctionComponent<FieldConfigProps> = (props) => {
    return (
        <Fragment>
            {
                renderConfig(props.fieldType, props.fieldProps)
            }
        </Fragment>
    )
}