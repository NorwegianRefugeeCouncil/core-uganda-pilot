/*
  This input is designed to be used as a date, month, or week input on web only.
  It copies some code from native base so that we can apply the input style to a html input tag.
*/

import * as React from 'react';
import { Box, usePropsResolution, useThemeProps, useToken } from 'native-base';
import { useHover } from '@react-native-aria/interactions';
import styled from 'styled-components';
import { Platform } from 'react-native';

import { makeStyledComponent } from './makeStyledComponent';
import { useResolvedFontFamily } from './useResolvedFontFamily';

type Props = {
  type: string;
  value: string;
  invalid?: boolean;
  disabled?: boolean;
  onChange: (value: string) => void;
} & React.HTMLProps<HTMLInputElement>;

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
        overflow: visible;
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
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    outlineWidth,
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
        role="search"
        value={value || ''}
        disabled={disabled}
        onChange={handleOnChange}
        {...otherProps}
      />
    </StyledWrapper>
  );
};
