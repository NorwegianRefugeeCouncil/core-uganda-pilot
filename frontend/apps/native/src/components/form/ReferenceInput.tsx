import {View} from "react-native";
import React from "react";
import {Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const ReferenceInput: React.FC<InputProps> = (
    {
        fieldDefinition,
        style,
        error,
        isTouched,
        isDirty,
    }) => {
    return (
        <View style={style}>
            {fieldDefinition.name && <Text theme={darkTheme}>{fieldDefinition.name}</Text>}
            {fieldDefinition.description && (
                <Text theme={darkTheme} style={{fontSize: 10}}>
                    {fieldDefinition.description}
                </Text>
            )}
            {isTouched && isDirty && error && (
                <Text>
                    {error.message == '' ? 'invalid' : error.message}
                </Text>
            )}
        </View>

    );
};

export default ReferenceInput;
