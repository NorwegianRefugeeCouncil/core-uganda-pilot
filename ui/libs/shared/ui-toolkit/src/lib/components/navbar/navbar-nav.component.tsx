import { classNames } from '@core/ui-toolkit/util/utils';
import { Nav, NavProps } from '../nav/nav.component';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface NavbarNavProps extends NavProps {}

export const NavbarNav: React.FC<NavbarNavProps> = ({
  className: customClass,
  ...rest
}) => {
  const className = classNames('navbar-nav', customClass);
  return <ul className={className} {...rest} />;
};
