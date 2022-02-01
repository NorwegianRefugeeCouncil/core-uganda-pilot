export default {
  name: {
    minLength: {
      value: 3,
      message: 'Form name needs to be at least 3 characters long',
    },
    maxLength: {
      value: 128,
      message: 'Form name can be at most 128 characters long',
    },
    pattern: {
      value: /[A-Za-z0-9]/,
      message: 'Form name contains invalid characters',
    },
    required: { value: true, message: 'Form name is a required field' },
  },
  selectedField: {
    name: {
      minLength: {
        value: 2,
        message: 'Field name needs to be at least 3 characters long',
      },
      maxLength: {
        value: 128,
        message: 'Field name can be at most 128 characters long',
      },
      required: { value: true, message: 'Field name is required' },
    },
    fieldType: {
      singleSelect: {
        options: {
          minLength: { value: 2, message: 'At least two options are required' },
          maxLength: { value: 60, message: 'At most 60 options are allowed' },
          required: {
            value: true,
            message: 'At least two options are required',
          },
        },
      },
      multiSelect: {
        options: {
          minLength: { value: 2, message: 'At least two options are required' },
          maxLength: { value: 60, message: 'At most 60 options are allowed' },
          required: {
            value: true,
            message: 'At least two options are required',
          },
        },
      },
    },
  },
};
