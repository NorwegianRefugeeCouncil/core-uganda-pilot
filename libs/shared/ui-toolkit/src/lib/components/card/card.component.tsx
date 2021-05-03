import {
  AnchorHTMLAttributes,
  forwardRef,
  FunctionComponent,
  HTMLAttributes,
} from 'react';
import classNames from 'classnames';

type CardProps = HTMLAttributes<HTMLDivElement>;

export const Card: FunctionComponent<CardProps> = (props, ref) => {
  const className = classNames('card', props.className);
  return (
    <div ref={ref} {...props} className={className}>
      {props.children}
    </div>
  );
};

type CardBodyProps = HTMLAttributes<HTMLDivElement>;

export const CardBody: FunctionComponent<CardBodyProps> = (props, ref) => {
  return (
    <div
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-body')}
    >
      {props.children}
    </div>
  );
};

type CardTextProps = HTMLAttributes<HTMLParagraphElement>;

export const CardText: FunctionComponent<CardTextProps> = (props, ref) => {
  return (
    <p
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-text')}
    >
      {props.children}
    </p>
  );
};

type CardTitleProps = HTMLAttributes<HTMLHeadingElement>;

export const CardTitle: FunctionComponent<CardTitleProps> = (props, ref) => {
  return (
    <h5
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-title')}
    >
      {props.children}
    </h5>
  );
};

type CardSubTitleProps = HTMLAttributes<HTMLHeadingElement>;

export const CardSubTitle: FunctionComponent<CardSubTitleProps> = forwardRef<
  HTMLHeadingElement,
  CardSubTitleProps
>((props, ref) => {
  return (
    <h6
      {...props}
      ref={ref}
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
});

type CardLinkProps = AnchorHTMLAttributes<HTMLAnchorElement>;

export const CardLink: FunctionComponent<CardLinkProps> = forwardRef<
  HTMLAnchorElement,
  CardLinkProps
>((props, ref) => {
  return (
    <a
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-link')}
    >
      {props.children}
    </a>
  );
});

type CardTopImageProps = HTMLAttributes<HTMLImageElement>;

export const CardTopImage: FunctionComponent<CardTopImageProps> = forwardRef<
  HTMLImageElement,
  CardTopImageProps
>((props, ref) => {
  // 'alt' attribute would be in props
  // eslint-disable-next-line jsx-a11y/alt-text
  return (
    <img
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-img-top')}
    >
      {props.children}
    </img>
  );
});

type CardHeaderProps = HTMLAttributes<HTMLDivElement>;

export const CardHeader: FunctionComponent<CardHeaderProps> = forwardRef<
  HTMLDivElement,
  CardHeaderProps
>((props, ref) => {
  return (
    <div
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-header')}
    >
      {props.children}
    </div>
  );
});

type CardHeaderFeaturedProps = HTMLAttributes<HTMLHeadingElement>;

export const CardHeaderFeatured: FunctionComponent<CardHeaderFeaturedProps> = forwardRef<
  HTMLDivElement,
  CardHeaderFeaturedProps
>((props, ref) => {
  return (
    <h5
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-header')}
    >
      {props.children}
    </h5>
  );
});

type CardFooterProps = HTMLAttributes<HTMLDivElement>;

export const CardFooter: FunctionComponent<CardFooterProps> = forwardRef<
  HTMLDivElement,
  CardFooterProps
>((props, ref) => {
  return (
    <div
      {...props}
      ref={ref}
      className={classNames(props.className, 'card-footer')}
    >
      {props.children}
    </div>
  );
});
