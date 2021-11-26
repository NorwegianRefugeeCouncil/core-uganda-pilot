import React from 'react';
import { View } from 'react-native';

import { darkTheme } from '../../constants/theme';
import { InputProps } from './FormControl';

const Switch: React.FC<InputProps> = ({
    fieldDefinition,
    style,
    value,
    onChange,
}) => {
    return (
        <View style={style}>
            {fieldDefinition.label && (
                <Text theme={darkTheme}>{fieldDefinition.label[0].value}</Text>
            )}
            {fieldDefinition.description && (
                <Text theme={darkTheme} style={{ fontSize: 10 }}>
                    {fieldDefinition.description[0].value}
                </Text>
            )}
            <SwitchRNP value={value} onValueChange={onChange} />
        </View>
    );
};

export default Switch;
