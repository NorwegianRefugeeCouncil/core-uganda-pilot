import textTheme from './Text';

export default {
  baseStyle: {
    backgroundColor: 'white',
    borderWidth: 1,
    borderStyle: 'solid',
    borderColor: 'neutral.300',
    borderRadius: 'nrc_xs',
    color: 'neutral.500',
    padding: 0,
    _text: {
      _android: {
        ...textTheme.variants.body({ level: '2' }),
      },
      _ios: {
        ...textTheme.variants.body({ level: '2' }),
      },
      _web: {
        ...textTheme.variants.body({ level: '2' }),
      },
      ml: 'nrc_2',
    },
    _checked: {
      borderColor: 'primary.500',
      borderRadius: 'nrc_xs',
      backgroundColor: 'primary.500',
      _icon: {
        backgroundColor: 'primary.500',
        name: 'check',
        color: 'white',
        viewBox: '-2 -2 18 18',
      },
    },
    _disabled: {
      opacity: 1,
      backgroundColor: 'white',
      borderColor: 'neutral.200',
      color: 'neutral.200',
    },
    _interactionBox: {
      backgroundColor: 'transparent',
      maxHeight: '100%',
    },
  },
};
