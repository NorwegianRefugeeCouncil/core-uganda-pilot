import * as React from 'react';
import classNames from 'classnames';

export interface TabsProps extends React.ComponentProps<typeof Tab> {
  align?: 'start';
  end;
  center;
}

const Tabs: React.FC<TabsProps> = ({
  align = 'start',
  className,
  children,
  ...baseProps
}) => {
  const [activeTab, setActiveTab] = React.useState(0);

  const handlePointerDown = (key: number) => setActiveTab(key);
  if (typeof children !== typeof Tab) throw Error();
  const classes = classNames('nav nav-tabs', className, {});
  return <nav className={classes}>{children}</nav>;
};

export { Tab, Tabs };
export interface TabProps
  extends React.ComponentProps<'li'>,
    React.ComponentProps<'a'> {
  key: number;
  onPointerDown: (e: React.PointerEvent<HTMLAnchorElement>) => void;
  isActive?: boolean;
  isDisabled?: boolean;
}
// TODO figure out how to properly type this shit

const Tab = React.forwardRef<HTMLLIElement, TabProps>(
  (
    {
      isActive = false,
      isDisabled = false,
      onPointerDown = () => {},
      className,
      children,
      ...baseProps
    }: TabProps,
    ref
  ) => {
    const classes = classNames('nav-link', className, {
      active: isActive,
      disabled: isDisabled,
    });
    return (
      <li ref={ref}>
        <a className={classes} onPointerDown={onPointerDown} {...baseProps}>
          {children}
        </a>
      </li>
    );
  }
);
