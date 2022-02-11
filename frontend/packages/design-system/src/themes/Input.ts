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
  baseStyle: () => {
    return {
      borderWidth: 1,
      borderStyle: 'solid',
      borderColor: 'neutral.300',
      backgroundColor: 'white',
      color: 'neutral.300',
      padding: '11px',
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
