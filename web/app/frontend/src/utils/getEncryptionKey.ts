import base64 from 'react-native-base64'
import 'react-native-get-random-values'

// const ENCRYPTION_KEY = 'UNIQUE_ID'; // TODO: curremt user id?

export const getEncryptionKey = () => {
    console.log('create key')
    var array = new Uint32Array(10);
    return base64.encode(crypto.getRandomValues(array).toString())
};
