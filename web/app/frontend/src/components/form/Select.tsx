import {View} from "react-native";
import React from "react";
import {Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";
import {Picker} from '@react-native-picker/picker';

const Select: React.FC<InputProps> = (
    {
        fieldDefinition,
        style,
        value,
        onChange,
    }) => {

    const [selectedValue, setSelectedValue] = React.useState(value);

    return (
        <View style={style}>
            {fieldDefinition.name && <Text theme={darkTheme}>{fieldDefinition.name}</Text>}
            {fieldDefinition.description &&(
                <Text theme={darkTheme} style={{fontSize: 10}}>
                    {fieldDefinition.description}
                </Text>
            )}
            <Picker
                selectedValue={selectedValue}
                style={{height: 50, width: 150}}
                onValueChange={(itemValue) => {
                    setSelectedValue(itemValue)
                    onChange(itemValue)
                }}
            >
                {fieldDefinition.options?.map((option) => (
                    <Picker.Item
                        key={option}
                        label={option}
                        value={option}
                    />
                ))}
            </Picker>
        </View>

    );
};

export default Select;
