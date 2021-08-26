type FormInputElement = HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement;

interface ValidationFormElement {
  type: string;
  attributes: {
    name: string
  };
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
      removeValidation(formElement);
      if (validation != null) {
        applyValidation(validation);
        return Promise.resolve(true);
      }
    } catch (e) {
      console.error(e);
    }
    return Promise.resolve(false);
  }));

  const passedValidation = validations.every(v => v.status !== 'rejected' && !v.value);
  if (passedValidation) {
    if (!location.pathname.includes('new')) {
      location.reload();
    } else {
      const redirect = location.origin + redirectPath;
      location.assign(redirect);
    }
  }
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

async function validateSubForm(formElement: HTMLFormElement): Promise<ValidationFormElement[]> {
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

function applyValidation(validation: ValidationFormElement[]) {
  for (const { type, attributes: { name }, errors } of validation) {
    if (!errors) {
      continue;
    }
    let domFormElement = document.querySelector(`#${name}`);
    if (type === 'taxonomyinput') {
      domFormElement = domFormElement.parentElement;
    }

    // Apply attributes
    domFormElement.classList.add('is-invalid');
    domFormElement.setAttribute('aria-describedby', `${name}Feedback`);
    const feedback = document.getElementById(`${name}Feedback`);
    // Append error messages
    if (feedback != null) {
      feedback.innerHTML = '';
      for (const error of errors) {
        const p = document.createElement('p');
        p.textContent = error.detail;
        feedback.appendChild(p);
      }
    }
  }
}

function removeValidation(formElement: HTMLFormElement) {
  const formElements = formElement.querySelectorAll('[name]');

  formElements.forEach((formElement: FormInputElement) => {
    let target: FormInputElement | HTMLDivElement;

    if (formElement.classList.contains('taxonomyinput')) {
      target = formElement.parentElement as HTMLDivElement;
    } else {
      target = formElement;
    }

    target.classList.remove('is-invalid');
    target.removeAttribute('aria-describedby');
    const feedback = document.getElementById(`${formElement.name}Feedback`);
    if (feedback != null) {
      feedback.innerHTML = '';
    }
  });
}



