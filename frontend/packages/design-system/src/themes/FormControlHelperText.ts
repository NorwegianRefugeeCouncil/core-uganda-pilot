import textTheme from './Text';

export default {
  baseStyle: () => {
    return {
      mt: '2px',
      _text: {
        ...textTheme.variants.caption,
        color: 'neutral.300',
      },
    };
  },
};
