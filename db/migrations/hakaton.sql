-- Заполнение таблицы surveys
INSERT INTO surveys (title)
VALUES
    ('Customer Satisfaction Survey'),
    ('Employee Feedback Survey'),
    ('Product Review Survey');

-- Заполнение таблицы questions
INSERT INTO questions (title, type, surveyId)
VALUES
    ('How satisfied are you with our service?', 'SMILE', 1),
    ('Rate our service quality.', 'STARS', 1),
    ('Rate your overall experience.', 'RATE', 1),
    ('How happy are you with your job?', 'SMILE', 2),
    ('Rate your work environment.', 'STARS', 2),
    ('Rate the management support.', 'RATE', 2),
    ('How likely are you to recommend our product?', 'STARS', 3),
    ('Rate the product quality.', 'RATE', 3),
    ('How satisfied are you with the price?', 'SMILE', 3);

-- Заполнение таблицы answers (примерные данные)
INSERT INTO answers (questionId, userId, value)
VALUES
    (1, '550e8400-e29b-41d4-a716-446655440000', 2),
    (2, '550e8400-e29b-41d4-a716-446655440000', 4),
    (3, '550e8400-e29b-41d4-a716-446655440000', 8),
    (4, '550e8400-e29b-41d4-a716-446655440001', 3),
    (5, '550e8400-e29b-41d4-a716-446655440001', 5),
    (6, '550e8400-e29b-41d4-a716-446655440001', 9),
    (7, '550e8400-e29b-41d4-a716-446655440002', 4),
    (8, '550e8400-e29b-41d4-a716-446655440002', 10),
    (9, '550e8400-e29b-41d4-a716-446655440002', 3);