const makeTextVariant =
  (sizes: Record<string, string>, styles: Record<string, any>) =>
  (props: Record<string, any>) => {
    const level = (props.level || 1).toString();
    return {
      fontSize: sizes[level],
      lineHeight: sizes[level],
      ...styles,
    };
  };

const textTheme = {
  defaultProps: {
    level: '1',
  },
  baseStyle: {
    textAlign: 'auto',
    color: 'neutral.500',
    fontFamily: 'body',
    fontStyle: 'normal',
  },
  variants: {
    display: makeTextVariant(
      { '1': '3xl', '2': '2xl' },
      { fontWeight: 'bold' },
    ),
    heading: makeTextVariant(
      { '1': 'xl', '2': 'lg', '3': 'md', '4': 'sm', '5': '3xs' },
      { fontWeight: 'medium' },
    ),
    title: makeTextVariant(
      { '1': 'md', '2': 'sm', '3': '2xs' },
      { fontWeight: 'regular' },
    ),
    body: makeTextVariant({ '1': 'xs', '2': '2xs' }, { fontWeight: 'regular' }),
    caption: makeTextVariant({ '1': '3xs' }, { fontWeight: 'regular' }),
    inline: makeTextVariant(
      { '1': 'xs' },
      {
        fontWeight: 'medium',
        textDecorationLine: 'underline',
      },
    ),
    label: makeTextVariant({ '1': 'xs' }, { fontWeight: 'medium' }),
  },
};

export default textTheme;
