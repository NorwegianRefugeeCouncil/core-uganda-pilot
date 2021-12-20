import React, { FC } from 'react';
import { Icon as IconNB, IIconProps } from 'native-base';

import { iconMap } from '../assets/iconMap';
import { IconName } from '../types/icons';

type Props = {
  name: keyof typeof iconMap;
} & Omit<IIconProps, 'name'> & { name: IconName };

const Icon: FC<Props> = (props) => {
  const { name } = props;
  const MappedIcon = iconMap[name];

  return (
    <IconNB
      viewBox="0 0 40 40"
      style={{
        height: '40px',
        width: '40px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
      {...props}
    >
      {MappedIcon}
    </IconNB>
  );
};

export default Icon;
