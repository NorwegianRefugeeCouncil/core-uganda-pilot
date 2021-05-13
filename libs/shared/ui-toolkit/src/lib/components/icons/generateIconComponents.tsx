import { FunctionComponent } from 'react';
import iconNames from './iconnames';

const Icons = {};

const bootstrapIcon: (className: string) => FunctionComponent = (className) => {
  return (props) => {
    return <i className={'bi bi-' + className}>{props.children}</i>;
  };
};

const kebabToCamel = (kebab: string) => {
  const words = kebab.match(/\w+/g);
  return words
    .map((word) => `${word.slice(0, 1).toUpperCase()}${word.slice(1)}`)
    .join('');
};

iconNames.forEach((kebabName) => {
  const camelName = kebabToCamel(kebabName);
  Icons[camelName] = bootstrapIcon(kebabName);
});

export default Icons;
