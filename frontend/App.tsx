import React from 'react';
import { Platform } from 'react-native';
import { CMSClient } from 'core-js-api-client';
import { Provider as PaperProvider } from 'react-native-paper';
import Constants from 'expo-constants';
import theme from './src/constants/theme';
import Layout from './src/components/Layout';

export default function App() {
	get();
	return (
		<PaperProvider theme={theme}>
			<Layout />
		</PaperProvider>
	);
}

const host = (Platform.OS === 'web' ? Constants.manifest?.extra?.server_default_hostname : Constants.manifest?.extra?.server_hostname);
const cmsClient = new CMSClient(host, 'http');

function get() {
	cmsClient.Cases().Get('dba43642-8093-4685-a197-f8848d4cbaaa').subscribe();
}

