import {ViewStyle} from "react-native";
import React from "react";
import {Control as ControlType} from "core-js-api-client/lib/types/models";
import TextInput from "./TextInput";
import Switch from "./Switch";
import Select from "./Select";
import CheckBox from "./Checkbox";
import {Control, Controller} from "react-hook-form";

export type InputProps = {
    formControl: ControlType,
    style?: ViewStyle,
    value: any,
    onChange: any,
    onBlur?: any,
};

type FormControlProps = {
    name: string,
    formControl: ControlType,
    style?: ViewStyle,
    value?: any,
    control: Control,
};

const FormControl: React.FC<FormControlProps> = (
    {
        formControl,
        style,
        control,
        name,
        value,
    }) => {
    return (
        <Controller
            name={name as '`${string}` | `${string}.${string}` | `${string}.${number}`'}
            control={control}
            defaultValue={value}
            rules={formControl.validation}
            render={(
                {
                    field: {onChange, onBlur, value, ref},
                    fieldState: {invalid, isTouched, isDirty, error},
                    formState,
                }) => {
                switch (formControl.type) {
                    case 'checkbox':
                        return (
                            <CheckBox
                                formControl={formControl}
                                style={style}
                                value={value}
                                onBlur={onBlur}
                                onChange={onChange}/>
                        )
                    case 'boolean':
                        return (
                            <Switch
                                formControl={formControl}
                                style={style}
                                value={value}
                                onChange={onChange}/>
                        )
                    case 'dropdown':
                        return (
                            <Select
                                formControl={formControl}
                                style={style}
                                value={value}
                                onBlur={onBlur}
                                onChange={onChange}/>
                        )
                    default:
                        return (
                            <TextInput
                                formControl={formControl}
                                style={style}
                                value={value}
                                onBlur={onBlur}
                                onChange={onChange}
                            />
                        )
                }
            }}
        />
    )
};

export default FormControl;
