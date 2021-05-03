import { AnchorHTMLAttributes, FunctionComponent, HTMLAttributes } from 'react';
import { addClasses } from '../../utils/utils';

export interface ListGroupProps extends HTMLAttributes<HTMLDivElement | HTMLUListElement> {
  flush?: boolean
  isActionListGroup?: boolean
  numbered?: boolean
}

const ListGroup: FunctionComponent<ListGroupProps> = ({flush, isActionListGroup, numbered, ...props}) => {

  const classes = ['list-group'];
  if (flush) {
    classes.push('list-group-flush');
  }
  if (isActionListGroup) {
    return (<div {...props} className={addClasses(props.className, ...classes)}>{props.children}</div>);
  } else if (numbered) {
    classes.push('list-group-numbered');
    return (<ol  {...props} className={addClasses(props.className, ...classes)}>{props.children}</ol>);
  } else {
    return (<ul  {...props} className={addClasses(props.className, ...classes)}>{props.children}</ul>);
  }
};

export interface ListGroupItemProps {
  active?: boolean
  disabled?: boolean
  isAction?: boolean
}

const ListGroupItem: FunctionComponent<(HTMLAttributes<HTMLLIElement> | AnchorHTMLAttributes<HTMLAnchorElement>) & ListGroupItemProps> = ({active, isAction, disabled, ...props}) => {
  const classes = ['list-group-item'];
  if (active) {
    classes.push('active');
  }
  if (disabled) {
    classes.push('disabled');
  }
  if (isAction) {
    classes.push('list-group-item-action');
  }
  if (isAction) {
    const p = props as AnchorHTMLAttributes<HTMLAnchorElement>;
    return (<a {...p} className={addClasses(props.className, ...classes)}>{props.children}</a>);
  } else {
    const p = props as HTMLAttributes<HTMLLIElement>;
    return (<li {...p} className={addClasses(props.className, ...classes)}>{props.children}</li>);
  }

};


export { ListGroup, ListGroupItem };
