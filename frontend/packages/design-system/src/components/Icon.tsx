import React, { FC } from 'react';
import { Icon as IconNB } from 'native-base';

import { IconProps } from '../types/icons';
import { iconMap } from '../assets/iconMap';

const Icon: FC<IconProps> = ({ name, ...props }) => {
  const MappedIcon = iconMap[name];

  return (
    <IconNB
      viewBox="0 0 40 40"
      style={{ height: '40px', width: '40px' }}
      {...props}
    >
      {MappedIcon && <MappedIcon />}
    </IconNB>
  );
};

export default Icon;
