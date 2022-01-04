import fontConfig from '../tokens/fontConfig';

export default {
  defaultProps: {
    letterSpacing: '0em',
    textAlign: 'auto',
  },
  variants: {
    display: {
      fontFamily: fontConfig.Roboto['700'].normal,
      fontSize: '44px',
      lineHeight: '52px',
    },
    heading: {
      fontFamily: fontConfig.Roboto['400'].normal,
      fontSize: '32px',
      lineHeight: '38px',
    },
    title: {
      fontFamily: fontConfig.Roboto['400'].normal,
      fontSize: '24px',
      lineHeight: '34px',
    },
    bodyText: {
      fontFamily: fontConfig.Roboto['400'].normal,
      fontSize: '18px',
      lineHeight: '26px',
    },
    caption: {
      fontFamily: fontConfig.Roboto['400'].normal,
      fontSize: '14px',
      lineHeight: '20px',
    },
    inline: {
      fontFamily: fontConfig.Roboto['400'].normal,
      fontSize: '18px',
      lineHeight: '26px',
      textDecorationLine: 'underline',
    },
    date: {
      fontFamily: fontConfig.Roboto['400'].italic,
      fontSize: '18px',
      lineHeight: '26px',
    },
    label: {
      fontFamily: fontConfig.Roboto['400'].normal,
      fontSize: '18px',
      lineHeight: '26px',
    },
  },
};
