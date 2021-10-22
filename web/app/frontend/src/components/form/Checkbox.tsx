import {View} from "react-native";
import React from "react";
import {Checkbox as CheckboxRNP, Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const CheckBox: React.FC<InputProps> = (
    {
        formControl,
        style,
        onChange
    }) => {

    const [isChecked, setIsChecked] = React.useState(formControl.checkboxOptions.map((o) =>
        o.value == 'yes'
    ))

    return (
        <View style={style}>
            {formControl.label && (
                <Text theme={darkTheme}>
                    {formControl.label[0].value}
                </Text>
            )}
            {formControl.checkboxOptions.map((option, i) => (
                <View key={option.label[0].value}>
                    <View style={{
                        display: 'flex',
                        flexDirection: "row",
                        alignItems: "center"
                    }}>
                        <CheckboxRNP
                            theme={darkTheme}
                            status={isChecked[i] ? 'checked' : 'unchecked'}
                            onPress={() => {
                                let isCheckedtmp = isChecked;
                                isCheckedtmp[i] = !isCheckedtmp[i];
                                setIsChecked(isCheckedtmp)
                                onChange(isCheckedtmp[i])
                            }}
                        />
                        {option.label && (
                            <Text theme={darkTheme} style={{fontSize: 12}}>
                                {option.label[0].value}
                            </Text>
                        )}
                    </View>
                </View>
            ))}
        </View>
    );
};

export default CheckBox;
