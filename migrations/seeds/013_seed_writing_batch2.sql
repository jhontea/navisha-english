-- Seed: Additional Writing exercises batch 2

INSERT INTO navisha_english_writing_exercises (title, type, context, prompt, template, key_phrases, difficulty) VALUES

('Slack Message: Requesting Code Review', 'email',
'You are a mid-level developer who has just opened a pull request for a significant feature — a new user notification system. The PR includes 12 files changed, 450 lines added. You need a senior developer (Alex) to review it before it can be merged. You want to be respectful of their time while also communicating the urgency since it is blocking the next sprint.',
'Write a Slack message to your senior developer Alex requesting a code review on your PR. Include: what the PR is about, why it is important, a link placeholder, and your availability to discuss.',
NULL,
ARRAY['would you be able to', 'blocking the next sprint', 'happy to walk you through', 'when you have a moment', 'pull request', 'appreciate your time'],
'B1'),

('Technical Documentation: API Endpoint', 'report',
'You are a backend developer who has just built a new REST API endpoint for user authentication. The endpoint is POST /api/v2/auth/login, accepts email and password in the request body (JSON), returns an access token and refresh token on success, and returns appropriate error codes for invalid credentials or missing fields.',
'Write clear technical documentation for this API endpoint. Include: endpoint details, request format, response format (success and error), and an example request/response.',
'## POST /api/v2/auth/login

**Description:**

### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{

}
```

### Response

**Success (200):**
```json
{

}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| | | |

### Example
',
ARRAY['request body', 'response payload', 'error handling', 'status code', 'authentication', 'access token'],
'B2'),

('Performance Review Self-Assessment', 'report',
'It is annual performance review season. You are a backend developer with 2 years at the company. This year you: led the migration of the payment service to microservices (reduced downtime by 60%), mentored 2 junior developers, improved API response time by 35% through caching optimizations, and missed one sprint deadline due to underestimated complexity of a third-party integration.',
'Write a professional self-assessment for your annual performance review. Cover your key achievements, areas where you grew, one honest challenge you faced, and your goals for next year.',
'# Annual Self-Assessment — [Your Name]
## Review Period: Q1–Q4 2026

## Key Achievements

## Growth & Learning

## Challenges & Lessons Learned

## Goals for Next Year
',
ARRAY['key achievement', 'contributed to', 'I took ownership of', 'areas for improvement', 'going forward', 'I aim to', 'lessons learned'],
'B2'),

('Meeting Agenda Email', 'email',
'You are a tech lead scheduling a technical architecture review meeting for next Tuesday at 2 PM (1 hour). Attendees: frontend lead, backend lead, DevOps engineer, and product manager. Topics: review current system bottlenecks, discuss proposed microservices migration, assign action items.',
'Write a professional meeting invitation email that includes the agenda, expected outcomes, and any pre-reading or preparation attendees should do before the meeting.',
'Subject: Technical Architecture Review — [Date] at 2:00 PM

Hi team,

I would like to invite you to a Technical Architecture Review meeting.

**Date & Time:**
**Location / Link:**
**Duration:**

**Agenda:**
1.
2.
3.

**Expected Outcomes:**

**Preparation:**

Please confirm your attendance by replying to this email.

Best regards,
[Your Name]',
ARRAY['please confirm your attendance', 'agenda items', 'expected outcomes', 'prior to the meeting', 'action items', 'I look forward to'],
'B1');
