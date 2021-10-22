import {Picker, View} from "react-native";
import React from "react";
import {Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const Select: React.FC<InputProps> = (
    {
        formControl,
        style,
        value,
        onChange,
    }) => {

    const [selectedValue, setSelectedValue] = React.useState(value);

    return (
        <View style={style}>
            {formControl.label && <Text theme={darkTheme}>{formControl.label[0].value}</Text>}
            {formControl.description &&
            <Text theme={darkTheme} style={{fontSize: 10}}>{formControl.description[0].value}</Text>}
            <Picker
                selectedValue={selectedValue}
                style={{height: 50, width: 150}}
                onValueChange={(itemValue, itemIndex) => {
                    setSelectedValue(itemValue)
                    onChange(itemValue)
                }}
            >
                {formControl.options.map((option) => (
                    <Picker.Item
                        key={option[0].value}
                        label={option[0].value}
                        value={option[0].value}
                    />
                ))}
            </Picker>
        </View>

    );
};

export default Select;
