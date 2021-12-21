import { extendTheme } from 'native-base';

import tokens from '../tokens';

import Button from './Button';

export default extendTheme({
  colors: tokens.colors,
  components: {
    Button,
  },
  space: {
    '1': tokens.spacing.spacing5,
    '2': tokens.spacing.spacing10,
    '3': tokens.spacing.spacing15,
    '4': tokens.spacing.spacing20,
    '5': tokens.spacing.spacing25,
    '6': tokens.spacing.spacing30,
    '8': tokens.spacing.spacing40,
  },
});