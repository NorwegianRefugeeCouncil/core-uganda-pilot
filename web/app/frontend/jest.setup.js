import 'react-native-gesture-handler/jestSetup';
import mockAsyncStorage from '@react-native-async-storage/async-storage/jest/async-storage-mock';

jest.mock('@react-native-async-storage/async-storage', () => mockAsyncStorage);

jest.mock('expo-secure-store');

jest.mock('react-native-crypto-js', () => {
    const CryptoJS = jest.requireActual('react-native-crypto-js');
    return {
        ...CryptoJS,
        AES: {
            ...CryptoJS.AES,
            encrypt: jest.fn(CryptoJS.AES.encrypt),
            decrypt: jest.fn(CryptoJS.AES.decrypt),
        },
    };
});

jest.mock('@react-native-community/datetimepicker', function () {
    const mockComponent = jest.requireActual('react-native/jest/mockComponent');
    return mockComponent('@react-native-community/datetimepicker');
});

// https://reactnavigation.org/docs/testing
jest.mock('react-native-reanimated', () => {
    const Reanimated = jest.requireActual('react-native-reanimated/mock');

    // The mock for `call` immediately calls the callback which is incorrect
    // So we override it with a no-op
    Reanimated.default.call = () => {};

    return Reanimated;
});

// Silence the warning: Animated: `useNativeDriver` is not supported because the native animated module is missing
jest.mock('react-native/Libraries/Animated/NativeAnimatedHelper');

// Mock being on android
jest.mock('react-native/Libraries/Utilities/Platform', () => {
    const Platform = jest.requireActual(
        'react-native/Libraries/Utilities/Platform'
    );
    Platform.OS = 'android';
    return Platform;
});
