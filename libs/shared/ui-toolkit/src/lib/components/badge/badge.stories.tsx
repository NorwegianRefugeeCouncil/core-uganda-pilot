import Button from '../button/button.component';
import Badge from './badge.component';

export default {
  title: 'Badge',
  decorators: [(Story: any) => <Story />],
};

export const basic = () => (
  <>
    <h1>
      Example heading{' '}
      <Badge className="m-2" theme="secondary">
        New
      </Badge>
    </h1>
    <h2>
      Example heading{' '}
      <Badge className="m-2" theme="secondary">
        New
      </Badge>
    </h2>
    <h3>
      Example heading{' '}
      <Badge className="m-2" theme="secondary">
        New
      </Badge>
    </h3>
    <h4>
      Example heading{' '}
      <Badge className="m-2" theme="secondary">
        New
      </Badge>
    </h4>
    <h5>
      Example heading{' '}
      <Badge className="m-2" theme="secondary">
        New
      </Badge>
    </h5>
    <h6>
      Example heading{' '}
      <Badge className="m-2" theme="secondary">
        New
      </Badge>
    </h6>
  </>
);

export const counter = () => (
  <Button>
    Notifications{' '}
    <Badge className="m-2" theme="secondary">
      4
    </Badge>
  </Button>
);

export const colors = () => (
  <>
    <Badge className="m-2">Primary</Badge>
    <Badge className="m-2" theme="secondary">
      Secondary
    </Badge>
    <Badge className="m-2" theme="success">
      Success
    </Badge>
    <Badge className="m-2" theme="danger">
      Danger
    </Badge>
    <Badge className="m-2" theme="warning">
      Warning
    </Badge>
    <Badge className="m-2" theme="info">
      Info
    </Badge>
    <Badge className="m-2" theme="light">
      Light
    </Badge>
    <Badge className="m-2" theme="dark">
      Dark
    </Badge>
  </>
);

export const pills = () => (
  <>
    <Badge className="m-2" pill>
      Primary
    </Badge>
    <Badge className="m-2" pill theme="secondary">
      Secondary
    </Badge>
    <Badge className="m-2" pill theme="success">
      Success
    </Badge>
    <Badge className="m-2" pill theme="danger">
      Danger
    </Badge>
    <Badge className="m-2" pill theme="warning">
      Warning
    </Badge>
    <Badge className="m-2" pill theme="info">
      Info
    </Badge>
    <Badge className="m-2" pill theme="light">
      Light
    </Badge>
    <Badge className="m-2" pill theme="dark">
      Dark
    </Badge>
  </>
);
