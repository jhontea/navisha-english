-- Seed: Additional Vocabulary — IT Business English (32 more words)

INSERT INTO navisha_english_vocabulary (word, definition, category, example_sentence, difficulty) VALUES

-- Project Management (additional)
('capacity planning', 'The process of determining the production capacity needed to meet changing demands', 'Project Management', 'We need to do capacity planning before committing to the new client project.', 'B2'),
('dependencies', 'Tasks or items that rely on other tasks being completed first', 'Project Management', 'Please list all dependencies before we finalize the sprint backlog.', 'B1'),
('sign-off', 'Formal approval or acceptance of a deliverable or decision', 'Project Management', 'We need sign-off from the client before we can move to the next phase.', 'B1'),
('turnaround time', 'The time required to complete a process or fulfill a request', 'Project Management', 'Our average turnaround time for bug fixes is 48 hours.', 'B1'),
('escalate', 'To raise an issue to a higher level of authority for resolution', 'Project Management', 'If the blocker isn''t resolved by tomorrow, we''ll need to escalate it to the director.', 'B2'),
('risk mitigation', 'Actions taken to reduce the probability or impact of a risk', 'Project Management', 'The risk mitigation plan includes daily backups and a failover server.', 'B2'),
('go-live', 'The moment when a system or product becomes available to end users', 'Project Management', 'The go-live date for the new payment system is scheduled for next Monday.', 'B1'),
('post-mortem', 'An analysis conducted after a project or incident to identify what went wrong', 'Project Management', 'We''ll hold a post-mortem meeting on Friday to review the production outage.', 'B2'),

-- Technical Communication (additional)
('proof of concept', 'A prototype or demo built to test the feasibility of an idea', 'Technical Communication', 'We built a proof of concept to validate the AI recommendation engine before full development.', 'B2'),
('technical debt', 'The accumulated cost of shortcuts and suboptimal solutions in a codebase', 'Technical Communication', 'We need to allocate time in the next sprint to address our growing technical debt.', 'B2'),
('edge case', 'An unusual or extreme scenario that occurs at the boundary of operating parameters', 'Technical Communication', 'The QA team found an edge case where the app crashes with special characters in the username.', 'B2'),
('throughput', 'The amount of work or data processed by a system in a given period', 'Technical Communication', 'After optimization, our API throughput increased from 500 to 2,000 requests per second.', 'C1'),
('idempotent', 'An operation that produces the same result regardless of how many times it is performed', 'Technical Communication', 'Make sure the payment endpoint is idempotent to prevent duplicate charges.', 'C1'),
('regression', 'A bug introduced by a recent change that breaks previously working functionality', 'Technical Communication', 'The latest deployment caused a regression in the login flow for mobile users.', 'B2'),
('load balancing', 'Distributing incoming network traffic across multiple servers to ensure reliability', 'Technical Communication', 'We implemented load balancing to handle traffic spikes during peak hours.', 'B2'),
('codebase', 'The complete body of source code for a software project', 'Technical Communication', 'New developers should take time to explore the codebase before making changes.', 'B1'),

-- Meeting & Email Phrases (additional)
('as per our discussion', 'Referring to what was agreed or said in a previous conversation', 'Meeting & Email Phrases', 'As per our discussion last Tuesday, I''m sending over the revised project timeline.', 'B1'),
('going forward', 'From now on; used to describe a new approach or process', 'Meeting & Email Phrases', 'Going forward, all API changes must be documented in the changelog before merging.', 'B1'),
('at your earliest convenience', 'As soon as you are able to, without rushing', 'Meeting & Email Phrases', 'Could you please review the attached proposal at your earliest convenience?', 'B2'),
('in the pipeline', 'Being planned or in progress, but not yet completed', 'Meeting & Email Phrases', 'We have three new features in the pipeline for the next quarter.', 'B1'),
('bandwidth permitting', 'If time and resources allow', 'Meeting & Email Phrases', 'Bandwidth permitting, we''d like to tackle the UI redesign this sprint.', 'B2'),
('table this', 'To postpone discussion of something to a later time', 'Meeting & Email Phrases', 'Let''s table this topic for now and revisit it in next week''s meeting.', 'B2'),
('align on', 'To reach a shared understanding or agreement on something', 'Meeting & Email Phrases', 'We need to align on the API contract before the frontend team starts building.', 'B2'),
('FYI', 'For Your Information — used to share information without requiring action', 'Meeting & Email Phrases', 'FYI, the staging environment will be unavailable this afternoon for maintenance.', 'B1'),

-- DevOps & Cloud (new category)
('CI/CD', 'Continuous Integration and Continuous Deployment — automated pipeline for building, testing, and deploying code', 'DevOps & Cloud', 'Our CI/CD pipeline runs automated tests on every pull request before allowing a merge.', 'B1'),
('containerization', 'Packaging an application and its dependencies into a portable container', 'DevOps & Cloud', 'We use Docker for containerization to ensure consistent environments across dev and production.', 'B2'),
('downtime', 'A period when a system is unavailable or not functioning', 'DevOps & Cloud', 'The planned maintenance window will cause approximately 30 minutes of downtime.', 'B1'),
('infrastructure as code', 'Managing and provisioning infrastructure through machine-readable configuration files', 'DevOps & Cloud', 'We use Terraform for infrastructure as code, so all cloud resources are version-controlled.', 'C1'),
('zero downtime deployment', 'A deployment strategy that keeps the application available throughout the update process', 'DevOps & Cloud', 'We implemented blue-green deployments to achieve zero downtime deployment for critical services.', 'C1'),
('monitoring', 'The process of continuously observing a system to detect issues and measure performance', 'DevOps & Cloud', 'Our monitoring dashboard alerts the on-call team within 1 minute of any service degradation.', 'B1'),
('SLA', 'Service Level Agreement — a contract defining the expected level of service between provider and client', 'DevOps & Cloud', 'Our SLA guarantees 99.9% uptime, which means no more than 8.7 hours of downtime per year.', 'B2'),
('environment', 'A distinct setup where software runs, such as development, staging, or production', 'DevOps & Cloud', 'Never test unreviewed changes directly in the production environment.', 'B1');
