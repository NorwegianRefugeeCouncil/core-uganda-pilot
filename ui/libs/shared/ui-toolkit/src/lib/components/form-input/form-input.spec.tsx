import { fireEvent, render, screen } from '@testing-library/react';
import { FC } from 'react';
import { FormInput, FormInputProps, FormLabel } from '@nrc.no/ui';
import {
  useForm,
  FormProvider,
  useFormContext,
  SubmitHandler,
  SubmitErrorHandler,
} from 'react-hook-form';

type ContextProps = {
  onValid: SubmitHandler<any>;
  onInvalid: SubmitErrorHandler<any>;
};

const Context: FC<ContextProps> = ({
  onValid,
  onInvalid,
  children,
}) => {
  const methods = useForm();
  methods.watch((a) => {
    console.log(a);
  });
  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(onValid, onInvalid)}>
        {children}
      </form>
    </FormProvider>
  );
};

type TestCase = {
  name: string;
  input: any;
  expected: any;
  props: Partial<FormInputProps>;
};

type TestCaseProps = {
  testCase: TestCase;
};

const FormField: FC<TestCaseProps> = ({
  testCase,
  ...props
}) => {
  const { register } = useFormContext();
  return (
    <div>
      <FormLabel htmlFor={'id'}>Label</FormLabel>
      <FormInput
        id={'id'}
        {...register('test')}
        type="number"
        {...testCase.props}
      />
    </div>
  );
};

describe('form-input', () => {
  describe('number', () => {
    describe('test cases', () => {
      const testCases: TestCase[] = [
        { name: 'regular', input: '2', expected: 2, props: { type: 'number' } },
        {
          name: 'min',
          input: '-3',
          expected: 2,
          props: { type: 'number', min: 2 },
        },
        {
          name: 'max',
          input: '4',
          expected: 3,
          props: { type: 'number', max: 3 },
        },
        {
          name: 'maxStr',
          input: '4',
          expected: 3,
          props: { type: 'number', max: '3' as any },
        },
        {
          name: 'minStr',
          input: '1',
          expected: 2,
          props: { type: 'number', min: '2' as any },
        },
        {
          name: 'step',
          input: '3',
          expected: 4,
          props: { type: 'number', step: '4' as any },
        },
        {
          name: 'decimalStep',
          input: '3.33',
          expected: 3.5,
          props: { type: 'number', step: '0.5' as any },
        },
        {
          name: 'null',
          input: null,
          expected: undefined,
          props: { type: 'number' },
        },
        {
          name: 'undefined',
          input: undefined,
          expected: undefined,
          props: { type: 'number' },
        },
        {
          name: 'badNumber',
          input: 'abc',
          expected: undefined,
          props: { type: 'number' },
        },
        {
          name: 'emptyString',
          input: '',
          expected: undefined,
          props: { type: 'number' },
        },
        {
          name: 'spaces',
          input: '  ',
          expected: undefined,
          props: { type: 'number' },
        },
        {
          name: 'hex',
          input: '0xFFF',
          expected: undefined,
          props: { type: 'number' },
        },
        { name: 'zeros', input: '000', expected: 0, props: { type: 'number' } },
        {
          name: 'zeros',
          input: '010',
          expected: 10,
          props: { type: 'number' },
        },
      ];

      testCases.forEach((testCase) => {
        it(testCase.name, function () {
          let value = testCase.props.value;
          const setValue = (e: any) => {
            value = e;
          };

          render(
            <Context onValid={setValue} onInvalid={(e) => console.error(e)}>
              <FormField testCase={testCase} />
            </Context>
          );

          fireEvent.input(screen.getByLabelText('Label'), {
            target: { value: testCase.input },
          });

          expect(value).toBe(testCase.expected);
        });
      });
    });
  });
});
