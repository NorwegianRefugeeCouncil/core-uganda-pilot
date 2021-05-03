import { AnchorHTMLAttributes, FunctionComponent, HTMLAttributes } from 'react';
import classNames from 'classnames';

type CardProps = HTMLAttributes<HTMLDivElement>;

export const Card: FunctionComponent<CardProps> = (props) => {
  const className = classNames('card', props.className);
  return (
    <div {...props} className={className}>
      {props.children}
    </div>
  );
};

type CardBodyProps = HTMLAttributes<HTMLDivElement>;

export const CardBody: FunctionComponent<CardBodyProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-body')}>
      {props.children}
    </div>
  );
};

type CardTextProps = HTMLAttributes<HTMLParagraphElement>;

export const CardText: FunctionComponent<CardTextProps> = (props) => {
  return (
    <p {...props} className={classNames(props.className, 'card-text')}>
      {props.children}
    </p>
  );
};

type CardTitleProps = HTMLAttributes<HTMLHeadingElement>;

export const CardTitle: FunctionComponent<CardTitleProps> = (props) => {
  return (
    <h5 {...props} className={classNames(props.className, 'card-title')}>
      {props.children}
    </h5>
  );
};

type CardSubTitleProps = HTMLAttributes<HTMLHeadingElement>;

export const CardSubTitle: FunctionComponent<CardSubTitleProps> = (props) => {
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

type CardLinkProps = AnchorHTMLAttributes<HTMLAnchorElement>;

export const CardLink: FunctionComponent<CardLinkProps> = (props) => {
  return (
    <a {...props} className={classNames(props.className, 'card-link')}>
      {props.children}
    </a>
  );
};

type CardTopImageProps = HTMLAttributes<HTMLImageElement>;

export const CardTopImage: FunctionComponent<CardTopImageProps> = (props) => {
  // 'alt' attribute would be in props
  // eslint-disable-next-line jsx-a11y/alt-text
  return (
    <img {...props} className={classNames(props.className, 'card-img-top')}>
      {props.children}
    </img>
  );
};

type CardHeaderProps = HTMLAttributes<HTMLDivElement>;

export const CardHeader: FunctionComponent<CardHeaderProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-header')}>
      {props.children}
    </div>
  );
};

type CardHeaderFeaturedProps = HTMLAttributes<HTMLHeadingElement>;

export const CardHeaderFeatured: FunctionComponent<CardHeaderFeaturedProps> = (
  props
) => {
  return (
    <h5 {...props} className={classNames(props.className, 'card-header')}>
      {props.children}
    </h5>
  );
};

type CardFooterProps = HTMLAttributes<HTMLDivElement>;

export const CardFooter: FunctionComponent<CardFooterProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-footer')}>
      {props.children}
    </div>
  );
};
