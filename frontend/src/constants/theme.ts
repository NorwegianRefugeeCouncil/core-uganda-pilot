import { DefaultTheme } from 'react-native-paper';

const lightTheme = {
  ...DefaultTheme,
  colors: {
    primary: '#24303E',
    accent: '#FF7602',
    background: '#E5E5E5',
    surface: '#FFFFFF',
    error: '#CE3616',
    text: '#000000',
    onSurface: '#000000',
    disabled: '#808080',
    placeholder: '#808080',
    backdrop: '#24303E',
    notification: '#FFF0E6',
    nrc: '#FF7602',
    green: '#47914A',
    blue: '#00ADD0',
    yellow: '#FDC82F',
    red: '#CE3616',
    white: '#FFFFFF'
  }
};

export default lightTheme;

export const darkTheme = {
  ...DefaultTheme,
  ...lightTheme,
  dark: true,
  colors: {
    primary: '#24303E',
    accent: '#FF7602',
    background: '#24303E',
    surface: '#24303E',
    error: '#CE3616',
    text: '#FFFFFF',
    onSurface: '#000000',
    disabled: '#808080',
    placeholder: '#808080',
    backdrop: '#24303E',
    notification: '#FFF0E6',
  }
};
