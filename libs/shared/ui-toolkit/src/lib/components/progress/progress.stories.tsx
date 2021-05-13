import Progress from './progress.component';

export default {
  title: 'Progress',
  component: Progress,
};

export const Basic = () => (
  <Progress
    labels={['Idea', 'Draft', 'Prototype', 'Invention']}
    progress={25}
  />
);
