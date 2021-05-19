import { Accordion } from './accordion.component';
import { loremIpsum } from 'lorem-ipsum';

const lp = loremIpsum({ count: 5 });

export default {
  title: 'Accordion',
  component: Accordion,
};

export const basic = () => (
  <Accordion>
    <Accordion.Item title="This doesn't work yet" body={lp} />
    <Accordion.Item title="This doesn't work yet" body={lp} />
    <Accordion.Item title="This doesn't work yet" body={lp} />
    <Accordion.Item title="This doesn't work yet" body={lp} />
  </Accordion>
);
