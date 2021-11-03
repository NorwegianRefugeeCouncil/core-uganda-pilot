import {Text, View, ViewStyle} from "react-native";
import React from "react";
import TextInput from "./TextInput";
import Select from "./Select";
import {Control, Controller} from "react-hook-form";
import {FieldDefinition, FieldKind} from "core-js-api-client/lib/types/types";
import {getFieldKind} from "core-js-api-client/lib/client";
import ReferenceInput from "./ReferenceInput";
// import {FieldKind} from "../../../../client/src/types/types";

// TODO: move & clean up types
export type InputProps = {
    fieldDefinition: FieldDefinition,
    style?: ViewStyle,
    value: any,
    onChange: any,
    onBlur?: any,
    error?: any,
    invalid?: boolean,
    isTouched?: boolean,
    isDirty?: boolean,
    isMultiple?: boolean,
    isQuantity?: boolean
};

type FormControlProps = {
    name: string,
    fieldDefinition: FieldDefinition,
    style?: ViewStyle,
    value?: any,
    control: Control<any, object>,
    errors?: object
};

const FormControl: React.FC<FormControlProps> = (
    {
        fieldDefinition,
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
                // rules={formControl.validation}
                render={(
                    {
                        field: {onChange, onBlur, value, ref},
                        fieldState,
                        formState,
                    }) => {

                    const fieldKind = getFieldKind(fieldDefinition.fieldType);

                    switch (fieldKind) {
                        // case 'checkbox':
                        //     return (
                        //         <CheckBox
                        //             formControl={formControl}
                        //             style={style}
                        //             value={value}
                        //             onBlur={onBlur}
                        //             onChange={onChange}
                        //         />
                        //     )
                        // case 'boolean':
                        //     return (
                        //         <Switch
                        //             formControl={formControl}
                        //             style={style}
                        //             value={value}
                        //             onChange={onChange}
                        //         />
                        //     )
                        case FieldKind.SubForm:
                            // console.log('subform', fieldDefinition)
                            return <Text>sdf</Text>
                        // <TextInput
                        //     formControl={formControl}
                        //     style={style}
                        //     value={value}
                        //     onBlur={onBlur}
                        //     onChange={onChange}
                        //     {...fieldState}
                        // />

                        case FieldKind.Date:
                            // console.log('date', value, fieldDefinition)
                            return <Text>sdf</Text>
                        // <Select
                        //     fieldDefinition={fieldDefinition}
                        //     style={style}
                        //     value={value}
                        //     onBlur={onBlur}
                        //     onChange={onChange}
                        // />

                        case FieldKind.Reference:
                            // console.log('reference', fieldDefinition)
                            return (
                                <ReferenceInput
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                />
                            )
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
                            )
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
                            )
                        case FieldKind.SingleSelect:
                            return (
                                <Select
                                    fieldDefinition={fieldDefinition}
                                    style={style}
                                    value={value}
                                    onBlur={onBlur}
                                    onChange={onChange}
                                />
                            )
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
                            )
                    }
                }}
            />
        </View>
    )
};

export default FormControl;
