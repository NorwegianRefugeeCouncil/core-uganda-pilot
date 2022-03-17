import { FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

export const prettifyData = (
  originalData: FormWithRecord<Recipient>[],
): FormWithRecord<Recipient>[] => {
  // don't show key fields
  const dataWithoutKeys = originalData.map((ancestor) => {
    const noKeys = ancestor.form.fields.filter((field) => !field.key);

    return {
      ...ancestor,
      form: {
        ...ancestor.form,
        fields: noKeys,
      },
    };
  });

  // merge fields of individual into individual beneficiary form
  if (dataWithoutKeys.length > 1) {
    const mergedIndividual = [
      {
        form: {
          ...dataWithoutKeys[1].form,
          fields: [
            ...dataWithoutKeys[0].form.fields,
            ...dataWithoutKeys[1].form.fields,
          ],
        },
        record: {
          ...dataWithoutKeys[1].record,
          values: [
            ...dataWithoutKeys[0].record.values,
            ...dataWithoutKeys[1].record.values,
          ],
        },
      },
    ];
    return mergedIndividual.concat(dataWithoutKeys.slice(2));
  }
  return dataWithoutKeys;
};
