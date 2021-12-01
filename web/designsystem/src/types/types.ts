import React from "react";

export enum ButtonVariant {
    PRIMARY = 'primary',
    SECONDARY = 'secondary'
}

export type ButtonProps = {
    children?: React.ReactNode,
    onPress: () => void,
    variant?: ButtonVariant
    disabled?: boolean,
    text?: string
}
