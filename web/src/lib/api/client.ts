import { csrfInterceptor } from './csrf';
import { client } from './gen/client.gen';

client.setConfig({ baseUrl: '', credentials: 'same-origin' });
client.interceptors.request.use(csrfInterceptor);

export { client };
