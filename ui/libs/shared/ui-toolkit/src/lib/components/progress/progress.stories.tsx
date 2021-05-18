import Card from '../card/card.component';
import Progress from './progress.component';

export default {
  title: 'Progress',
  component: Progress,
};

export const Basic = () => (
  <>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Basic</Card.Title>
        <Progress className="m-4" progress={0} />
        <Progress className="m-4" progress={25} />
        <Progress className="m-4" progress={50} />
        <Progress className="m-4" progress={75} />
        <Progress className="m-4" progress={100} />
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Labeled</Card.Title>
        <Progress className="m-4" showValue progress={0} />
        <Progress className="m-4" showValue progress={25} />
        <Progress className="m-4" showValue progress={50} />
        <Progress className="m-4" showValue progress={75} />
        <Progress className="m-4" showValue progress={100} />
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Height</Card.Title>
        <Progress className="m-4" height={1} progress={25} />
        <Progress className="m-4" height={10} progress={25} />
        <Progress className="m-4" progress={25} />
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Colors</Card.Title>
        <Progress className="m-4" theme="success" progress={25} />
        <Progress className="m-4" theme="info" progress={50} />
        <Progress className="m-4" theme="warning" progress={75} />
        <Progress className="m-4" theme="danger" progress={100} />
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Multiple bars</Card.Title>
        <Progress>
          <Progress.Bar progress={15} />
          <Progress.Bar progress={30} theme="success" />
          <Progress.Bar progress={20} theme="warning" />
        </Progress>
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Striped</Card.Title>
        <Progress className="m-4" striped progress={10} />
        <Progress className="m-4" striped theme="success" progress={25} />
        <Progress className="m-4" striped theme="info" progress={50} />
        <Progress className="m-4" striped theme="warning" progress={75} />
        <Progress className="m-4" striped theme="danger" progress={100} />
      </Card.Body>
    </Card>
    <Card className="mb-4">
      <Card.Body>
        <Card.Title>Animated</Card.Title>
        <Progress className="m-4" animated progress={10} />
        <Progress className="m-4" animated theme="success" progress={25} />
        <Progress className="m-4" animated theme="info" progress={50} />
        <Progress className="m-4" animated theme="warning" progress={75} />
        <Progress className="m-4" animated theme="danger" progress={100} />
      </Card.Body>
    </Card>
  </>
);
