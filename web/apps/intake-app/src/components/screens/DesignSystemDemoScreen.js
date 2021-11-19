"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var styles_1 = require("../../styles");
var react_native_1 = require("react-native");
// import {Button} from 'core-design-system'
var DesignSystemDemoScreen = function () {
    return (<react_native_1.View style={styles_1.layout.body}>
            <react_native_paper_1.Title>Design System Demo</react_native_paper_1.Title>

            {/*<Button onPress={() => console.log('integrated design system')}>*/}
            {/*    <Text>*/}
            {/*        button*/}
            {/*    </Text>*/}
            {/*</Button>*/}
        </react_native_1.View>);
};
exports["default"] = DesignSystemDemoScreen;
