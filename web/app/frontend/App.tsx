import React from 'react';
import {Provider as PaperProvider} from 'react-native-paper';
import theme from './src/constants/theme';
import Router from './src/components/Router';
import * as WebBrowser from 'expo-web-browser';
import {AuthWrapper} from "./src/components/AuthWrapper";

WebBrowser.maybeCompleteAuthSession();

export default function App() {
    return (
        <PaperProvider theme={theme}>
            <AuthWrapper>
                <Router/>
            </AuthWrapper>
        </PaperProvider>
    );
}


