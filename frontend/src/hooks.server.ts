import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const handle: Handle = async ({ event, resolve }) => {
  if (event.url.pathname.startsWith('/backend')) {
    const { pathname, search } = event.url;
    const { method, headers, json } = event.request;
    const contentLength = +(headers.get('Content-Length') || 0);
    const hasBody = contentLength > 0;
    const body = hasBody && await json() || undefined;

    const backendURL = `${env.API_BASE}${pathname.replace('/backend', '')}${search}`;

    const backendResponse = await fetch(backendURL, {
      body,
      headers,
      method
    });

    return new Response(await backendResponse.text(), { headers: { 'Content-Type': 'application/json' } });
  }

  const response = await resolve(event);
  return response;
};
