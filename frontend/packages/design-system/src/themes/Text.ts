import typography from '../tokens/typography';

const textTheme = {
  baseStyle: {
    textAlign: 'auto',
  },
  variants: {
    display: {
      fontFamily: typography.fontConfig.Roboto['700'].normal,
      fontSize: typography.fontSizes.xl,
      lineHeight: typography.lineHeights.xl,
    },
    heading: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.lg,
      lineHeight: typography.lineHeights.lg,
    },
    title: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.md,
      lineHeight: typography.lineHeights.md,
    },
    bodyText: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.xs,
    },
    inline: {
      fontFamily: typography.fontConfig.Roboto['500'].medium,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.xs,
      textDecorationLine: 'underline',
    },
    date: {
      fontFamily: typography.fontConfig.Roboto['400'].italic,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.xs,
    },
    label: {
      fontFamily: typography.fontConfig.Roboto['500'].medium,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.xs,
    },
    caption: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.xxs,
      lineHeight: typography.lineHeights.xxs,
    },
    button: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
      bold: 'true',
    },
  },
};

export default textTheme;
