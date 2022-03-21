/*
  This is a function copied from native base, but is not exported.
  It is required to apply the input style to a html input tag.
  It is not intended to be used outside of this component.
*/

import { ITheme, useTheme } from 'native-base';

export function useResolvedFontFamily(props: {
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
