import * as React from 'react';
import { VStack } from 'native-base';
import { FieldKind } from 'core-api-client';

import { SubformFieldValueComponent } from './SubformFieldValue.component';
import { NormalisedFieldValue } from './RecordView.types';
import { FieldValueComponent } from './FieldValue.component';

type Props = {
  fieldValues: NormalisedFieldValue[];
};

export const RecordViewComponent: React.FC<Props> = ({ fieldValues }) => {
  return (
    <VStack space={4}>
      {fieldValues.map((f, i) => {
        if (f.fieldType === FieldKind.SubForm) {
          return (
            <SubformFieldValueComponent
              key={i}
              header={f.header}
              labels={f.labels}
              items={f.values}
            />
          );
        }

        return <FieldValueComponent key={i} item={f} />;
      })}
    </VStack>
  );
};
