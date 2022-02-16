const baseFontSize = 18;
const baseLineHeight = 26;

const fontSizeKeys = ['3xs', '2xs', 'xs', 'sm', 'md', 'lg', 'xl', '2xl', '3xl'];
const lineHeightKeys = [
  '3xs',
  '2xs',
  'xs',
  'sm',
  'md',
  'lg',
  'xl',
  '2xl',
  '3xl',
  '4xl',
];

const fontSizeScales = [14, 16, 18, 20, 24, 28, 32, 40, 44].map(
  (fs) => fs / baseFontSize,
);
const lineHeightScales = [17, 20, 23, 24, 26, 28, 34, 38, 48, 52].map(
  (fs) => fs / baseLineHeight,
);

const fontSizes: { [p: string]: number } = fontSizeScales.reduce(
  (acc, cur, i) => ({
    ...acc,
    [fontSizeKeys[i]]: Math.round(cur * baseFontSize),
  }),
  {},
);

const lineHeights: { [p: string]: number } = lineHeightScales.reduce(
  (acc, cur, i) => ({
    ...acc,
    [lineHeightKeys[i]]: Math.round(cur * baseLineHeight),
  }),
  {},
);

export default {
  fontConfig: {
    Roboto: {
      400: {
        normal: 'Roboto_400Regular',
        italic: 'Roboto_400Regular_Italic',
      },
      500: {
        medium: 'Roboto_500Medium',
      },
      700: {
        normal: 'Roboto_700Bold',
      },
    },
  },
  fonts: {
    display1: 'Roboto',
    display2: 'Roboto',
    heading1: 'Roboto',
    heading2: 'Roboto',
    heading3: 'Roboto',
    heading4: 'Roboto',
    heading5: 'Roboto',
    heading6: 'Roboto',
    title1: 'Roboto',
    title2: 'Roboto',
    title3: 'Roboto',
    body1: 'Roboto',
    body2: 'Roboto',
    inline: 'Roboto',
    date: 'Roboto',
    label: 'Roboto',
    caption: 'Roboto',
  },
  fontSizes,
  letterSpacings: {
    xs: 0,
    sm: 0,
    md: 0,
    lg: 0,
    xl: 0,
  },
  lineHeights,
};
