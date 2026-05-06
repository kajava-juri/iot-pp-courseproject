 // Using the .env value if exists, otherwise call a default server

import axios, { type AxiosResponse } from "axios";

// Credit: https://www.reddit.com/r/reactjs/comments/11iuryi/how_does_a_proper_fetch_wrapper_look/
const baseUrl = import.meta.env.VITE_SERVER_URL || 'http://localhost';
const port = import.meta.env.VITE_WEB_PORT || '8080';
const fullBaseUrl = `${baseUrl}:${port}`;

export const fetcher: (query: string) => Promise<AxiosResponse> = (query) => {
  const request = axios.get(fullBaseUrl + query);
  return request
    .then((response) => response.data)
    .catch((error) => {
      if (error.response) {
        // The request was made and the server responded with a status code
        // that falls out of the range of 2xx
        console.log(`Error response status: ${error.response.status}`);
        console.log(`Error message data: `, error.response.data);
        console.log(`Error response headers: `, error.response.headers);
      } else if (error.request) {
        // The request was made but no response was received
        // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
        // http.ClientRequest in node.js
        console.log(`Error Request:`, error.request);
      } else {
        // Something happened in setting up the request that triggered an Error
        console.log('Error', error.message);
      }
      console.log(`Error config: `, error.config);
    });
};