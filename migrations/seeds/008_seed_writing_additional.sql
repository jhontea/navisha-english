-- Seed: Additional Writing exercises

INSERT INTO navisha_english_writing_exercises (title, type, context, prompt, template, key_phrases, difficulty) VALUES

('Request for Project Extension', 'email',
'You are a developer working on a 6-week project to build a data analytics dashboard for a client. You are now in week 5 and have discovered that integrating with the client''s legacy data warehouse is significantly more complex than estimated. You need 2 extra weeks to complete the project properly.',
'Write a professional email to your project manager (David Park) requesting a 2-week extension. Explain the technical reason clearly without using too much jargon, propose a revised timeline, and suggest what can be delivered by the original deadline.',
NULL,
ARRAY['I am writing to request', 'unforeseen technical challenges', 'revised timeline', 'interim deliverable', 'to ensure the quality', 'I apologize for any inconvenience'],
'B1'),

('Onboarding Guide for New Team Member', 'report',
'Your team has just hired a junior backend developer who will join in 2 weeks. Your manager has asked you to write a short onboarding guide covering the tech stack, development workflow, and key tools they need to know.',
'Write a concise onboarding guide (3-4 sections) for the new developer. Cover: tech stack overview, local development setup process, Git workflow, and team communication tools.',
'# Developer Onboarding Guide

Welcome to the team! This guide will help you get up and running quickly.

## Tech Stack

## Local Development Setup

## Git Workflow

## Communication & Tools

If you have any questions, don''t hesitate to reach out!',
ARRAY['tech stack', 'development environment', 'pull request', 'code review process', 'point of contact', 'best practices'],
'B1'),

('Incident Report: System Outage', 'report',
'Your company''s production API experienced a 45-minute outage yesterday (June 26, 2026, 14:15–15:00 UTC). Root cause: a misconfigured environment variable was pushed to production during a routine deployment. Impact: approximately 1,200 users could not access the service. Fix: the variable was corrected and a rollback was performed.',
'Write a formal incident report for your engineering manager and the affected stakeholders. Include: incident summary, timeline, root cause, impact, resolution, and preventive measures.',
'# Incident Report — API Outage
Date: June 26, 2026
Severity: High
Status: Resolved

## Summary

## Timeline
| Time (UTC) | Event |
|------------|-------|
| | |

## Root Cause

## Impact

## Resolution

## Preventive Measures
',
ARRAY['incident summary', 'root cause analysis', 'time of detection', 'resolution steps', 'preventive measures', 'lessons learned', 'affected users'],
'B2'),

('Slack Message: Announcing a Breaking Change', 'email',
'You are a backend engineer. You are about to deprecate the v1 /users endpoint and replace it with /v2/users which has a different response format. The change will go live in 2 weeks. You need to notify all frontend and mobile developers in the #engineering Slack channel.',
'Write a clear, professional Slack message announcing the breaking API change. Include what is changing, why it is changing, the timeline, what action developers need to take, and who to contact with questions.',
NULL,
ARRAY['breaking change', 'deprecated', 'migration guide', 'action required', 'backward compatibility', 'please update', 'reach out to'],
'B2'),

('Weekly Status Update Email', 'email',
'It is Friday afternoon. You are a senior developer working on a 3-month platform migration project. This week: completed database schema migration (ahead of schedule), started API endpoint refactoring (50% done), blocker: waiting for security team approval on new auth flow. Next week plan: finish API refactoring, start frontend integration.',
'Write a concise weekly status update email to your project manager and stakeholders. Keep it professional, scannable, and under 200 words.',
'Subject: Weekly Status Update — Platform Migration — Week 8

Hi [Manager],

Here is this week''s update on the platform migration project.

**Completed This Week**
-

**In Progress**
-

**Blockers**
-

**Plan for Next Week**
-

Overall the project remains [on track / slightly behind / ahead of schedule].

Best regards,
[Your Name]',
ARRAY['on track', 'completed ahead of schedule', 'pending approval', 'key blocker', 'next steps', 'overall status'],
'B1');
