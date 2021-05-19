import * as React from 'react';
import { classNames } from '../../helpers/utils';

export interface CardImgProps extends React.ComponentPropsWithoutRef<'img'> {
  position?: 'top' | 'bottom';
}

type CardImg = React.FC<CardImgProps>;

export const CardImg: CardImg = ({
  position = 'top',
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames(customClass, `card-img-${position}`);
  return <img className={className} {...rest} />;
};
