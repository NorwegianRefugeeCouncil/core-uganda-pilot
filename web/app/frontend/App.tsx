import React from 'react';
import { Platform } from 'react-native';
import { IAMClient } from 'core-js-api-client';
import { Provider as PaperProvider } from 'react-native-paper';
import Constants from 'expo-constants';
import theme from './src/constants/theme';
import Router from './src/components/Router';

export default function App() {
	get();
	return (
		<PaperProvider theme={theme}>
      <Router />
		</PaperProvider>
	);
}

export const host = (Platform.OS === 'web' ? Constants.manifest?.extra?.server_default_hostname : Constants.manifest?.extra?.server_hostname);
const iamClient = new IAMClient(host, 'http', {});

function get() {
	iamClient.Parties().Get()('c529d679-3bb6-4a20-8f06-c096f4d9adc1');
}
