-- Seed: Additional Grammar exercises — Reported Speech, Article Usage, Tense Review

INSERT INTO navisha_english_grammar_exercises (title, topic, instruction, content, difficulty) VALUES

('Reported Speech in Meeting Notes', 'Reported Speech', 'Convert the direct speech from a meeting into reported speech for meeting minutes. Choose the correct reported speech form.', '{
  "questions": [
    {
      "id": "q1",
      "text": "Direct: \"We will release the feature next week.\" Reported: The project manager said that they ___ the feature the following week.",
      "options": ["would release", "will release", "released", "had released"],
      "correct_answer": "would release",
      "explanation": "In reported speech, \"will\" becomes \"would\" when the reporting verb is in the past tense."
    },
    {
      "id": "q2",
      "text": "Direct: \"The client has approved the design.\" Reported: She confirmed that the client ___ the design.",
      "options": ["had approved", "has approved", "approved", "would approve"],
      "correct_answer": "had approved",
      "explanation": "In reported speech, present perfect (has approved) shifts back to past perfect (had approved)."
    },
    {
      "id": "q3",
      "text": "Direct: \"We are working on the bug fix right now.\" Reported: The developer mentioned that they ___ on the bug fix at that moment.",
      "options": ["were working", "are working", "worked", "had worked"],
      "correct_answer": "were working",
      "explanation": "Present continuous (are working) shifts to past continuous (were working) in reported speech."
    },
    {
      "id": "q4",
      "text": "Direct: \"Can you send me the report by Friday?\" Reported: The manager asked if I ___ send the report by Friday.",
      "options": ["could", "can", "would", "should"],
      "correct_answer": "could",
      "explanation": "\"Can\" becomes \"could\" in reported questions. The word order also changes to statement form."
    },
    {
      "id": "q5",
      "text": "Direct: \"Do not push directly to the main branch.\" Reported: The tech lead told the team ___ push directly to the main branch.",
      "options": ["not to", "to not", "do not", "not"],
      "correct_answer": "not to",
      "explanation": "Negative imperatives in reported speech use \"not to + infinitive\"."
    }
  ]
}', 'B2'),

('Articles in Technical Writing', 'Articles', 'Choose the correct article (a, an, the, or no article) for each gap in these IT business sentences.', '{
  "questions": [
    {
      "id": "q1",
      "text": "We need to schedule ___ meeting with the DevOps team to discuss the deployment.",
      "options": ["a", "an", "the", "no article"],
      "correct_answer": "a",
      "explanation": "Use \"a\" before singular countable nouns when mentioned for the first time. \"Meeting\" starts with a consonant sound."
    },
    {
      "id": "q2",
      "text": "Please review ___ API documentation before the integration call tomorrow.",
      "options": ["the", "a", "an", "no article"],
      "correct_answer": "the",
      "explanation": "Use \"the\" when both speaker and listener know which specific thing is being referred to — in this case, the specific API documentation for their project."
    },
    {
      "id": "q3",
      "text": "She has ___ excellent understanding of cloud architecture.",
      "options": ["an", "a", "the", "no article"],
      "correct_answer": "an",
      "explanation": "Use \"an\" before words that start with a vowel sound. \"Excellent\" starts with the vowel sound /e/."
    },
    {
      "id": "q4",
      "text": "___ software engineers are in high demand across all industries.",
      "options": ["no article", "The", "A", "An"],
      "correct_answer": "no article",
      "explanation": "No article is used when making a general statement about a category of people or things in the plural."
    },
    {
      "id": "q5",
      "text": "The team deployed ___ hotfix to resolve the critical issue in production.",
      "options": ["a", "an", "the", "no article"],
      "correct_answer": "a",
      "explanation": "Use \"a\" for singular countable nouns mentioned for the first time. A hotfix is one of many possible hotfixes."
    }
  ]
}', 'B1'),

('Present Perfect vs Simple Past in Business Emails', 'Tense', 'Choose the correct tense — present perfect or simple past — for each gap in these professional email sentences.', '{
  "questions": [
    {
      "id": "q1",
      "text": "I ___ the report you requested and have a few questions.",
      "options": ["have reviewed", "reviewed", "was reviewing", "had reviewed"],
      "correct_answer": "have reviewed",
      "explanation": "Use present perfect when the action is recent and relevant to the present moment. The result (questions) is still current."
    },
    {
      "id": "q2",
      "text": "The team ___ the new authentication system last Thursday.",
      "options": ["deployed", "has deployed", "had deployed", "was deploying"],
      "correct_answer": "deployed",
      "explanation": "Use simple past for completed actions at a specific time in the past (\"last Thursday\")."
    },
    {
      "id": "q3",
      "text": "We ___ not yet received confirmation from the vendor.",
      "options": ["have", "had", "did", "are"],
      "correct_answer": "have",
      "explanation": "Present perfect with \"not yet\" describes a situation that started in the past and is still ongoing — the confirmation still hasn''t arrived."
    },
    {
      "id": "q4",
      "text": "The client ___ three follow-up emails since Monday.",
      "options": ["has sent", "sent", "was sending", "had sent"],
      "correct_answer": "has sent",
      "explanation": "Present perfect is used for actions that happened multiple times in an unfinished time period (since Monday = up to now)."
    },
    {
      "id": "q5",
      "text": "___ you manage to reproduce the bug in the staging environment?",
      "options": ["Did", "Have", "Had", "Do"],
      "correct_answer": "Did",
      "explanation": "Simple past is used here because the question refers to a specific, completed investigation attempt. If the time frame is open, \"Have you managed\" would also be correct."
    }
  ]
}', 'B2'),

('Prepositions in IT Business Context', 'Prepositions', 'Choose the correct preposition to complete these common IT business phrases.', '{
  "questions": [
    {
      "id": "q1",
      "text": "The app crashed ___ peak hours due to high traffic.",
      "options": ["during", "while", "for", "in"],
      "correct_answer": "during",
      "explanation": "\"During\" is used with noun phrases (peak hours, the meeting, the deployment). \"While\" is used with clauses (while it was running)."
    },
    {
      "id": "q2",
      "text": "I''ll follow up ___ you once I have the test results.",
      "options": ["with", "to", "for", "at"],
      "correct_answer": "with",
      "explanation": "\"Follow up with someone\" is the correct business English phrase. \"Follow up to\" is grammatically incorrect."
    },
    {
      "id": "q3",
      "text": "The feature is currently ___ development and will be ready next sprint.",
      "options": ["under", "in", "on", "at"],
      "correct_answer": "under",
      "explanation": "\"Under development\" is a fixed phrase meaning something is being actively worked on. Similarly: under review, under construction."
    },
    {
      "id": "q4",
      "text": "Please refer ___ the attached documentation for setup instructions.",
      "options": ["to", "at", "on", "for"],
      "correct_answer": "to",
      "explanation": "\"Refer to\" is the correct phrasal verb when directing someone to a resource or document."
    },
    {
      "id": "q5",
      "text": "We are ___ track to deliver the MVP by the end of the quarter.",
      "options": ["on", "in", "at", "with"],
      "correct_answer": "on",
      "explanation": "\"On track\" is a fixed business phrase meaning progressing as planned. Other fixed phrases: on schedule, on time, on budget."
    }
  ]
}', 'B1');
