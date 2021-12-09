import { FC } from 'react';
import classNames from 'classnames';

export type SectionTitleProps = {
  title: string;
  className?: string;
};

export const SectionTitle: FC<SectionTitleProps> = (props) => {
  const { title, children } = props;
  return (
    <div
      className={classNames('border-bottom border-secondary pb-3 my-2 d-flex flex-row justify-content-center', props.className)}
    >
      <span className="flex-grow-1 fs-5">{title}</span>
      {children}
    </div>
  );
};
