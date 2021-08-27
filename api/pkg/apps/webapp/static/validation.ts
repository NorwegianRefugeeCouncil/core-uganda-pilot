type FormInputElement = HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement;

interface ServerSideFormElementValidation {
  type: string;
  attributes: {
    name: string
  };
  errors: ValidationError[];
}

interface ClientSideFormElementValidation {
  formInputElement: FormInputElement;
  errors: ValidationError[];
}

interface ValidationError {
  detail: string;
}

// validateServerSide initiates and handles server-side validation for document form entities. The handler sends form
// data to the provided endpoints and awaits a validation response object from the server. If validation is received,
// the handler applies it to the concerned DOM elements. If no validation is received, the handler redirects the browser
// to the appropriate location.
export async function validateServerSide(forms: HTMLFormElement[], redirectPath = '') {
  const validations = await Promise.allSettled(forms.map(async (formElement) => {
    try {
      const validation = await validateSubForm(formElement);
      removeFormValidation(formElement);
      if (validation != null) {
        formElement.classList.add('was-validated');
        applyServerSideValidation(validation);
        return Promise.resolve(true);
      }
    } catch (e) {
      console.error(e);
    }
    return Promise.resolve(false);
  }));

  const passedValidation = validations.every(v => v.status !== 'rejected' && !v.value);
  if (passedValidation) {
    const redirect = redirectPath ?? location.origin;
    location.assign(redirect);
  }
}

export function validateClientSide(forms: HTMLFormElement[]): boolean {
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

function collectSearchParams(formElement: HTMLFormElement): URLSearchParams {
  const formInputElements = formElement.querySelectorAll('[name]');
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

async function validateSubForm(formElement: HTMLFormElement): Promise<ServerSideFormElementValidation[]> {
  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    body: collectSearchParams(formElement)
  };

  const response = await fetch(formElement.action, options);

  if (response.ok) return null;

  let validation;

  try {
    validation = await response.json();
  } catch (e) {
    console.error(`Fetch error: ${e}`);
    return null;
  }
  return validation;
}

function applyServerSideValidation(validation: ServerSideFormElementValidation[]) {
  for (const { type, attributes: { name }, errors } of validation) {
    if (!errors) {
      continue;
    }
    let domFormElement = document.querySelector(`#${name}`) as FormInputElement | HTMLDivElement;
    if (type === 'taxonomyinput') {
      domFormElement = domFormElement.parentElement as HTMLDivElement;
    }
    applyFormInputElementValidation(domFormElement, name, errors);
  }
}

function applyClientSideValidation(validation: ClientSideFormElementValidation[]) {
  for (const { formInputElement, errors } of validation) {
    removeFormInputElementValidation(formInputElement);
    applyFormInputElementValidation(formInputElement, formInputElement.name, errors);
  }

}

function applyFormInputElementValidation(element: FormInputElement | HTMLDivElement, name: string, errors: ValidationError[]) {
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
    p.textContent = error.detail;
    feedback.appendChild(p);
  }
}

function removeFormValidation(formElement: HTMLFormElement) {
  const formInputElements = formElement.querySelectorAll('[name]');
  formInputElements.forEach(removeFormInputElementValidation);
}

function removeFormInputElementValidation(formInputElement: FormInputElement) {
  let target: FormInputElement | HTMLDivElement;

  if (formInputElement.classList.contains('taxonomyinput')) {
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

function appendFeedbackChild(element: FormInputElement | HTMLDivElement, name: string, errors: ValidationError[]): HTMLDivElement {
  const div = document.createElement('div');
  div.id = `${name}Feedback`;
  div.className = 'invalid-feedback';
  return element.insertAdjacentElement('afterend', div) as HTMLDivElement;
}

