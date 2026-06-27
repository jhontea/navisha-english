-- Seed: Additional Writing exercises batch 3

INSERT INTO navisha_english_writing_exercises (title, type, context, prompt, template, key_phrases, difficulty) VALUES

('Rejection Email: Vendor Proposal', 'email',
'Your company received a software proposal from a vendor called DataSync Solutions. After evaluation, your team decided not to proceed because the pricing is 40% above budget and the integration with your existing systems would require 3 months of custom development work. You want to decline professionally while leaving the door open for future opportunities.',
'Write a professional rejection email to the account manager at DataSync Solutions (Ryan Mitchell). Decline their proposal clearly, explain the main reasons briefly without being overly negative, and leave the relationship open for the future.',
NULL,
ARRAY['after careful consideration', 'we have decided not to proceed', 'primary concerns', 'does not align with', 'we appreciate your time', 'we will keep your solution in mind', 'should our needs change'],
'B2'),

('Technical Retrospective Report', 'report',
'Your team recently completed a 3-month project to migrate a monolithic e-commerce platform to microservices. The migration went mostly well but had several challenges: two production incidents during migration weekends, the timeline slipped by 3 weeks due to unexpected database compatibility issues, and one team member left mid-project. Positive outcomes: 40% improvement in deployment frequency, 60% reduction in mean time to recovery, and the team learned new skills.',
'Write a technical retrospective report for the engineering leadership team. Cover: project overview, what went well, what did not go well, key learnings, and recommendations for future migrations.',
'# Technical Retrospective: Microservices Migration
## Project Overview

## What Went Well

## Challenges & What Did Not Go Well

## Key Learnings

## Recommendations for Future Projects
',
ARRAY['overall the project', 'key achievement', 'contributed to', 'root cause', 'going forward', 'we recommend', 'lessons learned', 'mitigating factor'],
'C1'),

('Escalation Email: Unresponsive Vendor', 'email',
'You have been trying to get a critical bug fixed by your payment gateway vendor (PayFlow) for 3 weeks. You have sent 4 emails and had 2 support calls, but the issue (intermittent payment failures affecting 2% of transactions) remains unresolved. You need to escalate to their management team.',
'Write a professional but firm escalation email addressed to the VP of Customer Success at PayFlow (Diana Torres). Summarize the issue, provide a timeline of your attempts to resolve it, state the business impact, and request a concrete resolution plan with a deadline.',
'Subject: Escalation: Unresolved Payment Failure Issue — [Reference #]

Dear Ms. Torres,

I am writing to escalate a critical issue that has remained unresolved for three weeks despite multiple attempts to resolve it through your standard support channels.

**Issue Summary:**

**Timeline of Support Interactions:**

**Business Impact:**

**Requested Action:**

I would appreciate your personal attention to this matter and a response with a concrete resolution plan by [Date].

Best regards,
[Your Name]',
ARRAY['I am writing to escalate', 'despite multiple attempts', 'business impact', 'resolution plan', 'concrete timeline', 'I expect a response by', 'this has been ongoing for'],
'C1'),

('LinkedIn Post: Sharing a Technical Achievement', 'email',
'Your team just shipped a major performance improvement: you reduced the API response time of your SaaS platform from an average of 850ms to 120ms by implementing Redis caching, query optimization, and connection pooling. This took 2 sprints and involved 3 engineers. You want to share this on LinkedIn to build your professional profile.',
'Write a LinkedIn post (150-200 words) sharing this technical achievement. Make it engaging and informative for both technical and non-technical audiences. Include what the problem was, what you did, the result, and a key learning.',
NULL,
ARRAY['excited to share', 'the challenge was', 'we tackled this by', 'the result', 'key takeaway', 'shoutout to the team', 'if you are working on similar challenges'],
'B1');
