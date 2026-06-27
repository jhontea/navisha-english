-- Seed: Additional Vocabulary batch 2 — IT Business English (28 more words)

INSERT INTO navisha_english_vocabulary (word, definition, category, example_sentence, difficulty) VALUES

-- Project Management (additional)
('stakeholder alignment', 'The process of ensuring all key parties share a common understanding and agreement on goals', 'Project Management', 'We need stakeholder alignment before we commit to the new architecture approach.', 'B2'),
('change request', 'A formal proposal to modify a project scope, schedule, or budget', 'Project Management', 'The client submitted a change request to add two-factor authentication to the scope.', 'B1'),
('blocker removal', 'The act of identifying and resolving obstacles that prevent team progress', 'Project Management', 'Blocker removal is the Scrum Master''s top priority during the sprint.', 'B1'),
('acceptance criteria', 'The conditions a product must meet to be accepted by the stakeholder or client', 'Project Management', 'Before we start development, let''s define the acceptance criteria for each user story.', 'B2'),
('velocity', 'A measure of the amount of work a team can complete in a single sprint', 'Project Management', 'Our team''s velocity has increased from 30 to 42 story points over the last three sprints.', 'B2'),
('handover', 'The process of transferring responsibility for a project or task to another person or team', 'Project Management', 'Please prepare thorough documentation before the handover to the maintenance team.', 'B1'),
('kickoff meeting', 'The first meeting of a project team to align on goals, roles, and timeline', 'Project Management', 'We''ll hold a kickoff meeting with the client on Monday to set expectations for the project.', 'B1'),

-- Technical Communication (additional)
('API contract', 'A formal agreement between teams defining how an API will behave, its endpoints, and data formats', 'Technical Communication', 'We need to finalize the API contract before the frontend team starts integration.', 'B2'),
('hot fix', 'An urgent fix deployed directly to production to resolve a critical issue', 'Technical Communication', 'We deployed a hot fix to resolve the payment gateway timeout affecting customers.', 'B1'),
('feature flag', 'A technique that allows features to be enabled or disabled without deploying new code', 'Technical Communication', 'We''ll use a feature flag to gradually roll out the new dashboard to users.', 'B2'),
('code freeze', 'A period during which no new code changes are allowed before a major release', 'Technical Communication', 'Code freeze starts Friday — make sure all your PRs are merged before then.', 'B2'),
('smoke test', 'A preliminary test to check that the most critical functions of a system work correctly', 'Technical Communication', 'After deployment, run a smoke test to confirm the login and checkout flows are working.', 'B2'),
('payload', 'The data transmitted in a request or response body', 'Technical Communication', 'The API request payload must include the user ID and the action type.', 'B1'),
('endpoint', 'A specific URL where an API can be accessed to perform a particular action', 'Technical Communication', 'The /api/v2/users endpoint now returns paginated results with a default limit of 20.', 'B1'),

-- Meeting & Email Phrases (additional)
('reach out', 'To contact someone for information, help, or collaboration', 'Meeting & Email Phrases', 'Please reach out to the DevOps team if you encounter any issues during deployment.', 'B1'),
('on the same page', 'Having a shared understanding of the situation or plan', 'Meeting & Email Phrases', 'Let''s schedule a quick sync to make sure everyone is on the same page before the sprint starts.', 'B1'),
('take it to the next level', 'To improve or escalate something significantly', 'Meeting & Email Phrases', 'The new CI/CD pipeline really takes our deployment process to the next level.', 'B2'),
('quick win', 'A small, easily achievable improvement that delivers immediate value', 'Meeting & Email Phrases', 'Caching the database queries is a quick win that could reduce load times by 40%.', 'B1'),
('deep dive', 'A thorough and detailed investigation or discussion of a topic', 'Meeting & Email Phrases', 'Let''s schedule a deep dive into the security architecture with the infra team next week.', 'B2'),
('sign off', 'To give formal approval or agreement to something', 'Meeting & Email Phrases', 'I need the product manager to sign off on the updated wireframes before we start coding.', 'B1'),
('low-hanging fruit', 'Tasks or improvements that are easy to achieve and deliver quick results', 'Meeting & Email Phrases', 'Let''s address the low-hanging fruit items first before tackling the complex refactoring.', 'B2'),

-- DevOps & Cloud (additional)
('blue-green deployment', 'A deployment strategy using two identical environments to minimize downtime during releases', 'DevOps & Cloud', 'We use blue-green deployment to ensure zero downtime when releasing new backend versions.', 'C1'),
('canary release', 'Gradually rolling out a change to a small subset of users before a full release', 'DevOps & Cloud', 'We deployed the new recommendation engine as a canary release to 5% of users first.', 'C1'),
('autoscaling', 'Automatically adjusting computing resources based on current demand', 'DevOps & Cloud', 'With autoscaling enabled, our servers handled the Black Friday traffic spike without issues.', 'B2'),
('secrets management', 'Securely storing and managing sensitive credentials and API keys', 'DevOps & Cloud', 'Never hardcode credentials — use a secrets management tool like Vault or AWS Secrets Manager.', 'B2'),
('observability', 'The ability to understand a system''s internal state from its outputs (logs, metrics, traces)', 'DevOps & Cloud', 'Good observability helped us pinpoint the exact microservice causing the latency spike.', 'C1'),
('on-call', 'A rotation where team members are available to respond to production incidents outside business hours', 'DevOps & Cloud', 'As the on-call engineer this week, you need to respond to PagerDuty alerts within 15 minutes.', 'B1'),
('runbook', 'A documented set of procedures for handling operational tasks or incidents', 'DevOps & Cloud', 'We updated the runbook with the steps to follow when the payment service goes down.', 'B2');
