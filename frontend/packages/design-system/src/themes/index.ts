import { extendTheme } from 'native-base';

import tokens from '../tokens';

import Button from './Button';
import Icon from './Icon';
import Input from './Input';
import Text from './Text';
import FormControlErrorMessage from './FormControlErrorMessage';
import FormControlHelperText from './FormControlHelperText';
import FormControlLabel from './FormControlLabel';
import Checkbox from './Checkbox';
import Link from './Link';

export default extendTheme({
  colors: tokens.colors,
  fontConfig: tokens.fontConfig,
  fontSizes: tokens.fontSizes,
  fontWeights: tokens.fontWeights,
  fonts: tokens.fonts,
  lineHeights: tokens.lineHeights,
  letterSpacings: tokens.letterSpacings,
  components: {
    Button,
    Checkbox,
    FormControlErrorMessage,
    FormControlHelperText,
    FormControlLabel,
    Icon,
    Input,
    Link,
    Text,
  },
  space: tokens.spacing,
  sizes: tokens.spacing,
  radii: tokens.radii,
});
