import { extendTheme } from 'native-base';

import tokens from '../tokens';

import Button from './Button';
import Icon from './Icon';
import Input from './Input';
import Text from './Text';
import FormControlErrorMessage from './FormControlErrorMessage';
import FormControlHelperText from './FormControlHelperText';
import FormControlLabel from './FormControlLabel';

export default extendTheme({
  colors: tokens.colors,
  fontConfig: tokens.fontConfig,
  fontSizes: tokens.fontSizes,
  fonts: tokens.fonts,
  lineHeights: tokens.lineHeights,
  letterSpacings: tokens.letterSpacings,
  components: {
    Button,
    Icon,
    Input,
    FormControlErrorMessage,
    FormControlHelperText,
    FormControlLabel,
    Text,
  },
  space: {
    nrc_1: tokens.spacing.spacing5,
    nrc_2: tokens.spacing.spacing10,
    nrc_3: tokens.spacing.spacing15,
    nrc_4: tokens.spacing.spacing20,
    nrc_5: tokens.spacing.spacing25,
    nrc_6: tokens.spacing.spacing30,
    nrc_7: tokens.spacing.spacing35,
    nrc_8: tokens.spacing.spacing40,
    nrc_9: tokens.spacing.spacing45,
    nrc_10: tokens.spacing.spacing50,
  },
  sizes: {
    nrc_1: tokens.spacing.spacing5,
    nrc_2: tokens.spacing.spacing10,
    nrc_3: tokens.spacing.spacing15,
    nrc_4: tokens.spacing.spacing20,
    nrc_5: tokens.spacing.spacing25,
    nrc_6: tokens.spacing.spacing30,
    nrc_7: tokens.spacing.spacing35,
    nrc_8: tokens.spacing.spacing40,
    nrc_9: tokens.spacing.spacing45,
    nrc_10: tokens.spacing.spacing50,
  },
  radii: tokens.radii,
});
