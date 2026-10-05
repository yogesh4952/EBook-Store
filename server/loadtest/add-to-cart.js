import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 1,
  duration: '5s',
};

export default function () {
  const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InlvZ2VzaHNoYWgyMDYzQGdtYWlsLmNvbSIsInJvbGUiOiJjdXN0b21lciIsImlzcyI6Im15LWJhY2tlbmQtc2VydmljZSIsImV4cCI6MTc5MTE4Nzk3OSwibmJmIjoxNzkxMTg0Mzc5LCJpYXQiOjE3OTExODQzNzl9.By0F2KrLxzFdAjG_atK4r-tudVH4PChM0ULHymCcZqM"

  const res = http.post(
    'http://localhost:8080/api/cart/add-to-cart',
    JSON.stringify({
      book_id: 1,
      quantity: 1,
    }),
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