// See: https://orval.dev/guides/custom-axios

import Axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { ErrorResponse } from './panda_studio.gen';
 
export const AXIOS_INSTANCE = Axios.create({
    baseURL: import.meta.env.API_URL ?? '/api'
}); // use your own URL here or environment variable

// add a second `options` argument here if you want to pass extra options to each generated query
export const customInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const source = Axios.CancelToken.source();
  const promise = AXIOS_INSTANCE({
    ...config,
    ...options,
    cancelToken: source.token,
  }).then(({ data }) => data);

  // @ts-ignore
  promise.cancel = () => {
    source.cancel('Query was cancelled');
  };

  return promise;
};

export type ErrorType<T> = AxiosError<T>;