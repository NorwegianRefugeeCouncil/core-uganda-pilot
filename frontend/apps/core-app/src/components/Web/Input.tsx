/*
  This input is designed to be used as a date, month, or week input on web only.
  It copies some code from native base so that we can apply the input style to a html input tag.
*/

import * as React from 'react';
import {
  Box,
  ITheme,
  usePropsResolution,
  useStyledSystemPropsResolver,
  useTheme,
  useThemeProps,
  useToken,
} from 'native-base';
import { useHover } from '@react-native-aria/interactions';
import styled from 'styled-components';
import { Platform } from 'react-native';

type Props = {
  type: string;
  value: string;
  invalid?: boolean;
  disabled?: boolean;
  onChange: (value: string) => void;
} & React.HTMLProps<HTMLInputElement>;

// native-base doesn't export this function
const makeStyledComponent = (Comp: any) => {
  // eslint-disable-next-line react/display-name
  return React.forwardRef((props: any, ref: any) => {
    const [style, restProps] = useStyledSystemPropsResolver(props);
    return (
      <Comp {...restProps} style={style} ref={ref}>
        {props.children}
      </Comp>
    );
  });
};

// native-base doesn't export this function
function useResolvedFontFamily(props: {
  fontFamily?: keyof ITheme['fonts'];
  fontStyle?: string;
  fontWeight?: keyof ITheme['fontWeights'];
}) {
  const { fontFamily, fontStyle, fontWeight } = props;
  let newFontFamily = fontFamily;
  let newFontStyle = fontStyle;
  let newFontWeight = fontWeight;

  const { fontConfig, fontWeights, fonts } = useTheme();
  if (fontWeight && fontStyle && fontFamily && fontFamily in fonts) {
    // TODO: Fix typing remove any.
    const fontToken: any = fonts[fontFamily];
    if (fontConfig && fontConfig[fontToken]) {
      // If a custom font family is resolved, set fontWeight and fontStyle to undefined.
      // https://github.com/GeekyAnts/NativeBase/issues/3811
      // On Android, If a fontFamily and fontWeight both are passed, it behaves in a weird way and applies system fonts with passed fontWeight. This happens only for some fontWeights e.g. '700' or 'bold'. So, if we find a custom fontFamily, we remove fontWeight and fontStyle
      // @ts-ignore
      newFontWeight = undefined;
      // @ts-ignore
      newFontStyle = undefined;

      const fontWeightNumber =
        fontWeight in fontWeights ? fontWeights[fontWeight] : fontWeight;
      const fontVariant = fontConfig[fontToken][fontWeightNumber];

      if (typeof fontVariant === 'object') {
        if (fontVariant[fontStyle]) newFontFamily = fontVariant[fontStyle];
      } else {
        newFontFamily = fontVariant;
      }
    } else {
      newFontFamily = fonts[fontFamily];
    }
  }

  return {
    fontFamily: newFontFamily,
    fontWeight: newFontWeight,
    fontStyle: newFontStyle,
  };
}

const StyledWrapper = makeStyledComponent(Box);

const StyledInput =
  Platform.OS === 'web'
    ? styled.input`
        font-family: inherit;
        font-weight: inherit;
        font-style: inherit;
        font-size: inherit;
        color: inherit;
        background-color: transparent;
        border: none;
        outline: none;
      `
    : null;

export const Input: React.FC<Props> = ({
  type: t,
  value,
  invalid,
  disabled,
  onChange,
  ...otherProps
}) => {
  const [isFocused, setIsFocused] = React.useState(false);

  const ref = React.useRef(null);
  const { isHovered } = useHover({}, ref);
  const inputThemeProps = useThemeProps('Input', {
    isHovered,
    isInvalid: invalid,
  });

  const {
    isFullWidth,
    isDisabled,
    isReadOnly,
    ariaLabel,
    accessibilityLabel,
    placeholderTextColor,
    selectionColor,
    underlineColorAndroid,
    type,
    fontFamily,
    fontWeight,
    fontStyle,
    ...resolvedProps
  } = usePropsResolution('Input', inputThemeProps, {
    isDisabled: false,
    isHovered,
    isFocused,
    isInvalid: invalid,
    isReadOnly: false,
  });

  const resolvedFontFamily = useResolvedFontFamily({
    fontFamily,
    fontWeight: fontWeight ?? 400,
    fontStyle: fontStyle ?? 'normal',
  });
  const resolvedPlaceholderTextColor = useToken('colors', placeholderTextColor);
  const resolvedSelectionColor = useToken('colors', selectionColor);
  const resolvedUnderlineColorAndroid = useToken(
    'colors',
    underlineColorAndroid,
  );

  const handleOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange(event.target.value);
  };

  if (!StyledInput) return null;

  return (
    <StyledWrapper
      type={type}
      value={value}
      onChange={onChange}
      secureTextEntry={type === 'password'}
      accessible
      accessibilityLabel={ariaLabel || accessibilityLabel}
      editable={!(isDisabled || isReadOnly)}
      w={isFullWidth ? '100%' : undefined}
      {...resolvedProps}
      {...resolvedFontFamily}
      placeholderTextColor={resolvedPlaceholderTextColor}
      selectionColor={resolvedSelectionColor}
      underlineColorAndroid={resolvedUnderlineColorAndroid}
      onFocus={() => setIsFocused(true)}
      onBlur={() => setIsFocused(false)}
    >
      <StyledInput
        type={t}
        value={value || ''}
        disabled={disabled}
        onChange={handleOnChange}
        {...otherProps}
      />
    </StyledWrapper>
  );
};
