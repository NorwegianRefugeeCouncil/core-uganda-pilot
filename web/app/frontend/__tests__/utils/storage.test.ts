import {storeEncryptedLocalData} from '../../src/utils/storage';
import * as SecureStore from "expo-secure-store";
import {mocked} from "jest-mock"

jest.mock('@react-native-async-storage/async-storage');
jest.mock('expo-secure-store');
jest.mock('react-native-crypto-js');

describe("utils/storage", () => {
    it("should encrypt data and store it", () => {

        mocked(SecureStore).setItemAsync.mockImplementation(() => Promise.resolve());
        const result = storeEncryptedLocalData('recordId', 'key', {id: Symbol('ID')})

    });
});
