import 'expo-dev-client'
import React from 'react';
import {Provider as PaperProvider} from 'react-native-paper';
import theme from './src/constants/theme';
import Router from './src/components/Router';

export default function App() {
    return (
        <PaperProvider theme={theme}>
            <Router/>
        </PaperProvider>
    );
}
