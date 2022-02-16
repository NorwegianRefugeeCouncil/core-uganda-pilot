import typography from '../tokens/typography';

const textTheme = {
  baseStyle: {
    textAlign: 'auto',
    color: 'neutral.500',
  },
  variants: {
    display1: {
      fontFamily: typography.fontConfig.Roboto['700'].normal,
      fontSize: typography.fontSizes['3xl'],
      lineHeight: typography.lineHeights['4xl'],
    },
    display2: {
      fontFamily: typography.fontConfig.Roboto['700'].normal,
      fontSize: typography.fontSizes['2xl'],
      lineHeight: typography.lineHeights['3xl'],
    },
    heading1: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.xl,
      lineHeight: typography.lineHeights['2xl'],
    },
    heading2: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.lg,
      lineHeight: typography.lineHeights.xl,
    },
    heading3: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.md,
      lineHeight: typography.lineHeights.lg,
    },
    heading4: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
    },
    heading5: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes['3xs'],
      lineHeight: typography.lineHeights['3xs'],
    },
    heading6: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
    },
    title1: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.md,
      lineHeight: typography.lineHeights.xl,
    },
    title2: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.lg,
    },
    title3: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes['2xs'],
      lineHeight: typography.lineHeights.xs,
    },
    body1: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.md,
    },
    body2: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes['2xs'],
      lineHeight: typography.lineHeights.xs,
    },
    inline: {
      fontFamily: typography.fontConfig.Roboto['500'].medium,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.md,
      textDecorationLine: 'underline',
    },
    date1: {
      fontFamily: typography.fontConfig.Roboto['400'].italic,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.md,
    },
    date2: {
      fontFamily: typography.fontConfig.Roboto['400'].italic,
      fontSize: typography.fontSizes['2xs'],
      lineHeight: typography.lineHeights.md,
    },
    label: {
      fontFamily: typography.fontConfig.Roboto['500'].medium,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.md,
    },
    caption: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes['3xs'],
      lineHeight: typography.lineHeights['2xs'],
    },
  },
};

export default textTheme;
