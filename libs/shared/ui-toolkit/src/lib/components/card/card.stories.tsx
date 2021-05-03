import { storiesOf } from '@storybook/react';
import {
  Card,
  CardBody,
  CardFooter,
  CardHeaderFeatured,
  CardLink,
  CardSubTitle,
  CardText,
  CardTitle,
  CardTopImage,
} from './card.component';

storiesOf('Input', module).add('default', () => {
  return (
    <div className={'container'}>
      <div className={'row'}>
        <div className={'col-12 mb-5'}>
          <p>A simple card with just text</p>
          <Card>
            <CardBody>
              <CardText>some text</CardText>
            </CardBody>
          </Card>
        </div>
      </div>
      <div className={'row'}>
        <div className={'col-12 mb-5'}>
          <p>A card with a title and some text</p>
          <Card>
            <CardBody>
              <CardTitle>Title</CardTitle>
              <CardText>some text</CardText>
            </CardBody>
          </Card>
        </div>
      </div>
      <div className={'row'}>
        <div className={'col-12 mb-5'}>
          <p>
            A full blown card, with a top image, a header, a title, subtitle,
            text, links, shadow and footer
          </p>
          <Card className={'shadow'}>
            <CardTopImage
              style={{ backgroundColor: 'gray', height: '200px' }}
            />
            <CardHeaderFeatured>Wow</CardHeaderFeatured>
            <CardBody>
              <CardTitle>Title</CardTitle>
              <CardSubTitle>Subtitle</CardSubTitle>
              <CardText>some text</CardText>
              <CardLink href={'#'}>link 1</CardLink>
              <CardLink href={'#'}>link 2</CardLink>
            </CardBody>
            <CardFooter>Footer here</CardFooter>
          </Card>
        </div>
      </div>
    </div>
  );
});
