import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return <Path fill={variant} d="M18.9,20.1l-6,6l1.1,1.1l6-6l6,6l1.1-1.1l-6-6l6-6L26,13l-6,6l-6-6l-1.1,1.1L18.9,20.1z" />;
};
