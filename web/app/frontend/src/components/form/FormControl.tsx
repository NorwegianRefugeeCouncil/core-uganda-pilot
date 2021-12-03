import { View, ViewStyle } from "react-native";
import React from "react";
import TextInput from "./TextInput";
import Select from "./Select";
import { Control, Controller } from "react-hook-form";
import { FieldDefinition, FieldKind, getFieldKind } from "core-js-api-client";
import ReferenceInput from "./ReferenceInput";

// TODO: move & clean up types
export type InputProps = {
    fieldDefinition: FieldDefinition;
    style?: ViewStyle;
    value: any;
    onChange: any;
    onBlur?: any;
    error?: any;
    invalid?: boolean;
    isTouched?: boolean;
    isDirty?: boolean;
    isMultiple?: boolean;
    isQuantity?: boolean;
};

type FormControlProps = {
    name: string;
    fieldDefinition: FieldDefinition;
    style?: ViewStyle;
    value?: any;
    control: Control<any, object>;
    errors?: object;
    testID?: string;
};

const FormControl = ({
    fieldDefinition,
    style,
    control,
    name,
    value,
    testID,
}: FormControlProps) => {
    return (
        // TODO: apply errors to all input types
        <View style={{ margin: 10 }} testID={testID}>
            <Controller
                name={name}
                control={control}
                defaultValue={value}
                render={({ field: { onChange, onBlur, value, ref }, fieldState, formState }) => {
                    const fieldKind = getFieldKind(fieldDefinition.fieldType);

                    switch (fieldKind) {
                        case FieldKind.Reference:
                            return (
                                <ReferenceInput
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                />
                            );
                        case FieldKind.Quantity:
                            return (
                                <TextInput
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                    isQuantity={true}
                                    {...fieldState}
                                />
                            );
                        case FieldKind.MultilineText:
                            return (
                                <TextInput
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                    isMultiple={true}
                                    {...fieldState}
                                />
                            );
                        case FieldKind.SingleSelect:
                            return (
                                <Select
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                />
                            );
                        default:
                            return (
                                <TextInput
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                    {...fieldState}
                                />
                            );
                    }
                }}
            />
        </View>
    );
};

export default FormControl;
