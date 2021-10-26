import {Text, View} from "react-native";
import React from "react";
import {Checkbox as CheckboxRNP} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";
import _ from "lodash";


// FIXME: reacts only on third click
const CheckBox: React.FC<InputProps> = (
    {
        formControl,
        style,
        onChange
    }) => {

    const [isChecked, setIsChecked] = React.useState(formControl.checkboxOptions.map((o) => {
        if (formControl.value == null) {

            return o.value == 'yes'
        } else {

            // console.log(formControl.value, o.value, _.find(formControl.value, o.value))
            return !!_.find(formControl.value, (e) => {
                return e == o.value
            })
        }
    }))
    console.log(formControl.checkboxOptions.length, isChecked.length)
    return (
        <View style={style}>
            {formControl.label && (
                <Text>
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
                            <Text>
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
