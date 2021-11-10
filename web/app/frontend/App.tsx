import React from 'react';
import { Provider as PaperProvider } from 'react-native-paper';

import Router from './src/components/Router';
import theme from './src/constants/theme';

export default function App() {
    return (
        <PaperProvider theme={theme}>
            <Router />
        </PaperProvider>
    );
}
