import { FC, useCallback, useEffect, useMemo } from 'react';
import { useForm } from 'react-hook-form';
import classNames from 'classnames';

import { useApiClient, useFormValidation } from '../../../hooks/hooks';
import { Organization } from '../../../types/types';

type Props = {
  id?: string;
  organization: Organization;
};

type FormData = {
  name: string;
  issuer: string;
  clientId: string;
  clientSecret: string;
  organizationId: string;
  emailDomain: string;
};

export const IdentityProviderEditor: FC<Props> = (props) => {
  const { id, organization } = props;

  const isNew = useMemo(() => !id, [id]);

  const apiClient = useApiClient();

  const form = useForm<FormData>({ mode: 'onChange' });

  const {
    register,
    handleSubmit,
    setValue,
    formState: { isSubmitting },
  } = form;

  const { fieldErrors, fieldClasses } = useFormValidation(isNew, form);

  useEffect(() => {
    if (id) {
      apiClient.getIdentityProvider({ id }).then((resp) => {
        if (resp.response) {
          setValue('name', resp.response.name);
          setValue('clientId', resp.response.clientId);
          setValue('organizationId', resp.response.organizationId);
          setValue('issuer', resp.response.domain);
          setValue('emailDomain', resp.response.emailDomain);
          setValue('clientSecret', '');
        }
      });
    }
  }, [apiClient, id, setValue]);

  const onSubmit = useCallback(
    (args: FormData) => {
      const obj = {
        id,
        name: args.name,
        clientId: args.clientId,
        clientSecret: args.clientSecret,
        domain: args.issuer,
        organizationId: organization.id,
        emailDomain: args.emailDomain,
      };
      if (id) {
        return apiClient.updateIdentityProvider({
          object: obj,
        });
      }
      return apiClient.createIdentityProvider({
        object: obj,
      });
    },
    [apiClient, id, organization.id],
  );

  return (
    <div className={classNames('card bg-dark border-secondary')}>
      <div className="card-body">
        <form className="needs-validation" noValidate onSubmit={handleSubmit(onSubmit)}>
          <div className={classNames('form-group mb-2')}>
            <label className="form-label text-light">Name</label>
            <input
              {...register('name', {
                required: true,
                pattern: /^[a-zA-Z0-9\-_ ]+$/,
              })}
              className={classNames('form-control form-control-darkula', fieldClasses('name'))}
            />
            {fieldErrors('name')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Issuer</label>
            <input
              {...register('issuer', {
                required: true,
                pattern: /^https?:\/\/[a-zA-Z0-9.\-_]+(:[0-9]+)?$/,
              })}
              className={classNames('form-control form-control-darkula', fieldClasses('issuer'))}
            />
            {fieldErrors('issuer')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Email Domain</label>
            <input
              {...register('emailDomain', {
                required: true,
              })}
              className={classNames('form-control form-control-darkula', fieldClasses('emailDomain'))}
            />
            {fieldErrors('emailDomain')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Client ID</label>
            <input
              {...register('clientId', {
                required: true,
              })}
              className={classNames('form-control form-control-darkula', fieldClasses('clientId'))}
            />
            {fieldErrors('clientId')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Client Secret</label>
            <input
              type="password"
              {...register('clientSecret', {
                required: isNew,
              })}
              className={classNames('form-control form-control-darkula', fieldClasses('clientSecret'))}
              placeholder={isNew ? '' : '********'}
            />
            {fieldErrors('clientSecret')}
          </div>
          <button disabled={isSubmitting} className="btn btn-success mt-2">
            {props.id ? 'Update Identity Provider' : 'Create Identity Provider'}
          </button>
        </form>
      </div>
    </div>
  );
};
