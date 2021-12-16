import { extendTheme } from 'native-base';

import tokens from '../tokens';

import Button from './Button';

export default extendTheme({
  colors: tokens.colors,
  components: {
    Button,
  },
});
