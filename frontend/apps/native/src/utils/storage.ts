import * as SecureStore from 'expo-secure-store';
import AsyncStorage from '@react-native-async-storage/async-storage';
import * as CryptoJS from 'react-native-crypto-js';

export const storeEncryptedLocalData = (recordId: string, key: string, data: Record<any, any>) => {
  return SecureStore.setItemAsync(recordId, key).then(async () => {
    const encryptedData = CryptoJS.AES.encrypt(JSON.stringify(data), key);
    return AsyncStorage.setItem(recordId, encryptedData.toString());
  });
};

export const getEncryptedLocalData = (recordId: string) => {
  return SecureStore.getItemAsync(recordId).then(async (key) => {
    if (key == null) {
      return;
    }

    const data = await AsyncStorage.getItem(recordId);
    if (data == null) {
      return;
    }
    const bytes = await CryptoJS.AES.decrypt(data, key);
    return JSON.parse(bytes.toString(CryptoJS.enc.Utf8));
    // TODO: delete data, once extracted to save space. or only after online submit?
  });
};

export const deleteEncryptedLocalData = (recordId: string) => {
  return SecureStore.getItemAsync(recordId)
    .then((key) => {
      return AsyncStorage.removeItem(recordId);
    })
    .then((key) => {
      return SecureStore.deleteItemAsync(recordId);
    });
};
