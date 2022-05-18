import {useForm, useFieldArray, FieldNamesMarkedBoolean} from 'react-hook-form';
import {Column, Table} from './client';
import classnames from 'classnames';
import {useCallback} from 'react';
import {object, array, string} from 'yup';
import {yupResolver} from '@hookform/resolvers/yup';

type CreateTableProps = {
  onSubmit: (values: Table) => void;
  onCancel: () => void;
};

function CreateTable(props: CreateTableProps) {

  const {onSubmit, onCancel} = props;

  const validation = object().shape({
    name: string().required(),
    columns: array().of(object().shape({
      name: string().required(),
      type: string().required()
    })).required().min(1)
  });

  const {
    control,
    register,
    handleSubmit,
    formState: {errors, isValid, dirtyFields, touchedFields}
  } = useForm<Table>({
    mode: 'all',
    resolver: yupResolver(validation)
  });

  const {fields, append, remove} = useFieldArray<Table>({
    control,
    name: 'columns',
  });

  const validationClasses = useCallback((name: keyof FieldNamesMarkedBoolean<Table>) => {
    return {
      'dirty': dirtyFields[name],
      'touched': touchedFields[name],
    };
  }, [errors, touchedFields]);

  return <div>

    <div className={'card'}>
      <div className={'card-body'}>
        <h5 className={'card-title'}>Create Table</h5>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className={classnames('form-group', 'mb-3', validationClasses('name'))}>
            <label htmlFor="name">Table Name</label>
            <input type="text" className="form-control" id="name" {...register('name', {required: true})}/>
            {errors.name && <div className="invalid-feedback">Name is required</div>}
          </div>

          <button
            className={'btn btn-outline-secondary mb-3'}
            onClick={() => append({name: '', type: 'string'})}>
            Add column
          </button>

          <div className={'mb-3'}>
            {fields.map((field, index) => {
              return <div className={'row'} key={field.id}>
                <div className={'col'}>

                  <div className={classnames('form-group', 'mb-3', {
                    dirty: dirtyFields?.columns?.[index].name,
                    touched: touchedFields?.columns?.[index]?.name
                  })}>
                    <label htmlFor={`column-${index}`}>Column Name</label>
                    <input className={'form-control'}
                           placeholder="name"
                           {...register(`columns.${index}.name`, {required: true, minLength: 1})}
                    />
                    {errors?.columns?.[index]?.name && <div className="invalid-feedback">Name is required</div>}
                  </div>
                </div>
                <div className={'col mb-3'}>
                  <label htmlFor={`column-${index}`}>Column Type</label>
                  <select className={'form-control'} {...register(`columns.${index}.type`, {
                    required: true,
                    minLength: 1
                  })}>
                    <option value="varchar">varchar</option>
                    <option value="boolean">boolean</option>
                    <option value="int">int</option>
                  </select>
                </div>

                <div className={'col'}>
                  <label/>
                  <button className={'btn d-block btn-secondary'} onClick={() => remove(index)}>Remove</button>
                </div>
              </div>;
            })}

          </div>

          <button
            disabled={!isValid}
            className={'btn btn-primary'}
            type={'submit'}>
            Submit
          </button>

          <button
            className={'btn btn-secondary'}
            onClick={() => onCancel()}>
            Cancel
          </button>

        </form>
      </div>
    </div>

  </div>;
}

export default CreateTable;
