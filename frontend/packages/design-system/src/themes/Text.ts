const makeTextVariant =
  (
    levelStyles: Record<string, Record<string, any>>,
    styles: Record<string, any>,
  ) =>
  (props: Record<string, any>) => {
    const level = (props.level || 1).toString();
    return {
      ...levelStyles[level],
      ...styles,
    };
  };

const textTheme = {
  defaultProps: {
    level: '1',
    variant: 'body',
  },
  baseStyle: {
    textAlign: 'auto',
    color: 'neutral.500',
    fontFamily: 'body',
    fontStyle: 'normal',
  },
  variants: {
    display: makeTextVariant(
      {
        '1': { fontSize: '3xl', lineHeight: '3xl' },
        '2': { fontSize: '2xl', lineHeight: '2xl' },
      },
      { fontWeight: 'bold' },
    ),
    heading: makeTextVariant(
      {
        '1': { fontSize: 'xl', lineHeight: 'xl' },
        '2': { fontSize: 'lg', lineHeight: 'lg' },
        '3': { fontSize: 'md', lineHeight: 'md' },
        '4': { fontSize: 'sm', lineHeight: 'xs' },
        '5': { fontSize: '3xs', lineHeight: '4xs' },
      },
      { fontWeight: 'medium' },
    ),
    body: makeTextVariant(
      {
        '1': { fontSize: 'xs', lineHeight: 'sm' },
        '2': { fontSize: '2xs', lineHeight: '2xs' },
      },
      { fontWeight: 'regular' },
    ),
    caption: makeTextVariant(
      {
        '1': { fontSize: '3xs', lineHeight: '3xs' },
      },
      { fontWeight: 'regular' },
    ),
    inline: makeTextVariant(
      {
        '1': { fontSize: 'xs', lineHeight: 'sm' },
      },
      {
        fontWeight: 'medium',
        textDecorationLine: 'underline',
      },
    ),
    label: makeTextVariant(
      {
        '1': { fontSize: 'xs', lineHeight: 'sm' },
      },
      { fontWeight: 'medium' },
    ),
  },
};

export default textTheme;
