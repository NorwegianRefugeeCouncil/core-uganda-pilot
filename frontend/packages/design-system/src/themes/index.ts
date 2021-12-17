import { extendTheme } from 'native-base';

import tokens from '../tokens';

import Button from './Button';

export default extendTheme({
  colors: tokens.colors,
  components: {
    Button,
  },
  space: {
    '1': tokens.spacing['spacing-5'],
    '2': tokens.spacing['spacing-10'],
    '3': tokens.spacing['spacing-15'],
    '4': tokens.spacing['spacing-20'],
    '5': tokens.spacing['spacing-25'],
    '6': tokens.spacing['spacing-30'],
    '8': tokens.spacing['spacing-40'],
  },
});
