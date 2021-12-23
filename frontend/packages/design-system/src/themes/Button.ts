export default {
  variants: {
    major: ({ colorScheme }: any) => {
      return {
        backgroundColor: `${colorScheme}.500`,

        _text: {
          color: 'white',
          bold: true,
        },
        _hover: {
          backgroundColor:
            colorScheme === 'secondary' ? 'link' : `${colorScheme}.200`,
          _text: {
            color: colorScheme === 'secondary' ? 'white' : `${colorScheme}.500`,
          },
        },
        _disabled: {
          backgroundColor: 'neutral.300',
          _text: {
            color: 'neutral.200',
          },
        },
        _pressed: {
          backgroundColor: `${colorScheme}.300`,
          _text: {
            color: 'white',
          },
        },
        _focus: {
          backgroundColor: `${colorScheme}.300`,
          _text: {
            color: 'white',
          },
        },
      };
    },
    minor: ({ colorScheme }: any) => {
      return {
        backgroundColor: 'white',
        _text: {
          color: `${colorScheme}.500`,
          bold: true,
        },
        borderWidth: 1,
        borderStyle: 'solid',
        borderColor: `${colorScheme}.500`,
        _hover: {
          backgroundColor:
            colorScheme === 'secondary'
              ? `${colorScheme}.100`
              : `${colorScheme}.200`,
          borderColor:
            colorScheme === 'secondary' ? 'link' : `${colorScheme}.500`,
          _text: {
            color: colorScheme === 'secondary' ? 'link' : `${colorScheme}.500`,
          },
        },
        _disabled: {
          backgroundColor: 'neutral.200',
          borderColor: 'neutral.300',
          _text: {
            color: 'neutral.300',
          },
        },
        _pressed: {
          backgroundColor: `${colorScheme}.300`,
          borderWidth: 0,
          _text: {
            color: 'white',
          },
        },
        _focus: {
          backgroundColor: `${colorScheme}.300`,
          borderWidth: 0,
          _text: {
            color: 'white',
          },
        },
      };
    },
  },
};
