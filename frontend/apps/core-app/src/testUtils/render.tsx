import * as React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import {
  render as r,
  RenderAPI,
  RenderOptions,
} from '@testing-library/react-native';
import { render as rWeb } from '@testing-library/react';
import { NativeBaseProvider } from 'native-base';
import { theme } from 'core-design-system';

const inset = {
  frame: { x: 0, y: 0, width: 0, height: 0 },
  insets: { top: 0, left: 0, right: 0, bottom: 0 },
};

export const render = (
  component: React.ReactElement<any>,
  options?: RenderOptions | undefined,
): RenderAPI =>
  r(
    <NativeBaseProvider theme={theme} initialWindowMetrics={inset}>
      <NavigationContainer>{component}</NavigationContainer>
    </NativeBaseProvider>,
    options,
  );

export const renderWeb = (
  component: React.ReactElement<any>,
  options?: RenderOptions | undefined,
) =>
  rWeb(
    <NativeBaseProvider theme={theme} initialWindowMetrics={inset}>
      <NavigationContainer>{component}</NavigationContainer>
    </NativeBaseProvider>,
    options,
  );
