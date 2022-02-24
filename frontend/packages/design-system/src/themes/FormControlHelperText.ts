import textTheme from './Text';

export default {
  baseStyle: () => {
    return {
      mt: '0.5',
      _text: {
        ...textTheme.variants.caption({ level: '1' }),
        color: 'neutral.300',
      },
    };
  },
};
