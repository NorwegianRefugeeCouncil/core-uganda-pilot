import * as React from 'react';
import classNames from 'classnames';
import { ButtonProps } from '../button/button.component';

export interface CloseButtonProps extends ButtonProps {
  size?: 'sm' | 'lg';
  white?: boolean;
}

const CloseButton: React.FC<CloseButtonProps> = ({
  white = false,
  size,
  className: customClass,
  ...rest
}) => {
  const className = classNames('btn-close', customClass, {
    [`btn-${size}`]: size != null,
    'btn-close-white': white,
  });
  return <button type="button" className={className} {...rest} />;
};

export default CloseButton;
