type FormInputElement = HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement;

interface ServerSideFormcontrolValidation {
  type: string;
  name: string;
  errors: string[];
}

interface ClientSideFormcontrolValidation {
  formInputElement: FormInputElement;
  errors: string[];
}

interface ValidationOptions {
    redirectPath?: string;
}

export default function (forms: HTMLFormElement | HTMLFormElement[], submitBtn: HTMLButtonElement, options?: ValidationOptions) {
    let redirect = options?.redirectPath;

    const f: HTMLFormElement[] = forms instanceof HTMLFormElement ? [forms] : forms

    submitBtn.addEventListener('click', event => {
        event.preventDefault()
        if (validateClientSide(f)) {
            submitBtn.disabled = true;
            validateServerSide(f).then(isValid => {
                if (isValid) {
                    if (redirect) location.assign(redirect);
                } else {
                    submitBtn.disabled = false;
                }
            });
        }
    })
}

function validateClientSide(forms: HTMLFormElement[]): boolean {
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
async function validateServerSide(forms: HTMLFormElement[]): Promise<boolean> {
    const validations = await Promise.allSettled(forms.map(async (formcontrol) => {
        try {
            const validation = await validateSubForm(formcontrol);

            removeFormValidation(formcontrol);

            if (validation instanceof Response) {
                // follow redirect
                location.assign(validation.url);
            } else {
                formcontrol.classList.add('was-validated');
                applyServerSideValidation(validation);
                return Promise.resolve(true);
            }
        } catch (e) {
            console.error(e);
        }

        return Promise.resolve(false);
    }));

    return Promise.resolve(validations.every(v => v.status !== 'rejected' && !v.value));
}

function collectSearchParams(formcontrol: HTMLFormElement): URLSearchParams {
  const formInputElements = formcontrol.querySelectorAll('[name]');
  const searchParams = new URLSearchParams();

  formInputElements.forEach((formInputElement: FormInputElement) => {
    const isCheckboxOrRadio = (ele): ele is HTMLInputElement => ele.type === 'checkbox' || ele.type === 'radio';
    if (!isCheckboxOrRadio(formInputElement) || (isCheckboxOrRadio(formInputElement) && formInputElement.checked)) {
      searchParams.append(formInputElement.name, formInputElement.value);
    } else if (formInputElement instanceof HTMLSelectElement) {
      if (formInputElement.hasAttribute('multiple')) {
        const { children } = formInputElement;

        for (let i = 0; i < children.length; i++) {
          const child = children.item(i) as HTMLOptionElement;
          if (child.hasAttribute('selected')) {
            searchParams.append(formInputElement.name, child.value);
          }
        }
      } else {
        searchParams.append(formInputElement.name, formInputElement.options[formInputElement.selectedIndex].value);
      }
    }
  });
  return searchParams;
}

async function validateSubForm(formcontrol: HTMLFormElement): Promise<ServerSideFormcontrolValidation[] | Response> {
    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Accept': 'application/json'
        },
        body: collectSearchParams(formcontrol)
    };

    const response = await fetch(formcontrol.action, options);

    if (response.ok) return response;

    let validation;

    try {
        validation = await response.json();
    } catch (e) {
        console.error(`Fetch error: ${e}`);
        return null;
    }
    return validation;
}

function applyServerSideValidation(validation: ServerSideFormcontrolValidation[]) {
  for (const { type, name, errors } of validation) {
    if (!errors) {
      continue;
    }
    let domFormControl = document.querySelector(`#${name}, [name=${name}]`) as FormInputElement | HTMLDivElement;
    if (domFormControl == null) {
      console.error(`element with name "${name} not found`);
      continue;
    }
    if (type === 'taxonomy') {
      domFormControl = domFormControl.parentElement as HTMLDivElement;
    }
    applyFormInputElementValidation(domFormControl, name, errors);
  }
}

function applyClientSideValidation(validation: ClientSideFormcontrolValidation[]) {
  for (const { formInputElement, errors } of validation) {
    removeFormInputElementValidation(formInputElement);
    applyFormInputElementValidation(formInputElement, formInputElement.name, errors);
  }

}

function applyFormInputElementValidation(element: FormInputElement | HTMLDivElement, name: string, errors: string[]) {
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

function removeFormValidation(formcontrol: HTMLFormElement) {
  const formInputElements = formcontrol.querySelectorAll('[name]');
  formInputElements.forEach(removeFormInputElementValidation);
}

function removeFormInputElementValidation(formInputElement: FormInputElement) {
  let target: FormInputElement | HTMLDivElement;

  if (formInputElement.classList.contains('taxonomy')) {
    target = formInputElement.parentElement as HTMLDivElement;
  } else {
    target = formInputElement;
  }

  target.classList.remove('is-invalid');
  target.removeAttribute('aria-describedby');
  const feedback = document.getElementById(`${formInputElement.name}Feedback`);
  if (feedback != null) {
    feedback.innerHTML = '';
  }

}

function appendFeedbackChild(element: FormInputElement | HTMLDivElement, name: string, errors: string[]): HTMLDivElement {
  const div = document.createElement('div');
  div.id = `${name}Feedback`;
  div.className = 'invalid-feedback';
  return element.insertAdjacentElement('afterend', div) as HTMLDivElement;
}

