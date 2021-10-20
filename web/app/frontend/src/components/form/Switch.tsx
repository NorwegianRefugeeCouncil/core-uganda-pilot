import {View} from "react-native";
import React from "react";
import {Switch as SwitchRNP, Text} from "react-native-paper";
import {darkTheme} from "../../constants/theme";
import {InputProps} from "./FormControl";

const Switch: React.FC<InputProps> = ({formControl, style}) => {
    return (
        <View style={style}>
            {formControl.label && <Text theme={darkTheme}>{formControl.label[0].value}</Text>}
            {formControl.description &&
            <Text theme={darkTheme} style={{fontSize: 10}}>{formControl.description[0].value}</Text>}
            <SwitchRNP
                value={formControl.value ? formControl.value[0] == 'true' : false}
                // name={formControl.name}
                // placeholder={formControl.placeholder.toLocaleString()}
                // type={'text'}
            />
        </View>

    );
};

export default Switch;
