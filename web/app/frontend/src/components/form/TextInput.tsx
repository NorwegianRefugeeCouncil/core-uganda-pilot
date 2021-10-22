import {TextInput as TextInputRNP, View} from "react-native";
import React from "react";
import {Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const TextInput: React.FC<InputProps> = (
    {
        formControl,
        style,
        value,
        onChange,
        onBlur,
    }) => {

    return (
        <View style={style}>
            {formControl.label && (
                <Text theme={darkTheme}>
                    {formControl.label[0].value}
                </Text>
            )}
            {formControl.description && (
                <Text theme={darkTheme} style={{fontSize: 10}}>
                    {formControl.description[0].value}
                </Text>
            )}
            <TextInputRNP
                onChangeText={onChange}
                value={value}
                onBlur={onBlur}
            />
        </View>

    );
};

export default TextInput;
