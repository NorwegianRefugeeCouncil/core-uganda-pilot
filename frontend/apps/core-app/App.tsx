import * as React from 'react';
import * as WebBrowser from 'expo-web-browser';
import { NativeBaseProvider } from 'native-base';
import { theme } from 'core-design-system';
import { NavigationContainer } from '@react-navigation/native';

import { AuthWrapper } from './src/components/AuthWrapper';
import { formsClient } from './src/clients/formsClient';
import { RootNavigator } from './src/navigators';

WebBrowser.maybeCompleteAuthSession();

const App: React.FC = () => {
  return (
    <NativeBaseProvider theme={theme}>
      <NavigationContainer>
        <AuthWrapper onTokenChange={formsClient.setAuth}>
          <RootNavigator />
        </AuthWrapper>
      </NavigationContainer>
    </NativeBaseProvider>
  );
};

export default App;
