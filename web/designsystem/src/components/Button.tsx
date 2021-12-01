import React from 'react';
import {Button as ButtonRN} from 'react-native';
import {ButtonProps, ButtonVariant} from "../types/types";


const Button: React.FC<ButtonProps> = (
    {
        onPress,
        variant = ButtonVariant.PRIMARY,
        disabled = false,
        text= ''
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
        />
    );
}

export default Button
