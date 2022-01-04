// eslint-disable-next-line no-unused-vars,@typescript-eslint/no-unused-vars
import React from 'react';
import { NativeBaseProvider } from 'native-base';
import { theme } from 'core-design-system';
import {
  useFonts,
  // eslint-disable-next-line camelcase
  Roboto_400Regular,
  // eslint-disable-next-line camelcase
  Roboto_400Regular_Italic,
  // eslint-disable-next-line camelcase
  Roboto_700Bold,
} from '@expo-google-fonts/roboto';

import StorybookUIRoot from './storybook';

export default function App() {
  const [fontsLoaded] = useFonts({
    Roboto_400Regular,
    Roboto_400Regular_Italic,
    Roboto_700Bold,
  });

  return (
    fontsLoaded && (
      <NativeBaseProvider theme={theme}>
        <StorybookUIRoot />
      </NativeBaseProvider>
    )
  );
}
