export default {
  name: {
    minLength: {
      value: 3,
      message: 'This needs to be at least 3 characters long',
    },
    pattern: {
      value: /[A-Za-z]/,
      message: 'This contains invalid characters',
    },
    required: { value: true, message: 'This is required' },
  },
};
