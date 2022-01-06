const baseFontSize = 18;

const keys = ['xs', 'sm', 'md', 'lg', 'xl'];

const fontSizeScales = [14 / 18, 1, 24 / 18, 32 / 18, 44 / 18];
const lineHeightScales = [20 / 26, 1, 34 / 18, 38 / 18, 52 / 18];

const fontSizes: { [p: string]: number } = fontSizeScales.reduce(
  (acc, cur, i) => ({
    ...acc,
    [keys[i]]: Math.round(cur * baseFontSize),
  }),
  {},
);

const lineHeights: { [p: string]: number } = lineHeightScales.reduce(
  (acc, cur, i) => ({
    ...acc,
    [keys[i]]: Math.round(cur * baseFontSize),
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
      700: {
        normal: 'Roboto_700Bold',
      },
    },
  },
  fonts: {
    display: 'Roboto',
    heading: 'Roboto',
    title: 'Roboto',
    bodyText: 'Roboto',
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
