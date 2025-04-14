import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 100,
  duration: '10s',
  thresholds: {
    http_req_duration: ['p(95)<100'],
    http_req_failed: ['rate<0.0001'],
  },
};

export default function () {
  const loginRes = http.post('http://localhost:8080/dummyLogin', JSON.stringify({ role: 'employee' }), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginRes, {
    'login succeeded': (r) => r.status === 200,
  });
}

