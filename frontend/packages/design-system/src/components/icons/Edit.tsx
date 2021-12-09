import React from 'react';
import { G, Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return (
    <G>
      <Path
        fill={variant}
        d="M27.7,10.1c-0.1-0.1-0.2-0.1-0.3,0L25.6,12c-0.1,0.1-0.1,0.2,0,0.3l2.5,2.5c0.1,0.1,0.2,0.1,0.3,0l1.8-1.8
	c0.1-0.1,0.1-0.2,0-0.3L27.7,10.1z M24.6,13c0.1-0.1,0.2-0.1,0.3,0l2.5,2.5c0.1,0.1,0.1,0.2,0,0.3l-7.6,7.6l0,0L16.6,24
	c-0.1,0-0.3-0.1-0.2-0.2l0.6-3.2l0,0L24.6,13z"
      />
      <Path
        fill={variant}
        d="M19,13.6h-7.6c-0.4,0-0.8,0.3-0.8,0.8v14.2c0,0.4,0.3,0.8,0.8,0.8h14.2c0.4,0,0.8-0.3,0.8-0.8V22h-1.7v5.7H12.3
	V15.3H19V13.6z"
      />
    </G>
  );
};
