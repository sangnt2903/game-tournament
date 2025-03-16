import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 200 },   // Ramp-up: 200 VUs in 30 sec
        { duration: '20s', target: 1000 },   // Peak load: 1000 VUs for 1 min
        { duration: '20s', target: 3000 },   // Peak load: 1000 VUs for 1 min
        { duration: '1m', target: 10000 },   // Peak load: 1000 VUs for 1 min
        { duration: '20s', target: 3000 },   // Peak load: 1000 VUs for 1 min
        { duration: '20s', target: 1000 },   // Peak load: 1000 VUs for 1 min
        { duration: '10s', target: 200 },   // Ramp-down: Back to 200 VUs in 30 sec
        { duration: '10s', target: 0 },     // Cool-down: Gradual stop
    ],
};

export default function () {
    let userId = `user-${__VU}`;  // Unique user ID per VU

    let params = {
        headers: {
            "Content-Type": "application/json",
            "username": userId,
        },
    };

    let res = http.get('http://localhost:9000/leaderboard?username=' + userId + '&top=100', params);

    check(res, {
        "Status is 200": (r) => r.status === 200,
        "Response time < 200ms": (r) => r.timings.duration < 200,
    });

    sleep(0.01); // 10ms sleep before next request
}
