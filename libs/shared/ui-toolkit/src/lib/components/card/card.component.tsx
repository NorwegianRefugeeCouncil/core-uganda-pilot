/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@ui-helpers/utils';
import CardImg from './cardimg.component';
import CardBody from './cardbody.component';
import CardHeader from './cardheader.component';
import CardFooter from './cardfooter.component';
import CardText from './cardtext.component';
import CardTitle from './cardtitle.component';
import CardSubtitle from './cardsubtitle.component';
import CardLink from './cardlink.component';

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

export default Card;
