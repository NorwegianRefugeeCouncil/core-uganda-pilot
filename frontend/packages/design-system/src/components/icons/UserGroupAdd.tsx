import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return (
    <Path
      fill={variant}
      d="M28.9,26.1H27v-1.9c0-0.1-0.1-0.2-0.2-0.2h-1.1c-0.1,0-0.2,0.1-0.2,0.2v1.9h-1.9c-0.1,0-0.2,0.1-0.2,0.2v1.1
	c0,0.1,0.1,0.2,0.2,0.2h1.9v1.9c0,0.1,0.1,0.2,0.2,0.2h1.1c0.1,0,0.2-0.1,0.2-0.2v-1.9h1.9c0.1,0,0.2-0.1,0.2-0.2v-1.1
	C29.1,26.2,29,26.1,28.9,26.1z M16.8,19.7c0-0.2,0-0.4,0-0.6c0-0.4,0-0.7,0.1-1.1c0-0.1,0-0.2-0.1-0.2c-0.3-0.1-0.6-0.3-0.9-0.6
	c-0.3-0.3-0.5-0.6-0.7-1c-0.2-0.4-0.2-0.8-0.2-1.2c0-0.8,0.3-1.5,0.9-2c0.6-0.6,1.4-0.9,2.2-0.9c0.7,0,1.5,0.3,2,0.8
	c0.2,0.2,0.3,0.4,0.5,0.6c0,0.1,0.1,0.1,0.2,0.1c0.4-0.1,0.8-0.2,1.3-0.3c0.1,0,0.2-0.2,0.1-0.3c-0.8-1.5-2.3-2.5-4.1-2.6
	c-2.6,0-4.8,2.1-4.8,4.7c0,1.5,0.7,2.8,1.7,3.6c-0.7,0.3-1.4,0.8-2,1.4c-1.3,1.3-2,3-2.1,4.8c0,0,0,0,0,0.1c0,0,0,0,0,0.1
	c0,0,0,0,0.1,0c0,0,0,0,0.1,0h1.3c0.1,0,0.2-0.1,0.2-0.2c0-1.4,0.6-2.6,1.6-3.6c0.7-0.7,1.5-1.2,2.5-1.4
	C16.7,19.9,16.8,19.8,16.8,19.7z M27.3,19.1c0-2.6-2.1-4.6-4.6-4.7c-2.6,0-4.8,2.1-4.8,4.7c0,1.5,0.7,2.8,1.7,3.6
	c-0.8,0.3-1.4,0.8-2,1.4c-1.3,1.3-2,3-2.1,4.8c0,0,0,0,0,0.1c0,0,0,0,0,0.1c0,0,0,0,0.1,0c0,0,0,0,0.1,0h1.3c0.1,0,0.2-0.1,0.2-0.2
	c0-1.4,0.6-2.6,1.6-3.6c1-1,2.4-1.6,3.8-1.6C25.2,23.8,27.3,21.7,27.3,19.1z M24.7,21.2c-0.6,0.6-1.3,0.9-2.1,0.9
	c-0.8,0-1.6-0.3-2.1-0.9c-0.3-0.3-0.5-0.6-0.7-1c-0.2-0.4-0.2-0.8-0.2-1.2c0-0.8,0.3-1.5,0.9-2.1c0.6-0.6,1.3-0.9,2.1-0.9
	c0.8,0,1.6,0.3,2.1,0.9c0.6,0.6,0.9,1.3,0.9,2.1C25.6,19.9,25.3,20.6,24.7,21.2z"
    />
  );
};
