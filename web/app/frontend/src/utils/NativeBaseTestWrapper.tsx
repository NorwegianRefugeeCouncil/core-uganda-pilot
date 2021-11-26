import React from 'react';
import {NativeBaseProvider} from "native-base";

// For testing purposes: https://docs.nativebase.io/testing
export const testInset = {
    frame: { x: 0, y: 0, width: 0, height: 0 },
    insets: { top: 0, left: 0, right: 0, bottom: 0 },
};

export const NativeBaseTestWrapper = ({children}:{children: React.ReactNode}) => {
    return <NativeBaseProvider initialWindowMetrics={testInset}>
        {children}
    </NativeBaseProvider>
}
