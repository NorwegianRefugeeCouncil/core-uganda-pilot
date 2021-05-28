import { FormDefinitionVersion } from '@core/api-client';
import * as React from 'react';
import { RootBuilderContainer } from './root.builder';

type VersionBuilderProps = {
  path: string,
  version: FormDefinitionVersion
}

/**
 * Renders the form builder for a given section
 * @param props
 * @constructor
 */
export const VersionBuilder: React.FC<VersionBuilderProps> = (props: VersionBuilderProps) => {
  const { version, path } = props;

  if (!version) {
    return <div />;
  }
  return <div>

    {/* Is the version served ? */}
    <div className='form-check'>
      <input disabled
             className='form-check-input'
             type='checkbox'
             value=''
             checked={version.served}
      />

      <label className='form-check-label text-dark'>
        Served
      </label>

    </div>

    {/* Is it a storage version ? */}
    <div className='form-check'>

      <input disabled
             className='form-check-input'
             type='checkbox'
             value=''
             checked={version.storage} />

      <label className='form-check-label text-dark'>
        Storage
      </label>

    </div>

    {/* render the builder for the content of the form version */}
    <div className={'mt-3'}>
      <RootBuilderContainer
        path={path + '.schema.formSchema.root'}
        root={version.schema.formSchema.root}
      />
    </div>

  </div>;
};
