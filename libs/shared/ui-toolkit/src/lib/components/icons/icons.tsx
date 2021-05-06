import { FunctionComponent } from 'react';

const bootstrapIcon: (className: string) => FunctionComponent = className => {
  return props => {
    return <i className={'bi bi-' + className}>{props.children}</i>;
  };
};

export const ThreeDotsIcon = bootstrapIcon('three-dots');
export const PencilFillIcon = bootstrapIcon('pencil-fill');
export const PencilIcon = bootstrapIcon('pencil');
export const PlusIcon = bootstrapIcon('plus');
export const XIcon = bootstrapIcon('x')
