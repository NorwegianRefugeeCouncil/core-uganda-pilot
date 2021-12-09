import { DefaultTheme as RNPDefaultTheme } from 'react-native-paper';
import { DefaultTheme as NavigationDefaultTheme } from '@react-navigation/native';

const lightTheme = {
  ...RNPDefaultTheme,
  colors: {
    primary: '#24303E',
    accent: '#FF7602',
    background: '#E5E5E5',
    surface: '#E5E5E5',
    error: '#CE3616',
    text: '#000000',
    onSurface: '#FF7602',
    disabled: '#808080',
    placeholder: '#808080',
    backdrop: '#24303E',
    notification: '#FFF0E6',
    nrc: '#FF7602',
    green: '#47914A',
    blue: '#00ADD0',
    yellow: '#FDC82F',
    red: '#CE3616',
    white: '#FFFFFF',
  },
};

export default lightTheme;

export const darkTheme = {
  ...RNPDefaultTheme,
  ...lightTheme,
  dark: true,
  colors: {
    ...lightTheme.colors,
    primary: '#24303E',
    accent: '#FF7602',
    background: '#24303E',
    surface: '#808080',
    text: '#FFFFFF',
    onSurface: '#FF7602',
    backdrop: '#808080',
    notification: '#808080',
  },
};

export const NavigationTheme = {
  ...NavigationDefaultTheme,
  dark: false,
  colors: {
    ...NavigationDefaultTheme.colors,
    ...lightTheme.colors,
    card: '#24303E',
    text: '#FFFFFF',
    border: '#24303E',
  },
};
