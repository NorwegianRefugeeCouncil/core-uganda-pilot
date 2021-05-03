import { AnchorHTMLAttributes, forwardRef, FunctionComponent, HTMLAttributes } from 'react';
import { addClasses } from '../../utils/utils';

type CardProps = HTMLAttributes<HTMLDivElement>

export const Card = forwardRef<HTMLDivElement, CardProps>(
  (props, ref) => {
    const className = props.className ? props.className + ' card' : 'card';
    return (<div ref={ref} {...props} className={className}>{props.children}</div>);
  });

type CardBodyProps = HTMLAttributes<HTMLDivElement>

export const CardBody = forwardRef<HTMLDivElement, CardBodyProps>(
  (props, ref) => {
    return (<div {...props} ref={ref} className={addClasses(props.className, 'card-body')}>{props.children}</div>);
  });

type CardTextProps = HTMLAttributes<HTMLParagraphElement>

export const CardText = forwardRef<HTMLParagraphElement, CardTextProps>(
  (props, ref) => {
    return (<p {...props} ref={ref} className={addClasses(props.className, 'card-text')}>{props.children}</p>);
  });

type CardTitleProps = HTMLAttributes<HTMLHeadingElement>

export const CardTitle = forwardRef<HTMLHeadingElement, CardTitleProps>(
  (props, ref) => {
    return (<h5 {...props} ref={ref} className={addClasses(props.className, 'card-title')}>{props.children}</h5>);
  });

type CardSubTitleProps = HTMLAttributes<HTMLHeadingElement>

export const CardSubTitle = forwardRef<HTMLHeadingElement, CardSubTitleProps>(
  (props, ref) => {
    return (<h6 {...props} ref={ref}
                className={addClasses(props.className, 'card-subtitle', 'text-muted', 'mb-2')}>{props.children}</h6>);
  });

type CardLinkProps = AnchorHTMLAttributes<HTMLAnchorElement>

export const CardLink = forwardRef<HTMLAnchorElement, CardLinkProps>(
  (props, ref) => {
    return (<a {...props} ref={ref} className={addClasses(props.className, 'card-link')}>{props.children}</a>);
  });

type CardTopImageProps = HTMLAttributes<HTMLImageElement>

export const CardTopImage = forwardRef<HTMLImageElement, CardTopImageProps>(
  (props, ref) => {
    return (<img {...props} ref={ref} className={addClasses(props.className, 'card-img-top')}>{props.children}</img>);
  });

type CardHeaderProps = HTMLAttributes<HTMLDivElement>

export const CardHeader = forwardRef<HTMLDivElement, CardHeaderProps>(
  (props, ref) => {
    return (<div {...props} ref={ref} className={addClasses(props.className, 'card-header')}>{props.children}</div>);
  });

type CardHeaderFeaturedProps = HTMLAttributes<HTMLHeadingElement>

export const CardHeaderFeatured = forwardRef<HTMLDivElement, CardHeaderFeaturedProps>(
  (props, ref) => {
    return (<h5 {...props} ref={ref} className={addClasses(props.className, 'card-header')}>{props.children}</h5>);
  });

type CardFooterProps = HTMLAttributes<HTMLDivElement>

export const CardFooter = forwardRef<HTMLDivElement, CardFooterProps>(
  (props, ref) => {
    return (<div {...props} ref={ref} className={addClasses(props.className, 'card-footer')}>{props.children}</div>);
  });

