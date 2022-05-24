import * as React from 'react';
import * as WebBrowser from 'expo-web-browser';
import { NativeBaseProvider } from 'native-base';
import { theme, tokens } from 'core-design-system';
import { NavigationContainer, DefaultTheme } from '@react-navigation/native';
import {
  // eslint-disable-next-line camelcase
  Roboto_400Regular,
  // eslint-disable-next-line camelcase
  Roboto_400Regular_Italic,
  // eslint-disable-next-line camelcase
  Roboto_500Medium,
  // eslint-disable-next-line camelcase
  Roboto_700Bold,
  useFonts,
} from '@expo-google-fonts/roboto';

import { AuthWrapper } from './src/components/AuthWrapper';
import { adminClient, formsClient } from './src/clients';
import { RootNavigator } from './src/navigation';
import { linkingConfig } from './src/navigation/linking.config';

const navTheme = {
  ...DefaultTheme,
  colors: {
    ...DefaultTheme.colors,
    background: tokens.colors.white,
  },
};

WebBrowser.maybeCompleteAuthSession();

const App: React.FC = () => {
  const [fontsLoaded] = useFonts({
    Roboto_400Regular,
    Roboto_400Regular_Italic,
    Roboto_500Medium,
    Roboto_700Bold,
  });

  if (!fontsLoaded) return null;

  const handleTokenChange = (token: string) => {
    formsClient.setToken(token);
    adminClient.setToken(token);
  };

  return (
    <NativeBaseProvider theme={theme}>
      <NavigationContainer linking={linkingConfig} theme={navTheme}>
        <AuthWrapper onTokenChange={handleTokenChange}>
          <RootNavigator />
        </AuthWrapper>
      </NavigationContainer>
    </NativeBaseProvider>
  );
};

export default App;
