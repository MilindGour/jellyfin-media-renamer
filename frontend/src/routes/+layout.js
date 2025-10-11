import { v4 as uuidV4 } from 'uuid';
import { Log } from '$lib/services/logger';

const uuid = uuidV4();
const ws = new WebSocket(`ws://localhost:7749/api/ws/${uuid}`);
const log = new Log('WebSocket');

ws.onopen = () => {
	log.info('Connection opened');
};

ws.onclose = () => {
	log.info('Connection closed');
};

/** @type {(e: any) => void} */
ws.onmessage = (e) => {
	log.info('Message received:', e);
};
