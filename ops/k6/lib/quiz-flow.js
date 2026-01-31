import http from 'k6/http';
import { check, sleep } from 'k6';

/** ボトルネック分析用: エンドポイントごとに name タグを付与し、結果の http_req_duration をエンドポイント別に集計する */
const TAGS = {
  questions: { name: 'quiz_questions' },
  answer: { name: 'quiz_answer' },
};

/**
 * クイズ1回分のフロー: 問題取得 → 各問に回答
 *
 * ボトルネック確認: k6 実行後、サマリーの http_req_duration が
 * name タグ別（quiz_questions / quiz_answer）に表示される。
 * 遅い方のエンドポイントがボトルネック候補。
 *
 * @param {string} baseUrl - API のベースURL（例: http://localhost:8080）
 */
export function runQuizFlow(baseUrl) {
  const questionsRes = http.get(`${baseUrl}/api/v1/quiz/questions`, { tags: TAGS.questions });
  check(questionsRes, { 'questions status is 200': (r) => r.status === 200 });
  if (questionsRes.status !== 200) {
    return;
  }

  const questions = questionsRes.json().questions || [];
  if (questions.length === 0) {
    return;
  }

  for (const q of questions) {
    const answerPayload = JSON.stringify({
      question_id: q.id,
      selected_category: q.choices && q.choices[0] ? q.choices[0] : '可燃',
    });
    const answerRes = http.post(`${baseUrl}/api/v1/quiz/answer`, answerPayload, {
      headers: { 'Content-Type': 'application/json' },
      tags: TAGS.answer,
    });
    check(answerRes, { 'answer status is 200': (r) => r.status === 200 });
    sleep(0.2);
  }
}
