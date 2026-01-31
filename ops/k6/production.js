/**
 * 段階的負荷（本番想定）
 * 100人 5分（各段階の前に30秒で増加）
 *
 * ボトルネック分析: 実行後サマリーで http_req_duration が
 * quiz_questions / quiz_answer 別に表示される。
 * 出力: summary.html（HTML）, summary.json（JSON）, コンソール（テキスト）
 * パス変更: SUMMARY_HTML=path, SUMMARY_JSON=path で指定可
 */
import { sleep } from 'k6';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';
import { htmlReport } from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js';
import { runQuizFlow } from './lib/quiz-flow.js';

const BASE_URL = 'https://api.garbage-category.stg.dosukoi.click';

export const options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '5m', target: 100 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.05'],
  },
};

export default function () {
  runQuizFlow(BASE_URL);
  sleep(0.5);
}

export function handleSummary(data) {
  const htmlPath = __ENV.SUMMARY_HTML || 'summary.html';
  return {
    [htmlPath]: htmlReport(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}
