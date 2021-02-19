import { first } from '../util';
import { getValues } from './internal';

export const getFirstValue = async () => {
  return first(await getValues());
};
