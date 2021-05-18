import { LoremIpsum } from 'lorem-ipsum';
import { Card } from './card.component';

const lorem = new LoremIpsum();
const fillerText = lorem.generateParagraphs(2);

export default {
  title: 'Card',
  component: Card,
};

export const Basic = () => {
  return (
    <div className={'container'}>
      <div className={'row'}>
        <div className={'col-12 mb-5'}>
          <p>A simple card with just text</p>
          <Card>
            <Card.Body>
              <Card.Text>{fillerText}</Card.Text>
            </Card.Body>
          </Card>
        </div>
      </div>
      <div className={'row'}>
        <div className={'col-12 mb-5'}>
          <p>A card with a title and some text</p>
          <Card>
            <Card.Body>
              <Card.Title>Title</Card.Title>
              <Card.Text>{fillerText}</Card.Text>
            </Card.Body>
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
            <Card.Img style={{ backgroundColor: 'gray', height: '200px' }} />
            <Card.Header>Header</Card.Header>
            <Card.Body>
              <Card.Title>Title</Card.Title>
              <Card.Subtitle>Subtitle</Card.Subtitle>
              <Card.Text>{fillerText}</Card.Text>
              <Card.Link href={'#'}>link 1</Card.Link>
              <Card.Link href={'#'}>link 2</Card.Link>
            </Card.Body>
            <Card.Footer>Footer</Card.Footer>
          </Card>
        </div>
      </div>
    </div>
  );
};
