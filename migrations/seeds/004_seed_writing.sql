-- Seed: Writing exercises — Email, Proposal, Report

INSERT INTO navisha_english_writing_exercises (title, type, context, prompt, template, key_phrases, difficulty) VALUES

('Bug Report Email to Client', 'email',
'You are a backend developer at a software company. A critical bug was discovered in production that affected some clients'' ability to log in for approximately 2 hours this morning. The issue has been fixed and you need to inform the affected client.',
'Write a professional email to your client (Sarah Chen, IT Director at Nexus Corp) explaining the login issue, what caused it, what you did to fix it, and what steps you''re taking to prevent it from happening again.',
'Subject: [Action Required] Service Disruption Report - Login Issue (June 27, 2026)

Dear [Client Name],

I am writing to inform you of a service disruption that occurred...

Root Cause:

Resolution:

Preventive Measures:

We sincerely apologize for any inconvenience this may have caused...

Best regards,
[Your Name]',
ARRAY['service disruption', 'root cause', 'has been resolved', 'preventive measures', 'we sincerely apologize', 'impact assessment'],
'B1'),

('Technical Proposal: Cloud Migration', 'proposal',
'You are a senior developer presenting a proposal to migrate your company''s monolithic application to a cloud-based microservices architecture. Your audience is the CTO and non-technical managers.',
'Write a short technical proposal (3-4 paragraphs) recommending the migration to AWS microservices. Include the business benefits, estimated timeline, and potential risks with mitigation strategies.',
NULL,
ARRAY['we propose', 'the primary objective', 'in terms of scalability', 'estimated timeline', 'potential risks', 'mitigation strategy', 'return on investment', 'moving forward'],
'B2'),

('Sprint Review Summary', 'report',
'Your team just completed Sprint 14. You shipped: user authentication module, password reset flow, and API rate limiting. One item was not completed: email notification system (moved to next sprint). Two bugs were fixed from the backlog.',
'Write a sprint review summary to share with your team and manager. Include what was completed, what was not completed and why, key metrics, and the plan for the next sprint.',
'Sprint 14 Review Summary
Period: [Date Range]
Team: [Team Name]

## Completed
-

## Not Completed
-

## Key Metrics

## Next Sprint Preview
',
ARRAY['successfully delivered', 'carried over to', 'key achievements', 'impediments encountered', 'going forward', 'acceptance criteria'],
'B1'),

('Responding to Negative Feedback', 'email',
'A client sent an angry email complaining that your team''s latest update broke their reporting dashboard and they lost half a day of productivity. Your team has identified the issue and deployed a fix 30 minutes ago.',
'Write a professional response email to the client (Mark Davidson, Operations Manager). Acknowledge the issue, take responsibility, explain what was done to fix it, and offer a concrete next step.',
NULL,
ARRAY['we understand your frustration', 'take full responsibility', 'has been resolved', 'to prevent recurrence', 'as a next step', 'we value your partnership'],
'B2');
