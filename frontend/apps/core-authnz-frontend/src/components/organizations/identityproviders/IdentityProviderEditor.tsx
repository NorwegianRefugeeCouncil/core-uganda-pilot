import { FC, useCallback, useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import classNames from 'classnames';

import { useApiClient, useFormValidation } from '../../../hooks/hooks';
import { IdentityProvider, Organization } from '../../../types/types';

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
  scopes: string;
  claimMappings: {
    version: string;
    subject: string;
    displayName: string;
    fullName: string;
    email: string;
    emailVerified: string;
  };
};

export const IdentityProviderEditor: FC<Props> = (props) => {
  const { id, organization } = props;

  const isNew = useMemo(() => !id, [id]);

  const apiClient = useApiClient();

  const form = useForm<FormData>({ mode: 'onChange' });

  const [version, setVersion] = useState<string>('0');

  const {
    register,
    handleSubmit,
    setValue,
    formState: { isSubmitting },
  } = form;

  const { fieldErrors, fieldClasses } = useFormValidation(isNew, form);

  const setData = (data: IdentityProvider) => {
    setValue('name', data.name);
    setValue('clientId', data.clientId);
    setValue('organizationId', data.organizationId);
    setValue('issuer', data.domain);
    setValue('emailDomain', data.emailDomain);
    setValue('clientSecret', '');
    setValue('scopes', data.scopes);
    setValue('claimMappings.subject', data.claimMappings.subject);
    setValue('claimMappings.displayName', data.claimMappings.displayName);
    setValue('claimMappings.fullName', data.claimMappings.fullName);
    setValue('claimMappings.email', data.claimMappings.email);
    setValue('claimMappings.emailVerified', data.claimMappings.emailVerified);
    setVersion(data.claimMappings.version);
  };

  useEffect(() => {
    if (id) {
      apiClient.getIdentityProvider({ id }).then((resp) => {
        if (resp.response) {
          setData(resp.response);
        }
      });
    }
  }, [apiClient, id, setValue]);

  const onSubmit = useCallback(
    async (args: FormData) => {
      const newVersion = `${parseInt(version, 10) + 1}`;

      const object = {
        id,
        name: args.name,
        clientId: args.clientId,
        clientSecret: args.clientSecret,
        domain: args.issuer,
        organizationId: organization.id,
        emailDomain: args.emailDomain,
        scopes: args.scopes,
        claimMappings: {
          version: newVersion,
          subject: args.claimMappings.subject,
          displayName: args.claimMappings.displayName,
          fullName: args.claimMappings.fullName,
          email: args.claimMappings.email,
          emailVerified: args.claimMappings.emailVerified,
        },
      };
      let resp;
      if (id) {
        resp = await apiClient.updateIdentityProvider({
          object,
        });
      } else {
        resp = await apiClient.createIdentityProvider({
          object,
        });
      }
      return resp.response && setData(resp.response);
    },
    [apiClient, id, organization.id, version],
  );

  return (
    <div className={classNames('card bg-dark border-secondary')}>
      <div className="card-body">
        <form
          className="needs-validation"
          noValidate
          onSubmit={handleSubmit(onSubmit)}
        >
          <div className={classNames('form-group mb-2')}>
            <label className="form-label text-light">Name</label>
            <input
              {...register('name', {
                required: true,
                pattern: /^[a-zA-Z0-9\-_ ]+$/,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('name'),
              )}
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
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('issuer'),
              )}
            />
            {fieldErrors('issuer')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Email Domain</label>
            <input
              {...register('emailDomain', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('emailDomain'),
              )}
            />
            {fieldErrors('emailDomain')}
          </div>
          <div className="form-group mb-2">
            <label className="form-label text-light">Client ID</label>
            <input
              {...register('clientId', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('clientId'),
              )}
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
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('clientSecret'),
              )}
              placeholder={isNew ? '' : '********'}
            />
            {fieldErrors('clientSecret')}
          </div>

          <div className="form-group mb-2">
            <label className="form-label text-light">Scopes</label>
            <input
              {...register('scopes', {
                required: true,
              })}
              className={classNames(
                'form-control form-control-darkula',
                fieldClasses('scopes'),
              )}
            />
            {fieldErrors('scopes')}
          </div>

          <h6 className="text-light mt-5">
            Claim Mapping, Current Version: {version}
          </h6>
          <p className="form-text text-muted mb-4">
            Please use go template syntax:
            <a
              href="https://blog.gopheracademy.com/advent-2017/using-go-templates/"
              target="_blank"
              rel="noreferrer"
            >
              Examples
            </a>
          </p>

          <div className="form-group row mb-3">
            <div className="col-sm-2 text-light">
              <label htmlFor="subject">Subject</label>
            </div>
            <div className="col-sm-10">
              <input
                {...register('claimMappings.subject', {
                  required: true,
                })}
                id="subject"
                className={classNames(
                  'form-control form-control-darkula',
                  fieldClasses('claimMappings.subject'),
                )}
              />
            </div>
          </div>

          <div className="form-group row mb-3">
            <div className="col-sm-2 text-light">
              <label htmlFor="displayName">Display Name</label>
            </div>
            <div className="col-sm-10">
              <input
                {...register('claimMappings.displayName', {
                  required: true,
                })}
                id="displayName"
                className={classNames(
                  'form-control form-control-darkula',
                  fieldClasses('claimMappings.displayName'),
                )}
              />
            </div>
          </div>

          <div className="form-group row mb-3">
            <div className="col-sm-2 text-light">
              <label htmlFor="fullName">Full Name</label>
            </div>
            <div className="col-sm-10">
              <input
                {...register('claimMappings.fullName', {
                  required: true,
                })}
                id="fullName"
                className={classNames(
                  'form-control form-control-darkula',
                  fieldClasses('claimMappings.fullName'),
                )}
              />
            </div>
          </div>

          <div className="form-group row mb-3">
            <div className="col-sm-2 text-light">
              <label htmlFor="email">Email</label>
            </div>
            <div className="col-sm-10">
              <input
                {...register('claimMappings.email', {
                  required: true,
                })}
                id="email"
                className={classNames(
                  'form-control form-control-darkula',
                  fieldClasses('claimMappings.email'),
                )}
              />
            </div>
          </div>

          <div className="form-group row mb-3">
            <div className="col-sm-2 text-light">
              <label htmlFor="emailVerified">Email Verified</label>
            </div>
            <div className="col-sm-10">
              <input
                {...register('claimMappings.emailVerified', {
                  required: true,
                })}
                id="emailVerified"
                className={classNames(
                  'form-control form-control-darkula',
                  fieldClasses('claimMappings.emailVerified'),
                )}
              />
            </div>
          </div>

          {fieldErrors('claimMappings')}

          <button disabled={isSubmitting} className="btn btn-success mt-2">
            {props.id ? 'Update Identity Provider' : 'Create Identity Provider'}
          </button>
        </form>
      </div>
    </div>
  );
};
