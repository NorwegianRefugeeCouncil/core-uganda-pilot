"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const react_1 = __importDefault(require("react"));
const react_native_1 = require("react-native");
const Button = ({ onPress, variant, disabled, text }) => {
    const backgroundColor = {
        primary: 'orange',
        secondary: 'blue'
    };
    return (react_1.default.createElement(react_native_1.Button, { onPress: onPress, color: backgroundColor[variant], title: text, disabled: disabled, style: { borderColor: 'red' } }));
};
exports.default = Button;
