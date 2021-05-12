import React from 'react';
import { Container } from '../container/container.component';
import { Accordion, AccordionItem } from './accordion.component';
import { loremIpsum } from 'lorem-ipsum';

const lp = loremIpsum({ count: 5 });

export default {
  title: 'Accordion',
  decorators: [(Story: any) => <Story />],
};

export const basic = () => (
  <Accordion>
    <AccordionItem parentId="0" uniqueKey="1" title="dingle" body={lp} />
    <AccordionItem parentId="0" uniqueKey="2" title="dongle" body={lp} />
    <AccordionItem parentId="0" uniqueKey="3" title="dangle" body={lp} />
    <AccordionItem parentId="0" uniqueKey="4" title="dungle" body={lp} />
  </Accordion>
);
