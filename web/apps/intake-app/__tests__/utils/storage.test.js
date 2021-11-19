import {storeEncryptedLocalData} from '../../src/utils/storage';
import AsyncStorage from '@react-native-community/async-storage';
import * as SecureStore from "expo-secure-store";
import * as CryptoJS from "react-native-crypto-js";

jest.mock('@react-native-community/async-storage');
jest.mock('expo-secure-store');
jest.mock('react-native-crypto-js');

describe("utils/storage", () => {
    it("should encrypt data and store it", () => {

        SecureStore.setItemAsync.mockImplementation(() => Promise.resolve(resp));
        const result = storeEncryptedLocalData('recordId', 'key', {id: Symbol('ID')})

    });
});
