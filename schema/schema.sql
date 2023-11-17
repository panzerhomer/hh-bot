CREATE TABLE IF NOT EXISTS vacancies (
    id SERIAL PRIMARY KEY,
    city TEXT,
    company TEXT,
    industry TEXT,
    title TEXT,
    keywords TEXT,
    skills TEXT,
    experience TEXT,
    salary TEXT,
    url TEXT
)