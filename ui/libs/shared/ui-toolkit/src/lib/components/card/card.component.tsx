/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';
import { CardImg } from './card-img.component';
import { CardBody } from './card-body.component';
import { CardHeader } from './card-header.component';
import { CardFooter } from './card-footer.component';
import { CardText } from './card-text.component';
import { CardTitle } from './card-title.component';
import { CardSubtitle } from './card-subtitle.component';
import { CardLink } from './card-link.component';

export interface CardProps extends React.ComponentPropsWithoutRef<'div'> {}

type Card = React.FC<CardProps> & {
  Img: typeof CardImg;
  Title: typeof CardTitle;
  Subtitle: typeof CardSubtitle;
  Body: typeof CardBody;
  Link: typeof CardLink;
  Text: typeof CardText;
  Header: typeof CardHeader;
  Footer: typeof CardFooter;
};

const Card: Card = ({ className: customClass, children, ...rest }) => {
  const className = classNames('card', customClass);
  return (
    <div {...rest} className={className}>
      {children}
    </div>
  );
};

Card.displayName = 'Card';

Card.Img = CardImg;
Card.Title = CardTitle;
Card.Subtitle = CardSubtitle;
Card.Body = CardBody;
Card.Link = CardLink;
Card.Text = CardText;
Card.Header = CardHeader;
Card.Footer = CardFooter;

export { Card };
