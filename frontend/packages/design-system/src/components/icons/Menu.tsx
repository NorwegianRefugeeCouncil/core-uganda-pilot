import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return <Path fill={variant} d="M11,26h18v-2H11V26z M11,14v2h18v-2H11z M11,21h18v-2H11V21z" />;
};
