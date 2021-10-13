var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export default function (forms, options) {
    let redirect = options === null || options === void 0 ? void 0 : options.redirectPath;
    if (forms && validateClientSide(forms)) {
        validateServerSide(forms).then(isValid => {
            if (isValid) {
                if (redirect)
                    location.assign(redirect);
            }
        });
    }
}
function validateClientSide(forms) {
    // FIXME I'm really dumb.
    //  For instance, I validate input, select, and textarea elements but not custom form elements.
    let isValid = true;
    for (const form of forms) {
        if (!form.reportValidity()) {
            isValid = false;
        }
        form.classList.add('was-validated');
    }
    return isValid;
}
// validateServerSide initiates and handles server-side validation for document form entities. The handler sends form
// data to the provided endpoints and awaits a validation response object from the server. If validation is received,
// the handler applies it to the concerned DOM elements. If no validation is received, the handler redirects the browser
// to the appropriate location.
function validateServerSide(forms) {
    return __awaiter(this, void 0, void 0, function* () {
        const validations = yield Promise.allSettled(forms.map((formcontrol) => __awaiter(this, void 0, void 0, function* () {
            try {
                const validation = yield validateSubForm(formcontrol);
                removeFormValidation(formcontrol);
                if (validation instanceof Response) {
                    // follow redirect
                    location.assign(validation.url);
                }
                else {
                    formcontrol.classList.add('was-validated');
                    applyServerSideValidation(validation);
                    return Promise.resolve(true);
                }
            }
            catch (e) {
                console.error(e);
            }
            return Promise.resolve(false);
        })));
        return Promise.resolve(validations.every(v => v.status !== 'rejected' && !v.value));
    });
}
function collectSearchParams(formcontrol) {
    const formInputElements = formcontrol.querySelectorAll('[name]');
    const searchParams = new URLSearchParams();
    formInputElements.forEach((formInputElement) => {
        const isCheckboxOrRadio = (ele) => ele.type === 'checkbox' || ele.type === 'radio';
        if (!isCheckboxOrRadio(formInputElement) || (isCheckboxOrRadio(formInputElement) && formInputElement.checked)) {
            searchParams.append(formInputElement.name, formInputElement.value);
        }
        else if (formInputElement instanceof HTMLSelectElement) {
            if (formInputElement.hasAttribute('multiple')) {
                const { children } = formInputElement;
                for (let i = 0; i < children.length; i++) {
                    const child = children.item(i);
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
    return __awaiter(this, void 0, void 0, function* () {
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Accept': 'application/json'
            },
            body: collectSearchParams(formcontrol)
        };
        const response = yield fetch(formcontrol.action, options);
        if (response.ok)
            return response;
        let validation;
        try {
            validation = yield response.json();
        }
        catch (e) {
            console.error(`Fetch error: ${e}`);
            return null;
        }
        return validation;
    });
}
function applyServerSideValidation(validation) {
    for (const { type, name, errors } of validation) {
        if (!errors) {
            continue;
        }
        let domFormControl = document.querySelector(`#${name}, [name=${name}]`);
        if (domFormControl == null) {
            console.error(`element with name "${name} not found`);
            continue;
        }
        if (type === 'taxonomy') {
            domFormControl = domFormControl.parentElement;
        }
        applyFormInputElementValidation(domFormControl, name, errors);
    }
}
function applyClientSideValidation(validation) {
    for (const { formInputElement, errors } of validation) {
        removeFormInputElementValidation(formInputElement);
        applyFormInputElementValidation(formInputElement, formInputElement.name, errors);
    }
}
function applyFormInputElementValidation(element, name, errors) {
    element.classList.add('is-invalid');
    element.setAttribute('aria-describedby', `${name}Feedback`);
    // Append error messages
    let feedback = document.getElementById(`${name}Feedback`);
    if (feedback == null) {
        feedback = appendFeedbackChild(element, name, errors);
    }
    feedback.innerHTML = '';
    for (const error of errors) {
        const p = document.createElement('p');
        p.textContent = error;
        feedback.appendChild(p);
    }
}
function removeFormValidation(formcontrol) {
    const formInputElements = formcontrol.querySelectorAll('[name]');
    formInputElements.forEach(removeFormInputElementValidation);
}
function removeFormInputElementValidation(formInputElement) {
    let target;
    if (formInputElement.classList.contains('taxonomy')) {
        target = formInputElement.parentElement;
    }
    else {
        target = formInputElement;
    }
    target.classList.remove('is-invalid');
    target.removeAttribute('aria-describedby');
    const feedback = document.getElementById(`${formInputElement.name}Feedback`);
    if (feedback != null) {
        feedback.innerHTML = '';
    }
}
function appendFeedbackChild(element, name, errors) {
    const div = document.createElement('div');
    div.id = `${name}Feedback`;
    div.className = 'invalid-feedback';
    return element.insertAdjacentElement('afterend', div);
}
