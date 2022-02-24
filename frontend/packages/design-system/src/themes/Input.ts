import textTheme from './Text';

const disabledStyle = {
  borderColor: 'neutral.200',
  backgroundColor: 'neutral.100',
  color: 'neutral.200',
};

const activeStyle = {
  borderColor: 'primary.500',
  backgroundColor: 'secondary.100',
  color: 'neutral.500',
};

export default {
  baseStyle: {
    backgroundColor: 'white',
    borderWidth: 1,
    borderStyle: 'solid',
    borderColor: 'neutral.300',
    borderRadius: 'nrc_xs',
    color: 'neutral.300',
    padding: 3,
    _android: {
      ...textTheme.variants.body({ level: '1' }),
    },
    _ios: {
      ...textTheme.variants.body({ level: '1' }),
    },
    _web: {
      ...textTheme.variants.body({ level: '1' }),
    },
    _invalid: {
      borderColor: 'signalDanger',
      backgroundColor: 'tertiary2.100',
      color: 'signalDanger',
    },
    _focus: activeStyle,
    _hover: activeStyle,
    _disabled: disabledStyle,
    _readOnly: disabledStyle,
  },
};
