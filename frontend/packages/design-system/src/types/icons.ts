import { IIconProps } from 'native-base';
import { Color } from 'react-native-svg';

import { iconMap } from '../assets/iconMap';

export type IconProps = {
  name: keyof typeof iconMap;
  customIconProps?: {
    innerColor?: Color;
  };
} & Omit<IIconProps, 'name'>;
