import { Accordion } from './accordion.component';
import { loremIpsum } from 'lorem-ipsum';

const lp = loremIpsum({ units: 'paragraph', count: 6 });

export default {
  title: 'Accordion',
  component: Accordion,
};

export const basic = () => (
  <Accordion activeId="0">
    <Accordion.Item id="0" header="Accordion Item #1" body={lp} />
    <Accordion.Item id="1" header="Accordion Item #2" body={lp} />
    <Accordion.Item id="2" header="Accordion Item #3" body={lp} />
    <Accordion.Item id="3" header="Accordion Item #4" body={lp} />
  </Accordion>
);
