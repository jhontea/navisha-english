-- Seed: Additional Grammar exercises batch 2

INSERT INTO navisha_english_grammar_exercises (title, topic, instruction, content, difficulty) VALUES

('Formal vs Informal Tone in Emails', 'Tone & Register', 'Choose the more professional and formal version of each sentence for use in business emails.', '{
  "questions": [
    {
      "id": "q1",
      "text": "Which sentence is more appropriate for a professional email to a client?",
      "options": [
        "Hey, just wanted to check if you got my last email?",
        "I am writing to follow up on my previous email dated June 20th.",
        "Did you see my email or what?",
        "Just checking in lol"
      ],
      "correct_answer": "I am writing to follow up on my previous email dated June 20th.",
      "explanation": "Professional emails use complete sentences, formal greetings, and avoid slang or casual phrases like \"Hey\" or \"lol\"."
    },
    {
      "id": "q2",
      "text": "Which is the most appropriate way to decline a meeting request professionally?",
      "options": [
        "Can''t make it, sorry.",
        "I am afraid I have a prior commitment at that time. Could we reschedule to Thursday afternoon?",
        "That time does not work for me at all.",
        "I''m busy then, pick another time."
      ],
      "correct_answer": "I am afraid I have a prior commitment at that time. Could we reschedule to Thursday afternoon?",
      "explanation": "\"I am afraid\" softens the refusal politely. Offering an alternative time shows professionalism and willingness to cooperate."
    },
    {
      "id": "q3",
      "text": "Which sentence best expresses urgency in a professional context?",
      "options": [
        "This is super urgent!!!",
        "I need this ASAP.",
        "I would appreciate your prompt attention to this matter, as it is time-sensitive.",
        "Please hurry up with this."
      ],
      "correct_answer": "I would appreciate your prompt attention to this matter, as it is time-sensitive.",
      "explanation": "Business writing expresses urgency through formal phrases like \"prompt attention\" and \"time-sensitive\" rather than exclamation marks or slang."
    },
    {
      "id": "q4",
      "text": "Which is the best opening for a cold outreach email to a potential client?",
      "options": [
        "Hey there! We make awesome software!",
        "I hope this email finds you well. I am reaching out regarding a solution that may be relevant to your team.",
        "You probably get a lot of emails but...",
        "Our product is the best in the market, check it out."
      ],
      "correct_answer": "I hope this email finds you well. I am reaching out regarding a solution that may be relevant to your team.",
      "explanation": "\"I hope this email finds you well\" is a standard professional opener. Stating your purpose clearly and relevantly shows respect for the reader''s time."
    },
    {
      "id": "q5",
      "text": "Which closing is most appropriate for a formal business proposal email?",
      "options": [
        "Talk soon!",
        "Cheers",
        "I look forward to your feedback and remain available should you have any questions.",
        "Let me know what you think, bye!"
      ],
      "correct_answer": "I look forward to your feedback and remain available should you have any questions.",
      "explanation": "A formal closing signals professionalism. It invites a response and offers availability without being pushy or overly casual."
    }
  ]
}', 'B1'),

('Gerunds and Infinitives in Business Writing', 'Gerunds & Infinitives', 'Choose the correct verb form (gerund or infinitive) to complete each business sentence.', '{
  "questions": [
    {
      "id": "q1",
      "text": "We recommend ___ the codebase before adding new features.",
      "options": ["refactoring", "to refactor", "refactor", "that refactor"],
      "correct_answer": "refactoring",
      "explanation": "\"Recommend\" is followed by a gerund (-ing form): recommend + doing something."
    },
    {
      "id": "q2",
      "text": "The team agreed ___ the deadline by one week.",
      "options": ["to extend", "extending", "extend", "that extending"],
      "correct_answer": "to extend",
      "explanation": "\"Agree\" is followed by an infinitive (to + verb): agree to do something."
    },
    {
      "id": "q3",
      "text": "We need to avoid ___ breaking changes in the public API.",
      "options": ["introducing", "to introduce", "introduce", "that introduce"],
      "correct_answer": "introducing",
      "explanation": "\"Avoid\" is always followed by a gerund (-ing form): avoid doing something."
    },
    {
      "id": "q4",
      "text": "The manager suggested ___ a post-mortem meeting after the outage.",
      "options": ["holding", "to hold", "hold", "that to hold"],
      "correct_answer": "holding",
      "explanation": "\"Suggest\" is followed by a gerund (-ing form): suggest doing something."
    },
    {
      "id": "q5",
      "text": "The client expects ___ a working prototype by the end of the month.",
      "options": ["to receive", "receiving", "receive", "that receive"],
      "correct_answer": "to receive",
      "explanation": "\"Expect\" is followed by an infinitive (to + verb): expect to do something."
    }
  ]
}', 'B2'),

('Linking Words in Technical Reports', 'Linking Words', 'Choose the most appropriate linking word or phrase to connect ideas in these technical report sentences.', '{
  "questions": [
    {
      "id": "q1",
      "text": "The initial load time was high. ___, we implemented server-side caching.",
      "options": ["Therefore", "Although", "Despite", "Unless"],
      "correct_answer": "Therefore",
      "explanation": "\"Therefore\" introduces a result or conclusion. The caching was implemented as a result of the high load time."
    },
    {
      "id": "q2",
      "text": "___ the deployment was delayed, the final product met all quality standards.",
      "options": ["Although", "Therefore", "Consequently", "Furthermore"],
      "correct_answer": "Although",
      "explanation": "\"Although\" introduces a contrast between two facts — the delay vs. the quality outcome."
    },
    {
      "id": "q3",
      "text": "The new architecture reduces latency. ___, it significantly lowers infrastructure costs.",
      "options": ["Furthermore", "However", "Despite", "Unless"],
      "correct_answer": "Furthermore",
      "explanation": "\"Furthermore\" adds an additional positive point to support the same argument."
    },
    {
      "id": "q4",
      "text": "We completed the backend integration on time. ___, the frontend team encountered unexpected delays.",
      "options": ["However", "Therefore", "Furthermore", "As a result"],
      "correct_answer": "However",
      "explanation": "\"However\" introduces a contrasting point — the backend was on time, but the frontend was not."
    },
    {
      "id": "q5",
      "text": "The database queries were not optimized. ___, the API response times exceeded acceptable limits.",
      "options": ["As a result", "Although", "Furthermore", "Nevertheless"],
      "correct_answer": "As a result",
      "explanation": "\"As a result\" shows a direct cause-and-effect relationship between the unoptimized queries and slow response times."
    }
  ]
}', 'B2');
