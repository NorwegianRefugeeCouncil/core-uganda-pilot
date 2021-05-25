import { classNames } from '@core/ui-toolkit/util/utils';
import { Nav, NavProps } from '../nav/nav.component';

export interface NavbarTextProps
  extends React.ComponentPropsWithoutRef<'span'> {}

export const NavbarText: React.FC<NavbarTextProps> = ({
  className: customClass,
  ...rest
}) => {
  const className = classNames('navbar-text', customClass);
  return <span className={className} {...rest} />;
};
