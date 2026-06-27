-- Seed: Additional Vocabulary batch 3 — Advanced IT Business English (30 words)

INSERT INTO navisha_english_vocabulary (word, definition, category, example_sentence, difficulty) VALUES

-- Project Management batch 3
('resource allocation', 'The process of assigning available resources to tasks and projects efficiently', 'Project Management', 'Poor resource allocation caused three developers to be overbooked during the same sprint.', 'B2'),
('deliverable sign-off', 'Formal client or stakeholder approval that a deliverable meets agreed requirements', 'Project Management', 'We cannot move to phase two until we receive deliverable sign-off from the client.', 'B2'),
('proof of value', 'A demonstration that a solution delivers measurable business benefit', 'Project Management', 'The pilot program served as a proof of value before the company committed to full deployment.', 'C1'),
('escalation path', 'A defined sequence of people or teams to contact when an issue cannot be resolved at the current level', 'Project Management', 'If the bug is not fixed in 24 hours, follow the escalation path and notify the VP of Engineering.', 'B2'),
('hard deadline', 'A fixed, non-negotiable deadline that cannot be moved', 'Project Management', 'The regulatory submission is a hard deadline — missing it has serious legal consequences.', 'B1'),
('parking lot', 'A list of topics raised during a meeting that are deferred for later discussion', 'Project Management', 'That''s a great point — let''s put it in the parking lot and revisit it in next week''s meeting.', 'B2'),
('Definition of Done', 'A shared understanding of what criteria must be met before a task is considered complete', 'Project Management', 'According to our Definition of Done, code must pass all tests and be peer-reviewed before merging.', 'B2'),

-- Technical Communication batch 3
('dependency injection', 'A design pattern where dependencies are provided to a component rather than created by it', 'Technical Communication', 'We use dependency injection to make the service layer easier to test and maintain.', 'C1'),
('race condition', 'A bug that occurs when the outcome depends on the unpredictable timing of multiple processes', 'Technical Communication', 'The data corruption was caused by a race condition in the concurrent write operations.', 'C1'),
('event-driven architecture', 'A software design pattern where components communicate through events', 'Technical Communication', 'We migrated to an event-driven architecture to decouple the order service from the inventory system.', 'C1'),
('schema migration', 'A controlled change to a database schema, typically versioned and tracked', 'Technical Communication', 'Run the schema migration before deploying the new version, or the app will crash on startup.', 'B2'),
('graceful degradation', 'The ability of a system to continue functioning in a reduced capacity when part of it fails', 'Technical Communication', 'The app uses graceful degradation — if the recommendations service is down, it shows generic content instead.', 'C1'),
('retry logic', 'Code that automatically re-attempts a failed operation a set number of times', 'Technical Communication', 'The API client has retry logic with exponential backoff to handle transient network failures.', 'B2'),
('idempotency key', 'A unique identifier sent with a request to ensure the same operation is not performed twice', 'Technical Communication', 'Always include an idempotency key when making payment requests to prevent duplicate charges.', 'C1'),

-- Meeting & Email Phrases batch 3
('for your awareness', 'Used to share information that does not require action but the recipient should know', 'Meeting & Email Phrases', 'For your awareness, the staging environment will be offline this weekend for infrastructure upgrades.', 'B2'),
('as discussed', 'Referring to something previously agreed or mentioned in a meeting or conversation', 'Meeting & Email Phrases', 'As discussed in today''s standup, I will take ownership of the API integration task.', 'B1'),
('pending your approval', 'Waiting for someone''s formal agreement before proceeding', 'Meeting & Email Phrases', 'The deployment is ready and pending your approval to proceed to production.', 'B2'),
('revisit', 'To return to a topic or decision at a later time for review', 'Meeting & Email Phrases', 'Let''s revisit the pricing model in Q3 once we have more usage data from customers.', 'B1'),
('keep me in the loop', 'A request to be kept informed about progress or developments', 'Meeting & Email Phrases', 'Please keep me in the loop as you work through the integration issues with the vendor.', 'B1'),
('no-brainer', 'A decision or choice that is obviously correct and requires little thought', 'Meeting & Email Phrases', 'Switching to automated testing is a no-brainer — it will save us hours of manual QA every sprint.', 'B2'),
('food for thought', 'Something worth thinking about or considering further', 'Meeting & Email Phrases', 'Here''s some food for thought: what if we built the feature as a plugin rather than core functionality?', 'B2'),

-- DevOps & Cloud batch 3
('immutable infrastructure', 'An approach where servers are never modified after deployment — they are replaced entirely', 'DevOps & Cloud', 'We follow immutable infrastructure principles, so every deployment creates fresh containers from scratch.', 'C1'),
('service mesh', 'A dedicated infrastructure layer for managing service-to-service communication in microservices', 'DevOps & Cloud', 'We use Istio as our service mesh to handle load balancing and mTLS between microservices.', 'C1'),
('mean time to recovery', 'The average time it takes to restore a system to full operation after a failure (MTTR)', 'DevOps & Cloud', 'Our MTTR improved from 45 minutes to under 10 minutes after implementing automated rollback.', 'C1'),
('health check', 'An automated test that verifies a service is running and responding correctly', 'DevOps & Cloud', 'The load balancer runs a health check every 30 seconds and removes unhealthy instances automatically.', 'B1'),
('rate limiting', 'Controlling how many requests a client can make to an API within a given time period', 'DevOps & Cloud', 'We implemented rate limiting to prevent a single client from overloading the API with requests.', 'B2'),
('multi-tenancy', 'A software architecture where a single instance serves multiple customers (tenants)', 'DevOps & Cloud', 'Our SaaS platform uses multi-tenancy, meaning all customers share the same infrastructure but with isolated data.', 'B2'),
('disaster recovery', 'A set of policies and procedures to recover from a catastrophic system failure', 'DevOps & Cloud', 'Our disaster recovery plan ensures we can restore all customer data within 4 hours of a major outage.', 'B2'),
('egress cost', 'The cost charged by cloud providers for data transferred out of their network', 'DevOps & Cloud', 'Egress costs were unexpectedly high because we were transferring large video files across regions unnecessarily.', 'C1'),
('WAF', 'Web Application Firewall — a security layer that filters and monitors HTTP traffic to protect web applications', 'DevOps & Cloud', 'The WAF blocked the SQL injection attempt before it reached our application servers.', 'B2');
