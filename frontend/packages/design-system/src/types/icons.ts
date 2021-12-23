import { IIconProps } from 'native-base';

import { iconMap } from '../assets/iconMap';

export type IconProps = {
  name: keyof typeof iconMap;
} & Omit<IIconProps, 'name'>;
