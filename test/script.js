import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages:[
    {duration: '30s', target: 20},
    {duration: '1m30s', target: 70},
    {duration: '20s', target: 0},
  ]
};

export default function() {
  http.get('http://ec2-54-254-14-197.ap-southeast-1.compute.amazonaws.com:8080/v1/auth/login');
  sleep(1);
}
