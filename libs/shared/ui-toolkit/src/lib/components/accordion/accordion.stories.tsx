import React from 'react';
import { Container } from '../container/container.component';
import { Accordion, AccordionItem } from './accordion.component';

export default {
  title: 'Accordion',
  decorators: [
    (Story: any) => (
      <Container centerContent>
        <Story />
      </Container>
    ),
  ],
};

export const basic = () => (
  <Accordion>
    <AccordionItem parentId="0" uniqueKey="1" title="dingle" body="dongle" />
    <AccordionItem parentId="0" uniqueKey="2" title="dingle" body="dongle" />
    <AccordionItem parentId="0" uniqueKey="3" title="dingle" body="dongle" />
    <AccordionItem parentId="0" uniqueKey="4" title="dingle" body="dongle" />
  </Accordion>
);
