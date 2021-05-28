import {
  FormDefinition,
  FormDefinitionNames,
  FormDefinitionSpec,
  FormDefinitionValidation,
  FormDefinitionVersion,
  FormElement,
  FormRoot,
  TranslatedString,
  TranslatedStrings
} from './index';
import { Duplicate, ErrorList, Required } from '../../../field/error';
import { newPath, Path } from '../../../field';
import { validateTypeMeta } from '../../meta/v1/validation';


export const validateFormDefinition = (formDefinition: FormDefinition): ErrorList => {
  const errorList: ErrorList = [];
  errorList.push(...validateTypeMeta(formDefinition, newPath('')));
  errorList.push(...validateFormDefinitionSpec(formDefinition.spec, newPath('spec')));
  return errorList;
};

export const validateFormDefinitionSpec = (spec: FormDefinitionSpec, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!spec) {
    errorList.push(Required(fieldPath, 'Form Definition spec is required'));
    return errorList;
  }

  if (!spec.group) {
    errorList.push(Required(fieldPath.child('group'), 'group is required'));
  }

  errorList.push(...validateFormDefinitionNames(spec.names, fieldPath.child('names')));
  errorList.push(...validateFormDefinitionVersions(spec.versions, fieldPath.child('versions')));

  return errorList;
};

export const validateFormDefinitionNames = (names: FormDefinitionNames, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!names.singular) {
    errorList.push(Required(fieldPath.child('singular'), 'group is required'));
  }

  if (!names.plural) {
    errorList.push(Required(fieldPath.child('plural'), 'group is required'));
  }

  if (!names.kind) {
    errorList.push(Required(fieldPath.child('kind'), 'group is required'));
  }

  return errorList;
};

export const validateFormDefinitionVersions = (versions: FormDefinitionVersion[], fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!versions || versions.length === 0) {
    errorList.push(Required(fieldPath, 'must have at least one version'));
    return errorList;
  }

  for (let i = 0; i < versions.length; i++) {
    errorList.push(...validateFormDefinitionVersion(versions[i], fieldPath.index(i)));
  }

  return errorList;
};


export const validateFormDefinitionVersion = (version: FormDefinitionVersion, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!version) {
    errorList.push(Required(fieldPath, 'version cannot be null'));
    return errorList;
  }

  if (!version.name) {
    errorList.push(Required(fieldPath.child('name'), 'version name is required'));
  }

  errorList.push(...validateFormDefinitionSchema(version.schema, fieldPath.child('schema')));

  return errorList;
};

export const validateFormDefinitionSchema = (schema: FormDefinitionValidation, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];
  if (!schema) {
    errorList.push(Required(fieldPath, 'Version schema is required'));
    return errorList;
  }
  errorList.push(...validateFormRoot(schema.formSchema, fieldPath.child('formSchema')));
  return errorList;
};


export const validateFormRoot = (schema: FormRoot, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];
  if (!schema) {
    errorList.push(Required(fieldPath, 'form schema is required'));
    return errorList;
  }

  errorList.push(...validateRoot(schema.root, fieldPath.child('root')));

  return errorList;
};

export const validateRoot = (root: FormElement, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!root) {
    errorList.push(Required(fieldPath, 'root is required'));
    return errorList;
  }

  errorList.push(...validateFormElement(root, true, fieldPath));

  return errorList;
};

export const validateFormElement = (element: FormElement, isRoot: boolean, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!element) {
    errorList.push(Required(fieldPath, 'element is required'));
    return errorList;
  }

  if (!element.type) {
    errorList.push(Required(fieldPath.child('type'), 'type is required'));
  }

  if (element.type !== 'section') {

    if (!element.key) {
      errorList.push(Required(fieldPath.child('key'), 'key is required'));
    }

    errorList.push(...validateTranslatedStrings(element.label, fieldPath.child('label')));
    if (element.tooltip) {
      errorList.push(...validateTranslatedStrings(element.tooltip, fieldPath.child('tooltip')));
    }
  } else {
    if (!element.children || element.children.length === 0) {
      errorList.push(Required(fieldPath.child('children'), 'section elements must have at least one child'));
    }
  }

  if (element.children) {
    for (let i = 0; i < element.children.length; i++) {
      const child = element.children[i];
      errorList.push(...validateFormElement(child, false, fieldPath.child('children').index(i)));
    }
  }

  return errorList;
};

export const validateTranslatedStrings = (strs: TranslatedStrings, fieldPath: Path): ErrorList => {

  const errorList: ErrorList = [];

  if (!strs) {
    errorList.push(Required(fieldPath, 'translations are required'));
    return errorList;
  }

  if (strs.length === 0) {
    errorList.push(Required(fieldPath, 'must have at least one translation'));
  }

  // assert unique locales
  const seenLocales: { [locale: string]: boolean } = {};
  for (let i = 0; i < strs.length; i++) {
    const translation = strs[i];

    errorList.push(...validateTranslatedString(translation, fieldPath.index(i)));

    if (translation) {
      if (seenLocales[translation.locale]) {
        errorList.push(Duplicate(fieldPath.index(i), 'duplicate locale'));
      }
      seenLocales[translation.locale] = true;
    }
  }


  return errorList;

};

export const validateTranslatedString = (str: TranslatedString, fieldPath: Path): ErrorList => {
  const errorList: ErrorList = [];

  if (!str) {
    errorList.push(Required(fieldPath, 'translation cannot be null'));
    return errorList;
  }

  if (!str.locale) {
    errorList.push(Required(fieldPath.child('locale'), 'locale is required'));
  }
  if (!str.value) {
    errorList.push(Required(fieldPath.child('value'), 'value is required'));
  }

  return errorList;
};


























