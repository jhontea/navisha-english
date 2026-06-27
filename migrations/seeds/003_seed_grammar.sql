-- Seed: Grammar exercises — Passive Voice, Modal Verbs, Conditionals

INSERT INTO navisha_english_grammar_exercises (title, topic, instruction, content, difficulty) VALUES

('Passive Voice in Technical Reports', 'Passive Voice', 'Choose the correct passive voice form for each sentence. Focus on how passive voice is used in technical documentation.', '{
  "questions": [
    {
      "id": "q1",
      "text": "The bug ___ by the QA team yesterday.",
      "options": ["was identified", "identified", "has identify", "is identify"],
      "correct_answer": "was identified",
      "explanation": "Use simple past passive (was/were + past participle) for completed actions in the past."
    },
    {
      "id": "q2",
      "text": "The new API endpoints ___ in the next release.",
      "options": ["will be documented", "will document", "are documenting", "has been documented"],
      "correct_answer": "will be documented",
      "explanation": "Use future passive (will be + past participle) for actions planned in the future."
    },
    {
      "id": "q3",
      "text": "The server ___ for maintenance every Sunday night.",
      "options": ["is restarted", "restarts", "has restarted", "was restarting"],
      "correct_answer": "is restarted",
      "explanation": "Use present simple passive for regular, repeated actions or scheduled processes."
    },
    {
      "id": "q4",
      "text": "The requirements ___ with all stakeholders before development began.",
      "options": ["had been reviewed", "have reviewed", "were reviewing", "reviewed"],
      "correct_answer": "had been reviewed",
      "explanation": "Use past perfect passive (had been + past participle) for actions completed before another past action."
    },
    {
      "id": "q5",
      "text": "The deployment ___ since this morning due to a critical issue.",
      "options": ["has been blocked", "is blocking", "was blocked", "blocks"],
      "correct_answer": "has been blocked",
      "explanation": "Use present perfect passive (has/have been + past participle) for actions that started in the past and are still relevant now."
    }
  ]
}', 'B1'),

('Modal Verbs for Professional Emails', 'Modal Verbs', 'Select the most appropriate modal verb for each professional email context. Consider the level of formality and politeness required.', '{
  "questions": [
    {
      "id": "q1",
      "text": "I ___ appreciate it if you could send me the updated project timeline.",
      "options": ["would", "will", "should", "must"],
      "correct_answer": "would",
      "explanation": "\"Would appreciate\" is a polite, formal way to make a request in professional emails."
    },
    {
      "id": "q2",
      "text": "All team members ___ submit their timesheets by end of day Friday.",
      "options": ["must", "might", "could", "would"],
      "correct_answer": "must",
      "explanation": "\"Must\" expresses a strong obligation or requirement, suitable for company policies."
    },
    {
      "id": "q3",
      "text": "___ you please review the attached proposal and share your feedback?",
      "options": ["Could", "Must", "Shall", "Need"],
      "correct_answer": "Could",
      "explanation": "\"Could you please\" is a polite and standard way to make a request in business communication."
    },
    {
      "id": "q4",
      "text": "The system ___ experience brief downtime during the migration window.",
      "options": ["may", "must", "will always", "should never"],
      "correct_answer": "may",
      "explanation": "\"May\" expresses possibility or uncertainty, which is honest and appropriate when informing users of potential issues."
    },
    {
      "id": "q5",
      "text": "You ___ need to restart your browser after the update is applied.",
      "options": ["might", "must always", "will never", "shall"],
      "correct_answer": "might",
      "explanation": "\"Might\" expresses a weaker possibility — appropriate when an action may or may not be necessary."
    }
  ]
}', 'B1'),

('Conditional Sentences for Tech Proposals', 'Conditional Sentences', 'Complete the conditional sentences used in a technical proposal. Choose the correct form for each business context.', '{
  "questions": [
    {
      "id": "q1",
      "text": "If we migrate to microservices, we ___ able to scale each service independently.",
      "options": ["would be", "will been", "are", "were"],
      "correct_answer": "would be",
      "explanation": "Second conditional (if + past simple, would + infinitive) is used for hypothetical proposals or recommendations."
    },
    {
      "id": "q2",
      "text": "If we ___ the caching layer, the API response time will improve significantly.",
      "options": ["implement", "implemented", "had implemented", "would implement"],
      "correct_answer": "implement",
      "explanation": "First conditional (if + present simple, will + infinitive) is used for realistic, likely outcomes."
    },
    {
      "id": "q3",
      "text": "If the client ___ the requirements earlier, we would have met the deadline.",
      "options": ["had confirmed", "confirmed", "confirms", "would confirm"],
      "correct_answer": "had confirmed",
      "explanation": "Third conditional (if + past perfect, would have + past participle) is used to discuss past situations and their hypothetical outcomes."
    },
    {
      "id": "q4",
      "text": "If the team ___ additional training, productivity would increase by 30%.",
      "options": ["received", "receives", "had received", "will receive"],
      "correct_answer": "received",
      "explanation": "Second conditional is used here to make a business case for a hypothetical action and its positive outcome."
    },
    {
      "id": "q5",
      "text": "If the tests pass, the build ___ automatically deployed to staging.",
      "options": ["will be", "would be", "had been", "was"],
      "correct_answer": "will be",
      "explanation": "First conditional is used for automated, rule-based processes like CI/CD pipelines."
    }
  ]
}', 'B2');
