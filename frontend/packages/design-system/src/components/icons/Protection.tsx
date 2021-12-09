import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return (
    <Path
      fill={variant}
      d="M27.7,14.5c0-0.2-0.1-0.3-0.3-0.4C27.3,14.1,27.1,14,27,14h-1.7v-2.6c0-0.4-0.3-0.8-0.8-0.8h-9c-0.4,0-0.8,0.3-0.8,0.8V14
	H13c-0.2,0-0.3,0-0.4,0.1c-0.1,0.1-0.2,0.2-0.3,0.4l-1.7,5.2v8.9c0,0.4,0.3,0.8,0.8,0.8h17.2c0.4,0,0.8-0.3,0.8-0.8v-8.9L27.7,14.5z
	 M16.4,12.3h7.1V14h-7.1V12.3z M27.7,27.7H12.3V20l1.4-4.3h12.5l1.4,4.3V27.7z M23.3,21.4h-2.5v-2.5c0-0.1-0.1-0.2-0.2-0.2h-1.1
	c-0.1,0-0.2,0.1-0.2,0.2v2.5h-2.5c-0.1,0-0.2,0.1-0.2,0.2v1.1c0,0.1,0.1,0.2,0.2,0.2h2.5v2.5c0,0.1,0.1,0.2,0.2,0.2h1.1
	c0.1,0,0.2-0.1,0.2-0.2v-2.5h2.5c0.1,0,0.2-0.1,0.2-0.2v-1.1C23.5,21.5,23.4,21.4,23.3,21.4z"
    />
  );
};
