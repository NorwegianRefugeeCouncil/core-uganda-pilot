//
// Form handling
//

function Validator(forms: { formID: string, path: string }[], submit?: HTMLButtonElement) {
  const submitBtn = submit ? submit : document.querySelector('button[type=submit]') as HTMLButtonElement;
  submitBtn.onclick = async event => {
    // Do validation
    const validations = await Promise.allSettled<Promise<boolean>[]>(forms.map(async ({ formID, path }) => {
      try {
        const validation = await validateSubForm(formID, path);
        removeValidation(formID);
        if (validation != null) {
          applyValidation(validation);
          return Promise.resolve(true);
        }
      } catch (e) {
        console.error(e);
      }
      return Promise.resolve(false);
    }));

    if (validations.every(v => !v)) {
      if (!location.pathname.includes('new')) {
        location.reload();
      } else {
        // FIXME!! where should this point
        location.assign(location.origin);
      }
    }
  };
  return {};
}

type InputElement = HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement;

function collectSearchParams(formID) {
  const fields = document.querySelectorAll(`#${formID} [name]`);
  const searchParams = new URLSearchParams();
  fields.forEach((field: InputElement) => {
    const isCheckboxOrRadio = (field): field is HTMLInputElement => field.type === 'checkbox' || field.type === 'radio';
    if (!isCheckboxOrRadio || (isCheckboxOrRadio(field) && field.checked)) {
      searchParams.append(field.name, field.value);
    } else if (field instanceof HTMLSelectElement) {
      if (field.hasAttribute('multiple')) {
        const children = field.children;
        for (let i = 0; i < children.length; i++) {
          const child = children.item(i) as HTMLOptionElement;
          if (child.hasAttribute('selected')) {
            searchParams.append(field.name, child.value);
          }
        }
      } else {
        searchParams.append(field.name, field.options[field.selectedIndex].value);
      }
    }
  });
  return searchParams;
}

async function validateSubForm(formID, path) {
  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded'
    },
    body: collectSearchParams(formID)
  };

  const response = await fetch(path, options);
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

function applyValidation(validation) {
  for (const element of validation) {
    let { name } = element.attributes;
    if (!element.errors) {
      continue;
    }
    // Get DOM form element
    let el = document.querySelector(`#${name}`);
    if (element.type === 'taxonomyinput') {
      el = el.parentElement;
    }
    // Apply attributes
    el.classList.add('is-invalid');
    el.setAttribute('aria-describedby', `${name}Feedback`);
    // Get feedback div
    const feedback = document.getElementById(`${name}Feedback`);
    if (feedback != null) {
      feedback.innerHTML = '';
      // Populate attributes and error messages
      for (const error of element.errors) {
        const p = document.createElement('p');
        p.textContent = error.detail;
        feedback.appendChild(p);
      }
    }
  }
}

function removeValidation(formID: string) {
  const fields = document.querySelectorAll(`#${formID} [name]`);
  fields.forEach((field: InputElement) => {
    const name = field.name;
    let receiver: InputElement | HTMLDivElement;
    if (field.classList.contains('taxonomyinput')) {
      receiver = field.parentElement as HTMLDivElement;
    } else {
      receiver = field;
    }
    receiver.classList.remove('is-invalid');
    receiver.removeAttribute('aria-describedby');
    const feedback = document.getElementById(`${name}Feedback`);
    if (feedback != null) {
      feedback.innerHTML = '';
    }
  });
}



