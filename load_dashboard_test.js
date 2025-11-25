import http from 'k6/http';
import { sleep } from 'k6';

// Configure load pattern
export const options = {
  vus: 1000,          // 1000 concurrent virtual users
  duration: '30s',    // run for 30 seconds
};

// Basic exponential backoff retry function
function requestWithRetry(url, maxRetries = 3, baseWait = 1) {
  let attempt = 0;
  while (attempt < maxRetries) {
    const res = http.get(url);

    // success
    if (res.status === 200) {
      return res;
    }

    // If 429, retry with backoff
    if (res.status === 429) {
      const waitTime = baseWait * Math.pow(2, attempt); // 1s, 2s, 4s, ...
      console.warn(`Got 429 (Too Many Requests) — retrying in ${waitTime}s...`);
      sleep(waitTime);
      attempt++;
      continue;
    }

    // For other errors, no retry
    console.error(`Request failed: ${res.status}`);
    return res;
  }

  console.error(`Gave up after ${maxRetries} retries`);
  return { status: 429 };
}

export default function () {
  const res = requestWithRetry('http://172.16.2.19:30080/buy', 3, 1);
  sleep(1);
}
