import textTheme from './Text';

export default {
  baseStyle: () => {
    return {
      mb: '2px',
      _text: {
        ...textTheme.variants.label,
        color: 'neutral.300',
      },
    };
  },
};
