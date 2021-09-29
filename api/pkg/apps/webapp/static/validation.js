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
exports.validateClientSide = exports.validateServerSide = void 0;
// validateServerSide initiates and handles server-side validation for document form entities. The handler sends form
// data to the provided endpoints and awaits a validation response object from the server. If validation is received,
// the handler applies it to the concerned DOM elements. If no validation is received, the handler redirects the browser
// to the appropriate location.
function validateServerSide(forms, redirectPath) {
    if (redirectPath === void 0) { redirectPath = ''; }
    return __awaiter(this, void 0, void 0, function () {
        var validations, passedValidation, redirect;
        var _this = this;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, Promise.allSettled(forms.map(function (formcontrol) { return __awaiter(_this, void 0, void 0, function () {
                        var validation, e_1;
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0:
                                    _a.trys.push([0, 2, , 3]);
                                    return [4 /*yield*/, validateSubForm(formcontrol)];
                                case 1:
                                    validation = _a.sent();
                                    removeFormValidation(formcontrol);
                                    if (validation != null) {
                                        formcontrol.classList.add('was-validated');
                                        applyServerSideValidation(validation);
                                        return [2 /*return*/, Promise.resolve(true)];
                                    }
                                    return [3 /*break*/, 3];
                                case 2:
                                    e_1 = _a.sent();
                                    console.error(e_1);
                                    return [3 /*break*/, 3];
                                case 3: return [2 /*return*/, Promise.resolve(false)];
                            }
                        });
                    }); }))];
                case 1:
                    validations = _a.sent();
                    passedValidation = validations.every(function (v) { return v.status !== 'rejected' && !v.value; });
                    if (passedValidation) {
                        redirect = redirectPath !== null && redirectPath !== void 0 ? redirectPath : location.origin;
                        location.assign(redirect);
                    }
                    return [2 /*return*/];
            }
        });
    });
}
exports.validateServerSide = validateServerSide;
function validateClientSide(forms) {
    // FIXME I'm really dumb.
    //  For instance, I validate input, select, and textarea elements but not custom form elements.
    var isValid = true;
    for (var _i = 0, forms_1 = forms; _i < forms_1.length; _i++) {
        var form = forms_1[_i];
        if (!form.reportValidity()) {
            isValid = false;
        }
        form.classList.add('was-validated');
    }
    return isValid;
}
exports.validateClientSide = validateClientSide;
function collectSearchParams(formcontrol) {
    var formInputElements = formcontrol.querySelectorAll('[name]');
    var searchParams = new URLSearchParams();
    formInputElements.forEach(function (formInputElement) {
        var isCheckboxOrRadio = function (ele) { return ele.type === 'checkbox' || ele.type === 'radio'; };
        if (!isCheckboxOrRadio(formInputElement) || (isCheckboxOrRadio(formInputElement) && formInputElement.checked)) {
            searchParams.append(formInputElement.name, formInputElement.value);
        }
        else if (formInputElement instanceof HTMLSelectElement) {
            if (formInputElement.hasAttribute('multiple')) {
                var children = formInputElement.children;
                for (var i = 0; i < children.length; i++) {
                    var child = children.item(i);
                    if (child.hasAttribute('selected')) {
                        searchParams.append(formInputElement.name, child.value);
                    }
                }
            }
            else {
                searchParams.append(formInputElement.name, formInputElement.options[formInputElement.selectedIndex].value);
            }
        }
    });
    return searchParams;
}
function validateSubForm(formcontrol) {
    return __awaiter(this, void 0, void 0, function () {
        var options, response, validation, e_2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    options = {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/x-www-form-urlencoded'
                        },
                        body: collectSearchParams(formcontrol)
                    };
                    return [4 /*yield*/, fetch(formcontrol.action, options)];
                case 1:
                    response = _a.sent();
                    if (response.ok)
                        return [2 /*return*/, null];
                    _a.label = 2;
                case 2:
                    _a.trys.push([2, 4, , 5]);
                    return [4 /*yield*/, response.json()];
                case 3:
                    validation = _a.sent();
                    return [3 /*break*/, 5];
                case 4:
                    e_2 = _a.sent();
                    console.error("Fetch error: " + e_2);
                    return [2 /*return*/, null];
                case 5: return [2 /*return*/, validation];
            }
        });
    });
}
function applyServerSideValidation(validation) {
    for (var _i = 0, validation_1 = validation; _i < validation_1.length; _i++) {
        var _a = validation_1[_i], type = _a.type, name_1 = _a.name, errors = _a.errors;
        if (!errors) {
            continue;
        }
        var domFormControl = document.querySelector("#" + name_1 + ", [name=" + name_1 + "]");
        if (domFormControl == null) {
            console.error("element with name \"" + name_1 + " not found");
            continue;
        }
        if (type === 'taxonomy') {
            domFormControl = domFormControl.parentElement;
        }
        applyFormInputElementValidation(domFormControl, name_1, errors);
    }
}
function applyClientSideValidation(validation) {
    for (var _i = 0, validation_2 = validation; _i < validation_2.length; _i++) {
        var _a = validation_2[_i], formInputElement = _a.formInputElement, errors = _a.errors;
        removeFormInputElementValidation(formInputElement);
        applyFormInputElementValidation(formInputElement, formInputElement.name, errors);
    }
}
function applyFormInputElementValidation(element, name, errors) {
    element.classList.add('is-invalid');
    element.setAttribute('aria-describedby', name + "Feedback");
    // Append error messages
    var feedback = document.getElementById(name + "Feedback");
    if (feedback == null) {
        feedback = appendFeedbackChild(element, name, errors);
    }
    feedback.innerHTML = '';
    for (var _i = 0, errors_1 = errors; _i < errors_1.length; _i++) {
        var error = errors_1[_i];
        var p = document.createElement('p');
        p.textContent = error.detail;
        feedback.appendChild(p);
    }
}
function removeFormValidation(formcontrol) {
    var formInputElements = formcontrol.querySelectorAll('[name]');
    formInputElements.forEach(removeFormInputElementValidation);
}
function removeFormInputElementValidation(formInputElement) {
    var target;
    if (formInputElement.classList.contains('taxonomy')) {
        target = formInputElement.parentElement;
    }
    else {
        target = formInputElement;
    }
    target.classList.remove('is-invalid');
    target.removeAttribute('aria-describedby');
    var feedback = document.getElementById(formInputElement.name + "Feedback");
    if (feedback != null) {
        feedback.innerHTML = '';
    }
}
function appendFeedbackChild(element, name, errors) {
    var div = document.createElement('div');
    div.id = name + "Feedback";
    div.className = 'invalid-feedback';
    return element.insertAdjacentElement('afterend', div);
}
