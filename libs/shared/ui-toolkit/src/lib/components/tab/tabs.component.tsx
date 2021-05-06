import * as React from 'react';
import classNames from 'classnames';

export type TabProps = React.ComponentPropsWithRef<'a'> & {
  key: number;
  onPointerDown: (e: React.PointerEvent<HTMLAnchorElement>) => void;
  isActive?: boolean;
  isDisabled?: boolean;
}

const Tab: React.FC<TabProps> = (props) => {
  const classes = classNames('nav-item', props.className, {
    active: props.isActive,
    disabled: props.isDisabled
  });
  return (<a className={classes} href={props.href} onPointerDown={props.onPointerDown}>{props.children}</a>);
};

export type TabsProps = {
} // TODO

const Tabs: React.FC<> = ({ children, ...props }) => {
  const [activeTab, setActiveTab] = React.useState(0);

  const handlePointerDown = (key: number) => setActiveTab(key);

  const classes = classNames('nav nav-tabs', props.className, {});
  return (
    <nav className={classes}>
      {
        React.Children.map(children, (child, index) => {
          React.cloneElement(child);
        })
      }
      {/*<Tab key={index} onPointerDown={() => handlePointerDown(index)} isActive={activeTab === index} isDisabled={child.props.isDisabled}>child</Tab>*/}
    </nav>
  );
};


export { Tab, Tabs };
