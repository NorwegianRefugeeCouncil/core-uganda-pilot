import { storiesOf } from '@storybook/react';
import { ProgressBar } from './progress-bar.component';

storiesOf('ProgressBar', module).add('default', () => (
  <ProgressBar labels={['A', 'B', 'C']} progress={2} />
));
