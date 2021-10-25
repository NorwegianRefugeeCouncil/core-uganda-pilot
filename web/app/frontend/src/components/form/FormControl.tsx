import {View, ViewStyle} from "react-native";
import React from "react";
import {Control as ControlType} from "core-js-api-client/lib/types/models";
import TextInput from "./TextInput";
import Switch from "./Switch";
import Select from "./Select";
import CheckBox from "./Checkbox";
import {Control, Controller} from "react-hook-form";

// TODO: move & clean up types
export type InputProps = {
    formControl: ControlType,
    style?: ViewStyle,
    value: any,
    onChange: any,
    onBlur?: any,
    error?: any,
    invalid?: boolean,
    isTouched?: boolean,
    isDirty?: boolean
};

type FormControlProps = {
    name: string,
    formControl: ControlType,
    style?: ViewStyle,
    value?: any,
    control: Control<any, object>,
    errors?: object
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
        // TODO: apply errors to all input types
        <View style={{margin: 10}}>
            <Controller
                name={name}
                control={control}
                defaultValue={value}
                rules={formControl.validation}
                render={(
                    {
                        field: {onChange, onBlur, value, ref},
                        fieldState,
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
                                    onChange={onChange}
                                />
                            )
                        case 'boolean':
                            return (
                                <Switch
                                    formControl={formControl}
                                    style={style}
                                    value={value}
                                    onChange={onChange}
                                />
                            )
                        case 'dropdown':
                            return (
                                <Select
                                    formControl={formControl}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                />
                            )
                        default:
                            return (
                                <TextInput
                                    formControl={formControl}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                    {...fieldState}
                                />
                            )
                    }
                }}
            />
        </View>
    )
};

export default FormControl;
