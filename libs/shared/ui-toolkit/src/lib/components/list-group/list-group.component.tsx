import { AnchorHTMLAttributes, FC, HTMLAttributes } from 'react';
import classNames from 'classnames';
export interface ListGroupProps
  extends HTMLAttributes<HTMLDivElement | HTMLUListElement> {
  flush?: boolean;
  isActionListGroup?: boolean;
  numbered?: boolean;
}

const ListGroup: FC<ListGroupProps> = ({
  flush,
  isActionListGroup,
  numbered,
  ...props
}) => {
  const classes = (obj) => classNames(props.className, 'list-group', obj);
  if (isActionListGroup) {
    return (
      <div {...props} className={props.className}>
        {props.children}
      </div>
    );
  } else if (numbered) {
    return (
      <ol
        {...props}
        className={classes({
          'list-group-flush': flush,
          'list-group-numbered': true,
        })}
      >
        {props.children}
      </ol>
    );
  } else {
    return (
      <ul {...props} className={classes({ 'list-group-flush': flush })}>
        {props.children}
      </ul>
    );
  }
};

export interface ListGroupItemProps {
  active?: boolean;
  disabled?: boolean;
  isAction?: boolean;
}

const ListGroupItem: FC<
  (HTMLAttributes<HTMLLIElement> | AnchorHTMLAttributes<HTMLAnchorElement>) &
    ListGroupItemProps
> = ({ active, isAction, disabled, ...props }) => {
  const classes = classNames('list-group-item', {
    active,
    disabled,
    'list-group-item-action': isAction,
  });
  if (isAction) {
    const p = props as AnchorHTMLAttributes<HTMLAnchorElement>;
    return (
      <a {...p} className={classNames(props.className, classes)}>
        {props.children}
      </a>
    );
  } else {
    const p = props as HTMLAttributes<HTMLLIElement>;
    return (
      <li {...p} className={classNames(props.className, classes)}>
        {props.children}
      </li>
    );
  }
};

export { ListGroup, ListGroupItem };
