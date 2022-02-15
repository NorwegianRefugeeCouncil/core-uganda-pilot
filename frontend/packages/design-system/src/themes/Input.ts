const disabledStyle = {
  borderColor: 'neutral.200',
  backgroundColor: 'neutral.100',
  color: 'neutral.200',
};

const activeStyle = {
  borderColor: 'primary.500',
  backgroundColor: 'secondary.100',
};

export default {
  defaultProps: {
    size: 'xs',
    borderRadius: 'nrc_sm',
    padding: 3,
  },
  baseStyle: () => {
    return {
      borderWidth: 1,
      borderStyle: 'solid',
      borderColor: 'neutral.300',
      backgroundColor: 'white',
      boxSizing: 'border-box',
      color: 'neutral.300',
      _invalid: {
        borderColor: 'signalDanger',
        backgroundColor: 'tertiary2.100',
        color: 'signalDanger',
      },
      _focus: activeStyle,
      _hover: activeStyle,
      _disabled: disabledStyle,
      _readOnly: disabledStyle,
    };
  },
  sizes: {
    xxs: {
      fontSize: 'xxs',
      _web: {
        lineHeight: 'xxs',
      },
    },
    xs: {
      fontSize: 'xs',
      _web: {
        lineHeight: 'xs',
      },
    },
    sm: {
      fontSize: 'sm',
      _web: {
        lineHeight: 'sm',
      },
    },
    md: {
      fontSize: 'md',
      _web: {
        lineHeight: 'md',
      },
    },
    lg: {
      fontSize: 'lg',
      _web: {
        lineHeight: 'lg',
      },
    },
    xl: {
      fontSize: 'xl',
      _web: {
        lineHeight: 'xl',
      },
    },
  },
};
