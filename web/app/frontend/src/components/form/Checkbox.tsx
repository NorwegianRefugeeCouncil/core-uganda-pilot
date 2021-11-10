import _ from 'lodash';
import React from 'react';
import { Text, View } from 'react-native';
import { Checkbox as CheckboxRNP } from 'react-native-paper';

import { darkTheme } from '../../constants/theme';
import { InputProps } from './FormControl';

// FIXME: reacts only on third click
const CheckBox: React.FC<InputProps> = ({
    fieldDefinition,
    style,
    onChange,
}) => {
    const [isChecked, setIsChecked] = React.useState(
        fieldDefinition.checkboxOptions.map(o => {
            if (fieldDefinition.value == null) {
                return o.value == 'yes';
            } else {
                return !!_.find(fieldDefinition.value, e => {
                    return e == o.value;
                });
            }
        })
    );
    return (
        <View style={style}>
            {fieldDefinition.label && (
                <Text>{fieldDefinition.label[0].value}</Text>
            )}
            {fieldDefinition.checkboxOptions.map((option, i) => (
                <View key={option.label[0].value}>
                    <View
                        style={{
                            display: 'flex',
                            flexDirection: 'row',
                            alignItems: 'center',
                        }}
                    >
                        <CheckboxRNP
                            theme={darkTheme}
                            status={isChecked[i] ? 'checked' : 'unchecked'}
                            onPress={() => {
                                const isCheckedtmp = isChecked;
                                isCheckedtmp[i] = !isCheckedtmp[i];
                                setIsChecked(isCheckedtmp);
                                onChange(isCheckedtmp[i]);
                            }}
                        />
                        {option.label && <Text>{option.label[0].value}</Text>}
                    </View>
                </View>
            ))}
        </View>
    );
};

export default CheckBox;
