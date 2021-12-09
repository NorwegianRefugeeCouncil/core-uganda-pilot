import React from 'react';
import { Provider as PaperProvider } from 'react-native-paper';
import * as WebBrowser from 'expo-web-browser';

import theme from './src/constants/theme';
import Router from './src/components/Router';
import { AuthWrapper } from './src/components/AuthWrapper';
import { ErrorBoundary } from './src/components/ErrorBoundary';

WebBrowser.maybeCompleteAuthSession();

export default function App() {
  return (
    <PaperProvider theme={theme}>
      <ErrorBoundary>
        <AuthWrapper>
          <Router />
        </AuthWrapper>
      </ErrorBoundary>
    </PaperProvider>
  );
}
