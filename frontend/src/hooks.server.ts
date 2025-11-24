import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const handle: Handle = async ({ event, resolve }) => {
  if (event.url.pathname.startsWith('/backend')) {
    const backendURL = `${env.API_BASE}${event.url.pathname.replace('/backend', '')}${event.url.search}`;
    if (['GET', 'HEAD'].includes(event.request.method)) {
      const res = await fetch(backendURL, {
        method: "GET",
        headers: event.request.headers
      });
      return new Response(await res.text(), { headers: { 'Content-Type': 'application/json' } });
    }

    const contentLength = +(event.request.headers.get('Content-Length') || 0);
    const hasBody = contentLength > 0;
    const b = hasBody ? await event.request.text() : undefined;

    const backendResponse = await fetch(backendURL, {
      body: b,
      headers: event.request.headers,
      method: event.request.method
    });

    return new Response(await backendResponse.text(), { headers: { 'Content-Type': 'application/json' } });
  }

  const response = await resolve(event);
  return response;
};
