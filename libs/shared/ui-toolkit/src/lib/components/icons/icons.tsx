import * as React from 'react';

const bootstrapIcon: (className: string) => React.FC<React.ComponentPropsWithoutRef<'i'>> = className => {
  return ({ children, ...props }) => {
    return <i {...props} className={'bi bi-' + className}>{children}</i>;
  };
};

export const ThreeDotsIcon = bootstrapIcon('three-dots');
export const PencilFillIcon = bootstrapIcon('pencil-fill');
export const PencilIcon = bootstrapIcon('pencil');
export const PlusIcon = bootstrapIcon('plus');
export const XIcon = bootstrapIcon('x')
