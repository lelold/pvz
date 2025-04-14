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
  const loginRes = http.post('http://localhost:8080/dummyLogin', JSON.stringify({ role: 'moderator' }), {
    headers: { 'Content-Type': 'application/json' },
  });
  const payload = JSON.stringify({
    city: 'Санкт-Петербург'
  });
  const token = loginRes.json()['token'];
  const postPVZRes = http.post('http://localhost:8080/pvz', payload, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  check(postPVZRes, {
    'POST /pvz status is 201': (r) => r.status === 201,
  }); 
}
