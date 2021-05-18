import * as React from 'react';
import classNames from 'classnames';

export interface TabsProps extends React.ComponentPropsWithoutRef<'nav'> {
  align?: 'start' | 'end' | 'center';
  selected?: number;
}

const Tabs: React.FC<TabsProps> = ({
  align = 'start',
  className: customClassName,
  children,
  ...rest
}) => {
  const [activeTab, setActiveTab] = React.useState(0);
  const className = classNames('nav nav-tabs', customClassName, {});
  return (
    <nav className={className} {...rest}>
      {children}
    </nav>
  );
};

interface TabLinkProps extends React.ComponentPropsWithRef<'a'> {
  key: number;
  isDisabled?: boolean;
  handlePointerDown: () => void;
}

const TabLink = React.forwardRef<HTMLAnchorElement, TabLinkProps>(
  (
    {
      isDisabled = false,
      className: customClassName,
      handlePointerDown,
      children,
      ...rest
    },
    ref
  ) => {
    const className = classNames('nav-link', customClassName);
    return (
      <a
        ref={ref}
        {...rest}
        onPointerDown={handlePointerDown}
        className={className}
      >
        {children}
      </a>
    );
  }
);
export interface TabProps extends React.ComponentPropsWithRef<'li'> {
  key: number;
  clickCallBack: (key: number) => void;
  isActive?: boolean;
  isDisabled?: boolean;
}

const Tab = React.forwardRef<HTMLLIElement, TabProps>(
  (
    {
      isActive = false,
      isDisabled = false,
      className,
      children,
      key = 0,
      clickCallBack,
      ...baseProps
    }: TabProps,
    ref
  ) => {
    const tabClasses = classNames('nav-item', className, {
      active: isActive,
      disabled: isDisabled,
    });
    const handlePointerDown = () => clickCallBack(key);
    return (
      <li ref={ref} key={key} className={tabClasses}>
        {children}
      </li>
    );
  }
);

export { TabLink, Tab, Tabs };
