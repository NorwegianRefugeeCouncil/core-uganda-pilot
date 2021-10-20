import {View} from "react-native";
import React from "react";
import {Checkbox as CheckboxRNP, Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const CheckBox: React.FC<InputProps> = ({formControl, style}) => {
    // const [isChecked, setIsChecked] = React.useState(formControl.)
    return (
        <View style={style}>
            {formControl.checkboxOptions.map((option) => (
                <View key={option.value[0]}>
                    {formControl.label && <Text theme={darkTheme}>{formControl.label[0].value}</Text>}
                    {formControl.description &&
                    <Text theme={darkTheme} style={{fontSize: 12}}>{formControl.description[0].value}</Text>}
                    <CheckboxRNP
                        theme={darkTheme}
                        status={option.value == 'checked' ? 'checked' : 'unchecked'}
                        onPress={() => {
                        }}
                    />
                </View>
            ))}
        </View>
    );
};

export default CheckBox;
