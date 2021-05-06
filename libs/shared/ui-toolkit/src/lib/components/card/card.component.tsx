import * as React from 'react';
import classNames from 'classnames';

type CardProps = React.ComponentPropsWithRef<'div'>;

export const Card: React.FC<CardProps> = (props) => {
  const className = classNames('card', props.className);
  return (
    <div {...props} className={className}>
      {props.children}
    </div>
  );
};

type CardBodyProps = React.ComponentPropsWithRef<'div'>;

export const CardBody: React.FC<CardBodyProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-body')}>
      {props.children}
    </div>
  );
};

type CardTextProps = React.ComponentPropsWithRef<'p'>;

export const CardText: React.FC<CardTextProps> = (props) => {
  return (
    <p {...props} className={classNames(props.className, 'card-text')}>
      {props.children}
    </p>
  );
};

type CardTitleProps = React.ComponentPropsWithRef<'h5'>;

export const CardTitle: React.FC<CardTitleProps> = (props) => {
  return (
    <h5 {...props} className={classNames(props.className, 'card-title')}>
      {props.children}
    </h5>
  );
};

type CardSubTitleProps = React.ComponentPropsWithRef<'h6'>;

export const CardSubTitle: React.FC<CardSubTitleProps> = (props) => {
  return (
    <h6
      {...props}
      className={classNames(
        props.className,
        'card-subtitle',
        'text-muted',
        'mb-2'
      )}
    >
      {props.children}
    </h6>
  );
};

type CardLinkProps = React.ComponentPropsWithRef<'a'>;

export const CardLink: React.FC<CardLinkProps> = (props) => {
  return (
    <a {...props} className={classNames(props.className, 'card-link')}>
      {props.children}
    </a>
  );
};

type CardTopImageProps = React.ComponentPropsWithRef<'img'>;

export const CardTopImage: React.FC<CardTopImageProps> = (props) => {
  // 'alt' attribute would be in props
  // eslint-disable-next-line jsx-a11y/alt-text
  return (
    <img {...props} className={classNames(props.className, 'card-img-top')}>
      {props.children}
    </img>
  );
};

type CardHeaderProps = React.ComponentPropsWithRef<'div'>;

export const CardHeader: React.FC<CardHeaderProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-header')}>
      {props.children}
    </div>
  );
};

type CardHeaderFeaturedProps = React.ComponentPropsWithRef<'h5'>;

export const CardHeaderFeatured: React.FC<CardHeaderFeaturedProps> = (
  props
) => {
  return (
    <h5 {...props} className={classNames(props.className, 'card-header')}>
      {props.children}
    </h5>
  );
};

type CardFooterProps = React.ComponentPropsWithRef<'div'>;

export const CardFooter: React.FC<CardFooterProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-footer')}>
      {props.children}
    </div>
  );
};
