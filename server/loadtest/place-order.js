import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 1,
  duration: '5s',
};

export default function () {
  const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InlvZ2VzaHNoYWgyMDYzQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImlzcyI6Im15LWJhY2tlbmQtc2VydmljZSIsImV4cCI6MTc5MTI3OTQzMywibmJmIjoxNzkxMjc1ODMzLCJpYXQiOjE3OTEyNzU4MzN9._HAttR5cQ9jG5DR6M2D9amvwmfn_wIC9s5Qu3WqM3Rg"

  const res = http.get(
    'http://localhost:8080/api/order/list-user-order',
    {
      headers: {
        'Content-Type': 'application/json',
        "Authorization":`Bearer ${token}`
      },
    }
  );

  console.log(`status=${res.status} body=${res.body}`);

  check(res, {
    'status is 2xx': (r) => r.status >= 200 && r.status < 300,
  });
}