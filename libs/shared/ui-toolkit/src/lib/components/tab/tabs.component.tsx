import * as React from 'react';
import classNames from 'classnames';
import { ProgressPlugin } from 'webpack';

export interface TabsProps extends React.ComponentPropsWithoutRef<'nav'> {
  align?: 'start' | 'end' | 'center';
  selected?: number;
}

const Tabs: React.FC<TabsProps> = ({
  align = 'start',
  className,
  children,
  ...baseProps
}) => {
  const [activeTab, setActiveTab] = React.useState(0);
  const classes = classNames('nav nav-tabs', className, {});
  return (
    <nav className={classes} {...baseProps}>
      {/* {React.Children.map(children, (child, idx) =>
        {
          const handlePointerDown = () => setActiveTab(idx)
          return React.cloneElement(child, {handlePointerDown})
        }
      )} */}
      /// TODO just can't get this to work, can't find any helpful docs...
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
      className,
      handlePointerDown,
      children,
      ...baseProps
    },
    ref
  ) => {
    const tabLinkClass = classNames('nav-link', className);
    return (
      <a
        ref={ref}
        {...baseProps}
        onPointerDown={handlePointerDown}
        className={tabLinkClass}
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
    const handleClic = () => clickCallBack(key);
    return (
      <li ref={ref} key={key} className={tabClasses}>
        {children}
      </li>
    );
  }
);

export { TabLink, Tab, Tabs };
