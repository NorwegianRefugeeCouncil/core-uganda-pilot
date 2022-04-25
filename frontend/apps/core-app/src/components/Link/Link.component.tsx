import * as React from 'react';
import { ILinkProps, Link as L } from 'native-base';
import { NavigationAction, useLinkProps } from '@react-navigation/native';
import { To } from '@react-navigation/native/lib/typescript/src/useLinkTo';

type Props = {
  to: To<any>;
  action?: NavigationAction;
} & ILinkProps;

export const LinkComponent: React.FC<Props> = ({
  to,
  action,
  children,
  ...otherProps
}) => {
  const { href, onPress } = useLinkProps({ to, action });

  return (
    <L href={href} onPress={onPress} {...otherProps}>
      {children}
    </L>
  );
};
