export enum ButtonVariant {
    PRIMARY = 'primary',
    SECONDARY = 'secondary'
}

export type ButtonProps = {
    onPress: () => void,
    children: any,
    variant?: ButtonVariant
    disabled?: boolean,
    text?: string
}
