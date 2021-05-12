import { Accordion, AccordionItem } from './accordion.component';
import { loremIpsum } from 'lorem-ipsum';

const lp = loremIpsum({ count: 5 });

export default {
  title: 'Accordion',
  decorators: [(Story: any) => <Story />],
};

export const basic = () => (
  <Accordion>
    <AccordionItem title="This doesn't work yet" body={lp} />
    <AccordionItem title="This doesn't work yet" body={lp} />
    <AccordionItem title="This doesn't work yet" body={lp} />
    <AccordionItem title="This doesn't work yet" body={lp} />
  </Accordion>
);
