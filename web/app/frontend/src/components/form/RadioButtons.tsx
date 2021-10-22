import {View} from "react-native";
import React from "react";
import {RadioButton, Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const RadioButtons: React.FC<InputProps> = ({formControl}) => {
    const [checked, setChecked] = React.useState(formControl.value[0]);
    return (
        <View>
            {formControl.label && <Text theme={darkTheme}>{formControl.label[0].value}</Text>}
            {formControl.description &&
            <Text theme={darkTheme} style={{fontSize: 10}}>{formControl.description[0].value}</Text>}

            {formControl.options.map((option) => (
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
