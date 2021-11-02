import {View} from "react-native";
import React from "react";
import {Text, TextInput as TextInputRNP} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const TextInput: React.FC<InputProps> = (
    {
        fieldDefinition,
        style,
        value,
        onChange,
        onBlur,
        error,
        invalid,
        isTouched,
        isDirty
    }) => {

    // console.log(isDirty, isTouched, error)

    return (
        <View style={style}>
            {fieldDefinition.label && (
                <Text theme={darkTheme}>
                    {fieldDefinition.label[0].value}
                </Text>
            )}
            {fieldDefinition.description && (
                <Text theme={darkTheme} style={{fontSize: 10}}>
                    {fieldDefinition.description[0].value}
                </Text>
            )}
            <TextInputRNP
                onChangeText={onChange}
                value={value}
                onBlur={onBlur}
                error={isTouched && isDirty && error}
            />
            {isTouched && isDirty && error && (
                <Text>
                    {error.message == '' ? 'invalid' : error.message}
                </Text>
            )}
        </View>

    );
};

export default TextInput;
