// eslint-disable-next-line no-unused-vars,@typescript-eslint/no-unused-vars
import React from 'react';
import { NativeBaseProvider } from 'native-base';
import { theme } from 'core-design-system';

import StorybookUIRoot from './storybook';

export default function App() {
  return (
    <NativeBaseProvider theme={theme}>
      <StorybookUIRoot />
    </NativeBaseProvider>
  );
}
