import React from 'react';
import { IAMClient } from 'core-js-api-client';
import { Provider as PaperProvider } from 'react-native-paper';
import theme from './src/constants/theme';
import Router from './src/components/Router';
import { Subject } from 'rxjs';
import { Platform } from 'react-native';
import Constants from 'expo-constants';

const host = Platform.OS === 'web' ? Constants.manifest?.extra?.server_default_hostname : Constants.manifest?.extra?.server_hostname;
console.log('host: ', host);

const subject = new Subject();

const iamClient = new IAMClient('http', host, { 'X-Authenticated-User-Subject': ['test@user.email'] });

function get() {
	subject.pipe(iamClient.Parties().Get()).subscribe(console.log);
	subject.next('c529d679-3bb6-4a20-8f06-c096f4d9adc1');
}

export default function App() {
	get();
	return (
		<PaperProvider theme={theme}>
			<Router host={host} />
		</PaperProvider>
	);
}


