"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
exports.__esModule = true;
var react_1 = require("react");
var react_native_paper_1 = require("react-native-paper");
var styles_1 = require("../../styles");
var react_native_1 = require("react-native");
var clients_1 = require("../../utils/clients");
var react_hook_form_1 = require("react-hook-form");
var Network = require("expo-network");
var FormControl_1 = require("../form/FormControl");
var storage_1 = require("../../utils/storage");
var recordsReducers_1 = require("../../reducers/recordsReducers");
var getEncryptionKey_1 = require("../../utils/getEncryptionKey");
var AddRecordScreen = function (_a) {
    var route = _a.route, dispatch = _a.dispatch;
    var isWeb = react_native_1.Platform.OS === 'web';
    var _b = route.params, formId = _b.formId, recordId = _b.recordId;
    var _c = react_1["default"].useState(true), isLoading = _c[0], setIsLoading = _c[1];
    var _d = react_1["default"].useState(), form = _d[0], setForm = _d[1];
    var _e = react_1["default"].useState(!isWeb), simulateOffline = _e[0], setSimulateOffline = _e[1]; // TODO: for testing, remove
    var _f = react_1["default"].useState(!simulateOffline), isConnected = _f[0], setIsConnected = _f[1];
    var _g = react_1["default"].useState(!isConnected), showSnackbar = _g[0], setShowSnackbar = _g[1];
    var _h = react_1["default"].useState(false), hasLocalData = _h[0], setHasLocalData = _h[1];
    var client = (0, clients_1["default"])();
    var _j = (0, react_hook_form_1.useForm)(), control = _j.control, handleSubmit = _j.handleSubmit, formState = _j.formState, reset = _j.reset;
    react_1["default"].useEffect(function () {
        client.getForm({ id: formId })
            .then(function (data) {
            setForm(data.response);
            setIsLoading(false);
        });
        // TODO add catch
    }, []);
    var onSubmitOffline = function (data) { return __awaiter(void 0, void 0, void 0, function () {
        var key;
        return __generator(this, function (_a) {
            key = (0, getEncryptionKey_1.getEncryptionKey)();
            (0, storage_1.storeEncryptedLocalData)(recordId, key, data)
                .then(function () {
                setHasLocalData(true);
                dispatch({
                    type: recordsReducers_1.RECORD_ACTIONS.ADD_LOCAL_RECORD, payload: {
                        formId: formId,
                        localRecord: recordId
                    }
                });
            })["catch"](function () {
                setHasLocalData(false);
            });
            return [2 /*return*/];
        });
    }); };
    var onSubmit = function (data) {
        if (isConnected || isWeb) {
            client.createRecord({ object: { formId: formId, values: data } });
        }
        else {
            onSubmitOffline(data);
        }
    };
    // check for locally stored data on mobile device
    react_1["default"].useEffect(function () {
        if (!isWeb && recordId) {
            (0, storage_1.getEncryptedLocalData)(recordId)
                .then(function (data) {
                setHasLocalData(!!data);
                reset(data);
            });
        }
    }, [isWeb, recordId]);
    // react to network changes
    react_1["default"].useEffect(function () {
        Network.getNetworkStateAsync()
            .then(function (networkState) {
            // TODO: uncomment, use real network state
            // setIsConnected(networkState.type != NetworkStateType.NONE); // NONE
        })["catch"](function () { return setIsLoading(true); });
    }, [simulateOffline]);
    return (<react_native_1.ScrollView contentContainerStyle={[styles_1.layout.container, styles_1.layout.body, styles_1.common.darkBackground]}>

            <react_native_1.View style={[]}>
                {/* simulate network changes, for testing */}
                {!isWeb && (<react_native_1.View style={{ display: "flex", flexDirection: "row" }}>
                        <react_native_paper_1.Switch value={simulateOffline} onValueChange={function () {
                setSimulateOffline(!simulateOffline);
                setIsConnected(simulateOffline);
                setShowSnackbar(!simulateOffline);
            }}/>
                        <react_native_1.Text> simulate being offline </react_native_1.Text>
                    </react_native_1.View>)}

                {/* upload data collected offline */}
                {hasLocalData && (<react_native_1.View style={{ display: "flex", flexDirection: "column" }}>
                        <react_native_1.Text>
                            There is locally stored data for this individual.
                        </react_native_1.Text>
                    </react_native_1.View>)}
                {hasLocalData && isConnected && (<react_native_1.View style={{ display: "flex", flexDirection: "column" }}>
                        <react_native_1.Text>
                            Do you want to upload it?
                        </react_native_1.Text>
                        <react_native_1.Button title="Submit local data" onPress={handleSubmit(onSubmit)}/>
                    </react_native_1.View>)}
                {!isLoading && (<react_native_1.View style={{ width: '100%' }}>
                        {form === null || form === void 0 ? void 0 : form.fields.map(function (field) {
                return (<FormControl_1["default"] key={field.code} fieldDefinition={field} style={{ width: '100%' }} 
                // value={''} // take value from record
                control={control} name={field.id} errors={formState.errors}/>);
            })}
                        <react_native_1.Button title="Submit" onPress={handleSubmit(onSubmit)}/>
                    </react_native_1.View>)}
            </react_native_1.View>
            <react_native_paper_1.Snackbar visible={showSnackbar} onDismiss={function () { return setShowSnackbar(false); }} action={{
            label: 'Got it',
            onPress: function () { return setShowSnackbar(false); }
        }}>
                No internet connection. Submitted data will be stored locally.
            </react_native_paper_1.Snackbar>
        </react_native_1.ScrollView>);
};
exports["default"] = AddRecordScreen;
