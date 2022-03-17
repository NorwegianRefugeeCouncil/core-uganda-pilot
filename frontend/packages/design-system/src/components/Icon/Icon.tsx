import React, { FC } from 'react';
import { Icon as IconNB } from 'native-base';

import { IconProps } from '../../types/icons';
import { iconMap } from '../../assets/iconMap';

export const Icon: FC<IconProps> = ({ name, customIconProps, ...props }) => {
  const MappedIcon = iconMap[name];
  return (
    <IconNB {...props}>
      {MappedIcon && <MappedIcon {...customIconProps} />}
    </IconNB>
  );
};
