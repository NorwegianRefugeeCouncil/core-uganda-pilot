import React from 'react';
import {TouchableHighlight} from 'react-native';

type ButtonProps = {
    onPress: () => void,
    children: any
}

const Button: React.FC<ButtonProps> = ({onPress, children}) => {
    return (
        <TouchableHighlight
            onPress={onPress}
        >
            {children}
        </TouchableHighlight>
    );
}

export default Button
