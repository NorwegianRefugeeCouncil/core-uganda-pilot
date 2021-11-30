import React from 'react';
import {Button as ButtonRN} from 'react-native';

type ButtonProps = {
    onPress: () => void,
    children: any,
    variant?: 'primary' | 'secondary'
    disabled?: boolean,
    text?: string
}

const Button: React.FC<ButtonProps> = (
    {
        onPress,
        variant,
        disabled,
        text
    }
) => {

    const backgroundColor = {
        primary: 'orange',
        secondary: 'blue'
    }

    return (
        <ButtonRN
            onPress={onPress}
            color={backgroundColor[variant]}
            title={text}
            disabled={disabled}
            style={{borderColor: 'red'}}
        >
        </ButtonRN>
    );
}

export default Button
