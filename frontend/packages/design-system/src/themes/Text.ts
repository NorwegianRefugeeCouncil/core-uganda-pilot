import typography from '../tokens/typography';

const textTheme = {
  defaultProps: {
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
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
    },
    inline: {
      fontFamily: typography.fontConfig.Roboto['500'].medium,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
      textDecorationLine: 'underline',
    },
    date: {
      fontFamily: typography.fontConfig.Roboto['400'].italic,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
    },
    label: {
      fontFamily: typography.fontConfig.Roboto['500'].medium,
      fontSize: typography.fontSizes.sm,
      lineHeight: typography.lineHeights.sm,
    },
    caption: {
      fontFamily: typography.fontConfig.Roboto['400'].normal,
      fontSize: typography.fontSizes.xs,
      lineHeight: typography.lineHeights.xs,
    },
  },
};

export default textTheme;
