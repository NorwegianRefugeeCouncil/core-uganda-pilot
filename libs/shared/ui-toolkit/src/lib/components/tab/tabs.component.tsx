import { FunctionComponent, useState } from 'react';
import classNames from 'classnames';
import * as React from 'react';

type TabProps = any & {
  key: number;
  onPointerDown: (e: PointerEvent) => unknown;
  isActive?: boolean;
  isDisabled?: boolean;
}

const Tab: FunctionComponent<TabProps> = (props) => {
  const classes = classNames('nav-item', props.className, {
    active: props.isActive,
    disabled: props.isDisabled
  })
  return (<a className={classes} href={props.href} onPointerDown={props.onPointerDown}>{props.children}</a>)
}

type TabsProps = HTMLUListElement & {

}

const Tabs: FunctionComponent<TabsProps> = ({children, ...props}) => {
  const [activeTab, setActiveTab] = useState(0)

  const handlePointerDown = (key: number) => setActiveTab(key);

  const classes = classNames('nav nav-tabs', props.className, {})
  return (
    <nav className={classes}>
      {
        React.Children.map(children, (child, index) =>
          React.cloneElement(child, )
     }
     {/*<Tab key={index} onPointerDown={() => handlePointerDown(index)} isActive={activeTab === index} isDisabled={child.props.isDisabled}>child</Tab>*/}
    </nav>
  )
}


export { Tab, Tabs }
