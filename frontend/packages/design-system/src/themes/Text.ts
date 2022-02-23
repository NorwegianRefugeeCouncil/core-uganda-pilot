const makeTextVariant =
  (styles: Record<string, any>) => (props: Record<string, any>) => ({
    lineHeight: props.fontSize,
    ...styles,
  });

const textTheme = {
  baseStyle: {
    textAlign: 'auto',
    color: 'neutral.500',
    fontFamily: 'body',
    fontStyle: 'normal',
  },
  variants: {
    display: makeTextVariant({ fontWeight: 'bold' }),
    heading: makeTextVariant({ fontWeight: 'medium' }),
    title: makeTextVariant({ fontWeight: 'regular' }),
    body: makeTextVariant({ fontWeight: 'regular' }),
    inline: makeTextVariant({
      fontWeight: 'medium',
      textDecorationLine: 'underline',
    }),
    date: makeTextVariant({ fontWeight: 'regular', fontStyle: 'italic' }),
    label: makeTextVariant({ fontWeight: 'medium' }),
    caption: makeTextVariant({ fontWeight: 'regular' }),
  },
};

export default textTheme;
