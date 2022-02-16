import textTheme from './Text';

export default {
  baseStyle: () => {
    return {
      mt: '0.5',
      _text: {
        ...textTheme.variants.caption,
        color: 'neutral.300',
      },
    };
  },
};
