import * as React from 'react';
import { classNames, Box, PolymorphicComponentProps } from '@ui-helpers/utils';

export type FormTextOwnProps = {
  muted?: true;
};

export type FormTextProps<
  E extends React.ElementType
> = PolymorphicComponentProps<E, FormTextOwnProps>;

const defaultElement = 'div';

export const FormText = <E extends React.ElementType = typeof defaultElement>({
  muted,
  className: customClass,
  ...rest
}: FormTextProps<E>) => {
  const className = classNames(
    'form-text',
    {
      'text-muted': muted,
    },
    customClass
  );
  return <Box as={defaultElement} className={className} {...rest} />;
};
