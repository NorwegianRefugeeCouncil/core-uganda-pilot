import * as React from 'react';
import { classNames, Box, PolymorphicComponentProps } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/ban-types
export type FormTextOwnProps = {};

export type FormTextProps<
  E extends React.ElementType
> = PolymorphicComponentProps<E, FormTextOwnProps>;

const defaultElement = 'div';

export const FormText = <E extends React.ElementType = typeof defaultElement>({
  className: customClass,
  ...rest
}: FormTextProps<E>) => {
  const className = classNames('form-text', customClass);
  return <Box as={defaultElement} className={className} {...rest} />;
};
