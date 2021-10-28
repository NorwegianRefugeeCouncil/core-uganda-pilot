import {generateSecureRandom} from 'react-native-securerandom';
// import * as Keychain from 'react-native-keychain';
import base64 from 'react-native-base64'

// const ENCRYPTION_KEY = 'UNIQUE_ID'; // TODO: curremt user id?

export const getEncryptionKey = async () => {
    // check for existing credentials
    // const existingCredentials = await Keychain.getGenericPassword();
    // if (existingCredentials) {
    //     return {isFresh: false, key: existingCredentials.password};
    // }
    // generate new credentials based on random string
    const randomBytes = await generateSecureRandom(32);
    return base64.encodeFromByteArray(randomBytes);

    // const hasSetCredentials = await
    //     Keychain.setGenericPassword(ENCRYPTION_KEY, randomBytesString);
    // if (hasSetCredentials) {
    // // } else {
    //     return {isFresh: true, key: randomBytesString};
    // } else {
    //     return {isFresh: false, key: null}
    // }
};
