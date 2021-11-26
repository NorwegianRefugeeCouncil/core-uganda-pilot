import React from 'react';
import {NativeBaseProvider} from 'native-base';

import Router from './src/components/Router';

export default function App({inset}) {
    return (
        <NativeBaseProvider initialWindowMetrics={inset}>
            <Router/>
        </NativeBaseProvider>
    );
}
