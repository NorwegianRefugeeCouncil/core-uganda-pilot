"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var theme_1 = require("./src/constants/theme");
var Router_1 = require("./src/components/Router");
var WebBrowser = require("expo-web-browser");
var AuthWrapper_1 = require("./src/components/AuthWrapper");
WebBrowser.maybeCompleteAuthSession();
function App() {
    return (<react_native_paper_1.Provider theme={theme_1["default"]}>
            <AuthWrapper_1.AuthWrapper>
                <Router_1["default"] />
            </AuthWrapper_1.AuthWrapper>
        </react_native_paper_1.Provider>);
}
exports["default"] = App;
