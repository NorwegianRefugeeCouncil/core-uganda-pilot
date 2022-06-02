import React, { FC } from 'react';
import { Icon as IconNB } from 'native-base';

import { iconMap } from '../../assets/iconMap';
import { IconProps } from '../../types/icons';

export const Logo: FC<Partial<IconProps>> = (props) => {
  const MappedIcon = iconMap.nrc;

  return (
    <IconNB {...props}>
      {MappedIcon && <MappedIcon />}
    </IconNB>
  );
};
