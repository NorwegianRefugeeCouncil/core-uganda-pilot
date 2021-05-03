import { Fragment } from 'react';
import { storiesOf } from '@storybook/react';
import { FormCheck, FormCheckInput, FormCheckLabel, FormInput, FormRadioInput } from './form-input.component';
import { FormLabel } from '../form-label/label.component';

storiesOf('Input', module)
  .add('default', () => {
    return (
      <Fragment>
        <div className={'container'}>
          <div className={'row'}>
            <div className={'col-12 mb-5'}>
              <FormLabel>Blank:</FormLabel>
              <FormInput />
            </div>
            <div className={'col-12 mb-5'}>
              <FormLabel>With Placeholder:</FormLabel>
              <FormInput placeholder={'placeholder'} />
            </div>
            <div className={'col-12 mb-5'}>
              <FormLabel>Disabled Input:</FormLabel>
              <FormInput placeholder={'disabled'} disabled={true} />
            </div>
            <div className={'col-12 mb-5'}>
              <FormLabel>Plaintext Input:</FormLabel>
              <FormInput placeholder={'plaintext'} plaintext={true} readOnly={true} />
            </div>
            <div className={'col-12 mb-5'}>
              <FormLabel>Color Input:</FormLabel>
              <FormInput placeholder={'color'} colorInput={true} type={'color'} />
            </div>
            <div className={'col-12 mb-5'}>
              <p>Radio buttons: </p>
              <form>
                <FormCheck>
                  <FormCheckLabel>Label 1</FormCheckLabel>
                  <FormRadioInput name={'radio-1'} />
                </FormCheck>
                <FormCheck>
                  <FormCheckLabel>Label 2</FormCheckLabel>
                  <FormRadioInput name={'radio-1'} />
                </FormCheck>
                <FormCheck>
                  <FormCheckLabel>Label 3</FormCheckLabel>
                  <FormRadioInput name={'radio-1'} />
                </FormCheck>
              </form>
            </div>
            <div className={'col-12 mb-5'}>
              <p>Check Boxes: </p>
              <form>
                <FormCheck>
                  <FormCheckLabel>Label 1</FormCheckLabel>
                  <FormCheckInput />
                </FormCheck>
                <FormCheck>
                  <FormCheckLabel>Label 2</FormCheckLabel>
                  <FormCheckInput />
                </FormCheck>
                <FormCheck>
                  <FormCheckLabel>Label 3</FormCheckLabel>
                  <FormCheckInput />
                </FormCheck>
              </form>
            </div>
          </div>
        </div>
      </Fragment>
    );
  });
