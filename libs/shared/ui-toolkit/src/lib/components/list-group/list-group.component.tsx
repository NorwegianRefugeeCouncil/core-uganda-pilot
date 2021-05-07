import * as React from 'react';
import classNames from 'classnames';
export interface ListGroupProps
  extends React.HTMLAttributes<HTMLDivElement | HTMLUListElement> {
  flush?: boolean;
  isActionListGroup?: boolean;
  numbered?: boolean;
}

const ListGroup: React.FC<ListGroupProps> = ({
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

const ListGroupItem: React.FC<
  (
    | React.HTMLAttributes<HTMLLIElement>
    | React.AnchorHTMLAttributes<HTMLAnchorElement>
  ) &
    ListGroupItemProps
> = ({ active, isAction, disabled, ...props }) => {
  const classes = classNames('list-group-item', {
    active,
    disabled,
    'list-group-item-action': isAction,
  });
  if (isAction) {
    const p = props as React.AnchorHTMLAttributes<HTMLAnchorElement>;
    return (
      <a {...p} className={classNames(props.className, classes)}>
        {props.children}
      </a>
    );
  } else {
    const p = props as React.HTMLAttributes<HTMLLIElement>;
    return (
      <li {...p} className={classNames(props.className, classes)}>
        {props.children}
      </li>
    );
  }
};

export { ListGroup, ListGroupItem };
