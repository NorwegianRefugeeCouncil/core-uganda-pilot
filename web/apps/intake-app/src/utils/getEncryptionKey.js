"use strict";
exports.__esModule = true;
exports.getEncryptionKey = void 0;
require("react-native-get-random-values");
var getEncryptionKey = function () {
    var array = new Uint32Array(10);
    return crypto.getRandomValues(array).toString();
};
exports.getEncryptionKey = getEncryptionKey;
