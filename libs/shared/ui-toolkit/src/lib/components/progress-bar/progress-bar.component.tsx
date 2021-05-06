import * as React from 'react';

type ProgessBarProps =  React.ComponentPropsWithRef<'div'> & {
  labels : string[],
  progress : number
};

const renderProgressColumn = (i:number, numberOfColumns: number) => {
  if (i === 0 || i === numberOfColumns -1 ){
    return <td><span></span></td>
  }
  return <td><span>-</span></td>
}

const renderTextColumn = (i : number, labels: string[]) => {
  if (i % 2 === 0) {
    return <td colSpan={2} align='center' style={{width: 30}}><span>{labels[i/2]}</span></td>
  } else {
    return <td style={{width: 50}}><span></span></td>
  }
}

export const ProgressBar: React.FC<ProgessBarProps> = ({labels, progress, ...props}) => {

  const numberOfColumns = (labels.length * 2) + (labels.length - 1)
  const progressColumns = []
  for (let i = 0; i < numberOfColumns; i++) {
    progressColumns.push(<td>{renderProgressColumn(i, numberOfColumns)}</td>)
  }

  const numberOfTextColumns = labels.length * 2 - 1
  const textColumns = []
  for (let i = 0; i < numberOfTextColumns; i++){
    textColumns.push(renderTextColumn(i, labels))
  }

  return <table>
      <tr>
        {progressColumns}
      </tr>
      <tr>
        {textColumns}
      </tr>
  </table>;
};
