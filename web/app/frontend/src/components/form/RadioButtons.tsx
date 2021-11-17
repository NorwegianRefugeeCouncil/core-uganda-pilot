import {View} from "react-native";
import React from "react";
import {RadioButton, Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const RadioButtons: React.FC<InputProps> = ({fieldDefinition}) => {
    const [checked, setChecked] = React.useState(fieldDefinition.value[0]);
    return (
        <View>
            {fieldDefinition.label && <Text theme={darkTheme}>{fieldDefinition.label[0].value}</Text>}
            {fieldDefinition.description &&
            <Text theme={darkTheme} style={{fontSize: 10}}>{fieldDefinition.description[0].value}</Text>}

            {fieldDefinition.options.map((option) => (
                <RadioButton
                    key={option[0].value}
                    value={option[0].value}
                    status={checked === 'first' ? 'checked' : 'unchecked'}
                    onPress={() => setChecked('first')}
                />
            ))}
        </View>

    );
};

export default RadioButtons;
