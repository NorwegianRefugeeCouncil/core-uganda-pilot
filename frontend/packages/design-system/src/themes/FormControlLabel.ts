import textTheme from './Text';

export default {
  baseStyle: () => {
    return {
      mb: '0.5',
      _text: {
        ...textTheme.variants.label({ fontSize: 'xs' }),
        color: 'neutral.300',
      },
    };
  },
};
