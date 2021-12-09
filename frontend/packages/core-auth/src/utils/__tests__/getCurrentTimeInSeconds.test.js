import getCurrentTimeInSeconds from '../getCurrentTimeInSeconds';

Date.now = jest.fn(() => 1400000000000);

describe('utils/getCurrentTimeInSeconds', () => {
  it('should call qs.stringify', () => {
    expect(getCurrentTimeInSeconds()).toEqual(1400000000);
  });
});
