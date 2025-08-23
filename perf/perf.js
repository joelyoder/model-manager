import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const importPayload = open('./sample-import.json');

export const options = {
  scenarios: {
    models: {
      executor: 'constant-arrival-rate',
      exec: 'models',
      rate: 20,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 20,
    },
    sync: {
      executor: 'constant-arrival-rate',
      exec: 'sync',
      rate: 5,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 5,
    },
    importEndpoint: {
      executor: 'constant-arrival-rate',
      exec: 'importEndpoint',
      rate: 2,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 2,
    },
  },
  thresholds: {
    'http_req_duration{scenario:models}': ['p(95)<500'],
    'http_req_duration{scenario:sync}': ['p(95)<1000'],
    'http_req_duration{scenario:importEndpoint}': ['p(95)<1000'],
  },
};

export function models() {
  const res = http.get(`${BASE_URL}/api/models`);
  check(res, { 'status was 200': r => r.status === 200 });
}

export function sync() {
  const res = http.post(`${BASE_URL}/api/sync`, null, { headers: { 'Content-Type': 'application/json' } });
  check(res, { 'status was 200': r => r.status === 200 });
}

export function importEndpoint() {
  const data = { file: http.file(importPayload, 'import.json', 'application/json') };
  const res = http.post(`${BASE_URL}/api/import`, data);
  check(res, { 'status was 200': r => r.status === 200 });
}
